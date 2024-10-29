package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type MessageEnvelope struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Message string `json:"message"`
}

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
		q := r.URL.Query()
		localClientId := q.Get("clientId")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("client connected. clientId: %s", localClientId)

		connMap[localClientId] = conn

		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
					log.Printf("client closed. clientId: %s", localClientId)
					delete(connMap, localClientId)
				}
				log.Println(err)
				return
			}

			var d MessageEnvelope
			if err := json.Unmarshal(p, &d); err != nil {
				log.Printf("error unmarshalling message: %v", err)
				continue
			}
			log.Printf("message: %s", d)

			toClientId := d.To
			toConn, found := connMap[toClientId]
			if !found {
				log.Printf("client not found: %s", toClientId)
				continue
			}

			m := MessageEnvelope{To: d.To, From: localClientId, Message: d.Message}
			b, err := json.Marshal(m)
			if err != nil {
				log.Printf("error marshalling message: %v", err)
				continue
			}

			if err := toConn.WriteMessage(websocket.TextMessage, b); err != nil {
				log.Printf("error writing message: %v", err)
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
