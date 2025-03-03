package unityserver

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/unityserver/assetstruct"
	"github.com/gorilla/websocket"
)

const (
	ASSETS_DIR = "./internal/unityserver/3dModelsTest" // Asset directory
)

func GenerateAndSendContent() {
	// test 1
	fileName := "catLion"
	extension := "jpeg"
	filePath := filepath.Join(ASSETS_DIR, fmt.Sprintf("%s.%s", fileName, extension))
	ExportAsset(fileName, extension, filePath)

	// test 2
	fileName = "blueMan"
	extension = "glb"
	filePath = filepath.Join(ASSETS_DIR, fmt.Sprintf("%s.%s", fileName, extension))
	ExportAsset(fileName, extension, filePath)
}

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
