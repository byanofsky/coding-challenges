package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"myredis/internal"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Server configuration
type Config struct {
	Address         string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxMessageSize  int
	ShutdownTimeout time.Duration
}

// TCP server
type Server struct {
	config     Config
	listener   net.Listener
	logger     *slog.Logger
	handler    CommandHandler
	shutdownWg sync.WaitGroup
}

type KVRecord struct {
	value  string
	expire bool
	ttl    time.Time
}

// Dictionary stores key value pairs
type Dictionary struct {
	m  *sync.RWMutex
	kv map[string]KVRecord
}

type SetCommandOptions struct {
	ex        bool
	exSeconds int
}

// CommandHandler defines the interface for handling commands
type CommandHandler interface {
	Handle(ctx context.Context, command string, args []internal.Data) (*internal.Data, error)
}

// DefaultCommandHandler implements basic command handling
type DefaultCommandHandler struct {
	dict Dictionary
}

// TODO: Return pointer?
func NewDictionary() Dictionary {
	return Dictionary{m: &sync.RWMutex{}, kv: make(map[string]KVRecord)}
}

// TODO: Add mutex
func (d *Dictionary) Set(k string, v string) {
	d.m.Lock()
	defer d.m.Unlock()

	d.kv[k] = KVRecord{value: v, ttl: time.Time{}, expire: false}
}

func (d *Dictionary) SetWithExpire(k string, v string, expireMs int) {
	d.m.Lock()
	defer d.m.Unlock()

	ttl := time.Now().Add(time.Duration(expireMs) * time.Millisecond)

	d.kv[k] = KVRecord{value: v, ttl: ttl, expire: true}
}

func (d *Dictionary) Get(k string) (string, bool) {
	d.m.RLock()
	defer d.m.RUnlock()
	record, ok := d.kv[k]
	if record.expire {
		if !(time.Now().Before(record.ttl)) {
			return "", false
		}
	}
	return record.value, ok
}

// NewServer creates a new server instance
func NewServer(config Config, logger *slog.Logger, handler CommandHandler) *Server {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	return &Server{
		config:  config,
		logger:  logger,
		handler: handler,
	}
}

// Start begins listening for connections
func (s *Server) Start(ctx context.Context) error {
	var err error
	s.listener, err = net.Listen("tcp", s.config.Address)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	s.logger.Info("server started", "address", s.config.Address)

	go s.acceptConnections(ctx)
	return nil
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.listener.Close(); err != nil {
		return fmt.Errorf("failed to close listener: %w, err", err)
	}

	done := make(chan struct{})
	go func() {
		s.shutdownWg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func (s *Server) acceptConnections(ctx context.Context) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			s.logger.Error("failed to accept connection", "error", err)
		}

		s.shutdownWg.Add(1)
		go func() {
			defer s.shutdownWg.Done()
			s.handleConnection(ctx, conn)
		}()
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	logger := s.logger.With(
		"remote_addr", conn.RemoteAddr().String(),
	)
	logger.Info("new connection established")

	for {
		// TODO: Investigate whether read deadline is correct appraoch.
		// If it is, gracefully handle read request after deadline.
		if err := conn.SetReadDeadline(time.Now().Add(s.config.ReadTimeout)); err != nil {
			logger.Error("failed to set read deadline", "error", err)
		}

		request, err := s.readRequest(conn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			logger.Error("failed to read request", "error", err)
			s.sendError(conn, "failed to read request")
			return
		}

		if err := s.processRequest(ctx, conn, request); err != nil {
			logger.Error("failed to process request", "error", err)
			s.sendError(conn, "internal server error")
			return
		}
	}
}

func (s *Server) readRequest(conn net.Conn) (*internal.Data, error) {
	buffer := make([]byte, s.config.MaxMessageSize)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	request := string(buffer[:n])

	s.logger.Info("request received", "request", request)

	return internal.Deserialize(request)
}

func (s *Server) processRequest(ctx context.Context, conn net.Conn, request *internal.Data) error {
	// TODO: Use ok instead of error
	command, err := request.GetArray()
	if err != nil {
		return fmt.Errorf("failed to get command array: %w", err)
	}

	if len(command) < 1 {
		return s.sendError(conn, "empty command")
	}

	cmdStr, err := command[0].GetString()
	if err != nil {
		return fmt.Errorf("failed to get command string: %w", err)
	}

	response, err := s.handler.Handle(ctx, strings.ToUpper(cmdStr), command[1:])
	if err != nil {
		return s.sendError(conn, err.Error())
	}

	return s.sendResponse(conn, response)
}

