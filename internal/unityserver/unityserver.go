package unityserver

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var ClientReady = make(chan bool)
var Conn *websocket.Conn

func StartWebSocketServer() {
	go startWebSocketServer()

	// HANDSHAKE MECHANISM
	<-ClientReady // This will block until the Unity client sends a "READY" message
}

func startWebSocketServer() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("WebSocket server running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	var err error
	Conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}
	log.Println("WebSocket client connected")

	// Wait for a handshake message from the Unity client
	_, msg, err := Conn.ReadMessage()
	if err != nil {
		log.Println("Read Error:", err)
		return
	}

	if string(msg) == "READY" {
		log.Println("Unity client is ready")
		ClientReady <- true
	}
}
