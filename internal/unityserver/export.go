package unityserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/unityserver/assetstruct"
	"github.com/gorilla/websocket"
)

const (
	ASSETS_DIR = "./src/3dModelsTest" // Asset directory
)

func GenerateAndSendContent() {
	// test 1
	// fileName := "blueMan"
	// extension := "glb"

	// test 2
	fileName := "catLion"
	extension := "jpeg"
	fileData, err := ioutil.ReadFile(filepath.Join(ASSETS_DIR, fmt.Sprintf("%s.%s", fileName, extension)))
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	ExportAsset(fileName, extension, fileData)
}

func ExportAsset(fileName, extension string, fileData []byte) {
	assetMsg := assetstruct.AssetMessage{
		Type:      "asset",
		Name:      fileName,
		Extension: extension,
		FileData:  fileData,
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
