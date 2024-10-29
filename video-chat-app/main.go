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
	To      string
	From    string
	Message string
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
			// if err := conn.WriteMessage(messageType, p); err != nil {
			// 	log.Println(err)
			// 	return
			// }
			// m, err := strconv.Unquote(string(p))
			// if err != nil {
			// 	fmt.Printf("error processing message: %v", err)
			// }
			// if strings.Contains(m, "findClient") {
			// 	parts := strings.Split(m, ":")
			// 	clientToFind := parts[1]
			// 	log.Printf("findClient. requester: %s. to: %s", localClientId, clientToFind)

			// 	toConn, exists := connMap[clientToFind]
			// 	var bytes []byte
			// 	if exists {
			// 		bytes = []byte("ok")
			// 	} else {
			// 		bytes = []byte("not found")
			// 	}
			// 	if err := conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
			// 		log.Printf("error writing message %v", err)
			// 	}
			// 	after, found := strings.CutPrefix(m, "findClient:")
			// 	if !found {
			// 		log.Println("unable to find findClient:")
			// 	}
			// 	message := fmt.Sprintf("%s:%s", "remote", after)
			// 	if err := toConn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			// 		log.Printf("error writing message %v", err)
			// 	}
			// }
			// if strings.Contains(m, "answer") {
			// 	parts := strings.Split(m, ":")
			// 	toClientId := parts[1]
			// 	log.Printf("answer. requester: %s. to: %s", localClientId, toClientId)

			// 	toConn, exists := connMap[toClientId]
			// 	if !exists {
			// 		log.Printf("cannot find client %s", toClientId)
			// 	}
			// 	if err := toConn.WriteMessage(websocket.TextMessage, []byte(m)); err != nil {
			// 		log.Printf("error writing message %v", err)
			// 	}
			// }
		}
	})

	// Start server
	log.Println("Server starting on :8000...")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
