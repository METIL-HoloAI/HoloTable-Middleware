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

func GenerateAndSendContent() {
	// test 1
	// fileName := "v2"
	// extension := "obj"
	// filePath := "C:\\Users\\anala\\Desktop\\HoloAIDocuments\\Assets\\v2\\source\\v2\\v2.obj"
	// ExportAsset(fileName, extension, filePath)
	// filePath := filepath.Join(ASSETS_DIR, fmt.Sprintf("%s.%s", fileName, extension))
	fileName := "miau"
	extension := "gif"
	filePath := "C:\\Users\\anala\\Desktop\\HoloAIDocuments\\Assets\\miau.gif"
	ExportAsset(fileName, extension, filePath)
	fileName = "bunny"
	extension = "gif"
	filePath = "C:\\Users\\anala\\Desktop\\HoloAIDocuments\\Assets\\bunny.gif"
	ExportAsset(fileName, extension, filePath)
	fileName = "catLion"
	extension = "jpeg"
	filePath = "C:\\Users\\anala\\Desktop\\HoloAIDocuments\\Assets\\catLion.jpeg"
	ExportAsset(fileName, extension, filePath)
	ExportAsset(fileName, extension, filePath)
	fileName = "ghostwire"
	extension = "png"
	filePath = "C:\\Users\\anala\\Desktop\\HoloAIDocuments\\Assets\\ghostwire.png"
	ExportAsset(fileName, extension, filePath)
	// ExportAsset(fileName, extension, filePath)
	// test 2
	fileName = "blueMan"
	extension = "glb"
	filePath = "C:\\Users\\anala\\Desktop\\HoloAIDocuments\\Assets\\blueMan.glb"
	ExportAsset(fileName, extension, filePath)
	ExportAsset(fileName, extension, filePath)
	fileName = "blueGirl"
	extension = "glb"
	filePath = "C:\\Users\\anala\\Desktop\\HoloAIDocuments\\Assets\\blueGirl.glb"
	ExportAsset(fileName, extension, filePath)
	// fileName = "baby"
	// extension = "mp4"
	// filePath = "C:\\Users\\anala\\Desktop\\HoloAIDocuments\\Assets\\baby.mp4"
	fileName = "fish"
	extension = "mp4"
	filePath = "C:\\Users\\anala\\Desktop\\HoloAIDocuments\\Assets\\fish.mp4"
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
