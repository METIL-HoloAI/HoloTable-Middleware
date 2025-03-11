package websocket

import (
	"log"
	"net/http"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/listeners"
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

	// Start the vosk client
	listeners.InitializeVosk()
	defer listeners.CloseVosk()

	voskResponse := make(chan string)
	quitVosk := make(chan bool)

	go listeners.GetResponse(voskResponse, quitVosk)

	// Continuously read messages from the client
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// Process only binary messages (the audio data)
		if messageType == websocket.BinaryMessage {
			listeners.SendAudio(message)

			select {
			case text := <-voskResponse:
				log.Printf("Got response from Vosk: %s", text)
				if listeners.CheckForKeyword(text) {
					err = conn.WriteMessage(websocket.TextMessage, []byte("Keyword Detected"))
					if err != nil {
						log.Println("Failed to send keyword detected message to client:", err)
					}
				}
			default:
				// Keep going
			}
		} else if messageType == websocket.TextMessage {
			log.Println("Text message recieved: ", message)
		} else {
			log.Println("Non-binary message received; ignoring.")
		}
	}

	quitVosk <- true
}

func EstablishConnection() {
	http.HandleFunc("/ws/audio", wsHandler)
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
