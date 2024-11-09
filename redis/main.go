package main

import (
	"io"
	"log"
	"myredis/internal"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", "localhost:6379")
	if err != nil {
		log.Fatalf("error startup: %v", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("error accepting conn: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("new connection from %s\n", conn.RemoteAddr().String())

	buffer := make([]byte, 1024) // 1KB buffer
	for {
		n, err := conn.Read(buffer)
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Println("error reading:", err)
			return
		}

		// Process the received data
		received := string(buffer[:n])
		log.Printf("received %d bytes: %q\n", n, received)

		request, err := internal.Deserialize(received)
		if err != nil {
			log.Printf("error deserializing: %q", received)
			continue
		}

		handleRequest(request, conn)
	}
}

func handleRequest(request *internal.Data, conn net.Conn) {
	switch request.GetKind() {
	case internal.ArrayKind:
		command, err := request.GetArray()
		if err != nil {
			// TODO: Return error
			log.Printf("error getting array: %v", err)
			return
		}

		handleCommand(conn, command[0], command[1:])
	default:
		log.Printf("error unhandled kind: %s", request.GetKind())
		respondError(conn)
	}
}

func handleCommand(conn net.Conn, commandData internal.Data, args []internal.Data) {
	command, err := commandData.GetString()
	if err != nil {
		log.Printf("error command should be simple string: %v", command)

	}
	switch command {
	case "PING":
		response, err := internal.Serialize(*internal.NewSimpleStringData("PONG"))
		if err != nil {
			log.Printf("error serializing: %v", err)
			return
		}
		sendResponse(response, conn)
	case "COMMAND":
		response, err := internal.Serialize(*internal.NewSimpleStringData("CONNECTED"))
		if err != nil {
			log.Printf("error serializing: %v", err)
			return
		}
		sendResponse(response, conn)
	case "ECHO":
		response, err := internal.Serialize(*internal.NewArrayData(args))
		if err != nil {
			log.Printf("error serializing: %v", err)
			return
		}
		sendResponse(response, conn)
	default:
		log.Printf("unknown command: %s", command)
	}
}

func respondError(conn net.Conn) {
	response, err := internal.Serialize(*internal.NewSimpleError("unexpected command"))
	if err != nil {
		log.Printf("error serializing: %v", err)
	}
	sendResponse(response, conn)
}

func sendResponse(response string, conn net.Conn) {
	_, err := conn.Write([]byte(response))
	// TODO: Return error
	if err != nil {
		log.Printf("error writing: %v", err)
	}
	log.Printf("write response: %q", response)
}
