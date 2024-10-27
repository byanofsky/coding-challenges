package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// Using defaults from https://pkg.go.dev/github.com/gorilla/websocket#hdr-Overview
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var connMap = make(map[string]*websocket.Conn)

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.Handle("/", fs)

	// API routes
	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	http.HandleFunc("/signal", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		var clientId string

		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
					log.Printf("Close client %s", clientId)
					delete(connMap, clientId)
				}
				log.Println(err)
				return
			}
			// if err := conn.WriteMessage(messageType, p); err != nil {
			// 	log.Println(err)
			// 	return
			// }
			m := string(p)
			// TODO: More robust handling of message protocol
			if strings.Contains(m, "clientId") {
				if clientId != "" {
					log.Printf("client called twice")
					conn.Close()
					return
				}
				parts := strings.Split(m, ":")
				clientId = parts[1]
				connMap[clientId] = conn
				log.Printf("client connected: %s", clientId)
			}
			if strings.Contains(m, "findClient") {
				parts := strings.Split(m, ":")
				clientToFind := parts[1]
				log.Printf("findClient. requester: %s. to: %s", clientId, clientToFind)

				_, exists := connMap[clientToFind]
				var bytes []byte
				if exists {
					bytes = []byte("ok")
				} else {
					bytes = []byte("not found")
				}
				if err := conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
					log.Printf("error writing message %v", err)
				}
			}
		}
	})

	// Start server
	log.Println("Server starting on :8000...")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
