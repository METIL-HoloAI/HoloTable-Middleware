package listeners

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
	"github.com/gorilla/websocket"
)

// Models vosk server response
type Message struct {
	Result []struct {
		Conf  float64
		End   float64
		Start float64
		Word  string
	}
	Text string
}

var vosk *websocket.Conn

func InitializeVosk() {
	// Error has to be declared out here to avoid an error
	// of vosk not being used
	var err error
	vosk, _, err = websocket.DefaultDialer.Dial(config.SpeechToText.LiveTranscription.WebsocketURL, nil)
	if err != nil {
		log.Fatal("Failed to open Vosk WebSocket connection: ", err)
	}
}

func CloseVosk() {
	vosk.Close()
}

func TranscribeAudio(audio []byte) bool {
	err := vosk.WriteMessage(websocket.BinaryMessage, audio)
	if err != nil {
		log.Fatal("Failed to send audio to vosk, ", err)
	}

	_, _, err = vosk.ReadMessage()
	if err != nil {
		log.Fatal("Failed to recieve response from vosk for audio, ", err)
	}

	err = vosk.WriteMessage(websocket.TextMessage, []byte("{\"eof\" : 1}"))
	if err != nil {
		log.Fatal("Failed to send vosk EOF, ", err)
	}

	_, jsonMessage, err := vosk.ReadMessage()
	if err != nil {
		log.Fatal("Failed to recieve final message from vosk, ", err)
	}

	var message Message
	err = json.Unmarshal(jsonMessage, &message)
	if err != nil {
		log.Fatal("Failed to unmarshal json from vosk, ", err)
	}

	log.Print(message)

	return strings.Contains(message.Text, config.SpeechToText.LiveTranscription.Keyword)
}
