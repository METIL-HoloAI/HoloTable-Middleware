package unityserver

import (
	"encoding/json"
	"log"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/unityserver/assetstruct"
	"github.com/gorilla/websocket"
)

const (
	ASSETS_DIR = "./internal/unityserver/3dModelsTest" // Asset directory
)

func ExportAsset(fileName, extension, filePath string) {
	assetMsg := assetstruct.AssetMessage{
		Type:      "asset",
		Name:      fileName,
		Extension: extension,
		FilePath:  filePath,
	}

	response, err := json.Marshal(assetMsg)
	if err != nil {
		log.Println("Failed to marshal asset message:", err)
		return
	}

	SendToUnity(response)
}

func SendToUnity(response []byte) {
	log.Println("Sending message to Unity")
	if Conn == nil {
		log.Println("Connection is not initialized")
		return
	}
	if err := Conn.WriteMessage(websocket.TextMessage, response); err != nil {
		log.Println("Write Error:", err)
	}
}
