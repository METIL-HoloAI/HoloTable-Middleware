package listeners

import (
	"log"
	"strings"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
	"github.com/gorilla/websocket"
)

func TranscribeAudio(audio []byte, clientWebsocket *websocket.Conn) {
	vosk, _, err := websocket.DefaultDialer.Dial(config.SpeechToText.WebsocketURL, nil)
	if err != nil {
		log.Fatal("Failed to open vosk websocket connection, ", err)
	}
	defer vosk.Close()

	err = vosk.WriteMessage(websocket.BinaryMessage, audio)
	if err != nil {
		log.Fatal("Failed to send audio to vosk, ", err)
	}

	var response string
	for {
		_, responsebytes, err := vosk.ReadMessage()
		if err != nil {
			log.Fatal("Failed to read message from vosk, ", err)
		}

		if responsebytes != nil {
			response = string(responsebytes)
			break
		}
	}

	if strings.Contains(response, config.SpeechToText.Keyword) {
		// TODO: send something back to user websocket to tell them a keyword has been found
		// TODO: Record next sentence from user
		log.Print("Found keyword")
	}
}
