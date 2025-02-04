package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Create an upgrader with a simple origin check.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// For production, replace this with a proper origin check.
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Echo loop: this reads the message sent (from terminal where wscat is connected to go server) and writes it back (in terminal where main.go is running)
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received: %s", message)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func EstablishConnection() {
	http.HandleFunc("/ws", wsHandler)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
