package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
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
			m, err := strconv.Unquote(string(p))
			if err != nil {
				fmt.Printf("error processing message: %v", err)
			}
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

				toConn, exists := connMap[clientToFind]
				var bytes []byte
				if exists {
					bytes = []byte("ok")
				} else {
					bytes = []byte("not found")
				}
				if err := conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
					log.Printf("error writing message %v", err)
				}
				after, found := strings.CutPrefix(m, "findClient:")
				if !found {
					log.Println("unable to find findClient:")
				}
				message := fmt.Sprintf("%s:%s", "remote", after)
				if err := toConn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Printf("error writing message %v", err)
				}
			}
			if strings.Contains(m, "answer") {
				parts := strings.Split(m, ":")
				toClientId := parts[1]
				log.Printf("answer. requester: %s. to: %s", clientId, toClientId)

				toConn, exists := connMap[toClientId]
				if !exists {
					log.Printf("cannot find client %s", toClientId)
				}
				if err := toConn.WriteMessage(websocket.TextMessage, []byte(m)); err != nil {
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
