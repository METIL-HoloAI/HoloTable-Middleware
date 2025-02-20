package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// In production, be sure to check origins properly.
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to a WebSocket connection.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	// Continuously read messages from the client
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// Process only binary messages (the audio data)
		if messageType == websocket.BinaryMessage {
			log.Printf("Received %d bytes of audio data", len(message))
			// Here you might decode or forward the audio data for further processing
		} else {
			log.Println("Non-binary message received; ignoring.")
		}
	}
}

func EstablishConnection() {
	http.HandleFunc("/ws/audio", wsHandler)
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