func (s *Server) sendResponse(conn net.Conn, response *internal.Data) error {
	if err := conn.SetWriteDeadline(time.Now().Add(s.config.WriteTimeout)); err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	serialized, err := internal.Serialize(*response)
	if err != nil {
		return fmt.Errorf("failed to serialize response: %w", err)
	}

	s.logger.Info("send response", "response", serialized)

	_, err = conn.Write([]byte(serialized))
	return err
}

func (s *Server) sendError(conn net.Conn, message string) error {
	return s.sendResponse(conn, internal.NewSimpleError(message))
}

// Handle implements the CommandHandler interface for DefaultCommandHanlder
func (h *DefaultCommandHandler) Handle(ctx context.Context, command string, args []internal.Data) (*internal.Data, error) {
	switch command {
	case "PING":
		return internal.NewSimpleStringData("PONG"), nil
	case "ECHO":
		return internal.NewArrayData(args), nil
	case "COMMAND":
		return internal.NewSimpleStringData("CONNECTED"), nil
	case "SET":
		return h.handleSetCommand(args)
	case "GET":
		return h.handleGetCommand(args)
	case "HELLO":
		return h.handleHelloCommand(args)
	default:
		return nil, fmt.Errorf("unknown command: %s", command)
	}
}

func (h *DefaultCommandHandler) handleSetCommand(args []internal.Data) (*internal.Data, error) {
	// Validate input
	if len(args) < 2 {
		return nil, fmt.Errorf("invalid set command")
	}

	key, err := args[0].GetString()
	if err != nil {
		return nil, fmt.Errorf("first arg must be string")
	}
	// TODO: Relax string requirement
	value, err := args[1].GetString()
	if err != nil {
		return nil, fmt.Errorf("second arg must be string")
	}

	options := SetCommandOptions{}
	for i := 2; i < len(args); i++ {
		option, err := args[i].GetString()
		if err != nil {
			// TODO: include command in error for logging
			return nil, fmt.Errorf("invalid option flag")
		}
		switch strings.ToUpper(option) {
		case "EX":
			i++
			options.ex = true
			if i >= len(args) {
				return nil, fmt.Errorf("command error. EX must be followed by int")
			}
			s, err := args[i].GetString()
			if err != nil {
				// TODO: use generic error
				return nil, fmt.Errorf("EX not followed by int")
			}
			seconds, err := strconv.Atoi(s)
			if err != nil {
				// TODO: use generic error
				return nil, fmt.Errorf("EX not followed by int")
			}
			options.exSeconds = seconds
		}
	}
	// TODO: Remove
	// fmt.Printf("Options: %v", options)

	if options.ex {
		h.dict.SetWithExpire(key, value, options.exSeconds*1000)
	} else {
		h.dict.Set(key, value)
	}

	return internal.NewBulkStringData("OK"), nil
}

func (h *DefaultCommandHandler) handleGetCommand(args []internal.Data) (*internal.Data, error) {
	// Validate input
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid get command")
	}

	key, err := args[0].GetString()
	if err != nil {
		return nil, fmt.Errorf("first arg must be string")
	}

	value, ok := h.dict.Get(key)
	if !ok {
		return internal.NewNullData(), nil
	}
	return internal.NewBulkStringData(value), nil
}

func (h *DefaultCommandHandler) handleHelloCommand(args []internal.Data) (*internal.Data, error) {
	return internal.NewArrayData([]internal.Data{
		*internal.NewBulkStringData("server"),
		*internal.NewBulkStringData("redis"),
	}), nil
}

func main() {
	config := Config{
		Address:         "localhost:6379",
		ReadTimeout:     30 * time.Minute,
		WriteTimeout:    30 * time.Second,
		MaxMessageSize:  1024 * 1024, // 1MB
		ShutdownTimeout: 30 * time.Second,
	}

	logger := slog.New((slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	server := NewServer(config, logger, &DefaultCommandHandler{dict: NewDictionary()})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := server.Start(ctx); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	<-sigChan
	logger.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("failed to shutdown server gracefully", "error", err)
		os.Exit(1)
	}
}
