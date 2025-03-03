package unityserver

import (
	"bufio"
	"log"
	"net/http"
	"os"

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

	// Listen for "close" message from the terminal
	go listenForClose()
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

	// Wait for messages from the Unity client
	for {
		_, msg, err := Conn.ReadMessage()
		if err != nil {
			log.Println("Read Error:", err)
			return
		}

		if string(msg) == "READY" {
			log.Println("Unity client is ready")
			ClientReady <- true
		} else {
			log.Printf("Received message: %s", msg)
		}
	}
}

func listenForClose() {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		if input == "close\n" {
			log.Println("Received close message from terminal, closing connection")
			if Conn != nil {
				Conn.Close()
			}
			os.Exit(0)
		}
	}
}
