package listeners

import (
	"encoding/json"
	"log"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
	"github.com/gorilla/websocket"
)

// Models vosk server response
type Message struct {
	Result []struct {
		Conf  float64 `json:"conf"`
		End   float64 `json:"end"`
		Start float64 `json:"start"`
		Word  string  `json:"word"`
	} `json:"result,omitempty"`
	Text    string `json:"text,omitempty"`
	Partial string `json:"partial,omitempty"`
}

var vosk *websocket.Conn

func SendAudio(audio []byte) {
	err := vosk.WriteMessage(websocket.BinaryMessage, audio)
	if err != nil {
		log.Print("Failed to send audio to vosk, ", err)
	}
}

func InitializeVosk() {
	var err error
	vosk, _, err = websocket.DefaultDialer.Dial(config.SpeechToText.VoskURL, nil)
	if err != nil {
		log.Fatal("Failed to open Vosk WebSocket connection: ", err)
	}

	configMsg := map[string]any{
		"config": map[string]any{
			"sample_rate": 16000,
		},
	}
	configBytes, _ := json.Marshal(configMsg)
	err = vosk.WriteMessage(websocket.TextMessage, configBytes)
	if err != nil {
		log.Fatal("Failed to send config to Vosk: ", err)
	}
}

func CloseVosk() {
	vosk.Close()
}

func GetResponse(response chan string, quit chan bool) {
	done := make(chan bool)

	go func() {
		<-quit
		done <- true
	}()

	for {
		// Check if we should quit
		select {
		case <-done:
			log.Println("Stopping Vosk response listener")
			return
		default:
			// Continue processing
		}

		// Read from Vosk (this is blocking)
		_, jsonMessage, err := vosk.ReadMessage()
		if err != nil {
			log.Println("Error reading from Vosk:", err)
			return
		}

		var message Message
		err = json.Unmarshal(jsonMessage, &message)
		if err != nil {
			log.Println("Failed to unmarshal JSON from Vosk:", err)
			continue
		}

		// Only send non-empty final results
		if message.Text != "" {
			log.Printf("Sending text to channel: %s", message.Text)
			response <- message.Text
		}
	}
}
