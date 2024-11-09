package main

import (
	"io"
	"log"
	"myredis/internal"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":6379")
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
		data := buffer[:n]
		log.Printf("received %d bytes: %q\n", n, string(data))

		d, err := internal.Serialize(*internal.NewSimpleStringData("PONG"))
		if err != nil {
			log.Printf("error deserializing: %v", err)
		}
		_, err = conn.Write([]byte(d))
		if err != nil {
			log.Printf("error writing: %v", err)
		}
	}
}
