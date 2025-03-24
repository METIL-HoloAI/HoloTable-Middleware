package main

import (
	// "database/sql"
	// "os"

	// "github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
	// "github.com/METIL-HoloAI/HoloTable-Middleware/internal/database"
	// "github.com/METIL-HoloAI/HoloTable-Middleware/internal/endpoints/restapi"
	// "github.com/METIL-HoloAI/HoloTable-Middleware/internal/endpoints/websocket"
	// "github.com/METIL-HoloAI/HoloTable-Middleware/internal/listeners"
	"log"
	"time"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/unityserver"
	// "github.com/METIL-HoloAI/HoloTable-Middleware/internal/utils"
	// _ "github.com/mattn/go-sqlite3"
	// "github.com/sirupsen/logrus"
)

func main() {

	// // Load configuration
	// config.LoadYaml()

	// // Initialize logger
	// config.InitLogger()

	// if err := os.MkdirAll(config.General.DataDir, os.ModePerm); err != nil {
	// 	logrus.Fatal("Failed to create data directory:", err)
	// }

	// db, err := sql.Open("sqlite3", config.General.DataDir+"database.db?_mode=shared&_journal_mode=WAL")
	// if err != nil {
	// 	logrus.Fatal("Failed to open database:", err)
	// }
	// defer db.Close()

	// database.Init(db)

	// Start WebSocket server
	go unityserver.StartWebSocketServer()
	log.Println("Waiting for Unity client to connect...")
	<-unityserver.ClientReady
	log.Println("Unity client connected")

	unityserver.ExportAsset("catLion", "PNG", "../test/catLion.PNG")
	time.Sleep(3 * time.Second)
	unityserver.ExportAsset("table", "glb", "../test/table.glb")
	time.Sleep(5 * time.Second)
	unityserver.ExportAsset("movingImage", "mp4", "../test/movingImage.mp4")
	time.Sleep(5 * time.Second)
	unityserver.ExportAsset("Dragonflying2", "mp4", "../test/Dragonflying2.mp4")
	time.Sleep(5 * time.Second)
	unityserver.ExportAsset("Cat2", "mp4", "../test/Cat2.mp4")
	time.Sleep(5 * time.Second)
	unityserver.ExportAsset("cat", "mp4", "../test/cat.mp4")
	time.Sleep(5 * time.Second)
	// unityserver.ExportAsset(fileName, extension, filePath) // HOW TO call my function

	// if config.General.OpenWebsocket {
	// 	go websocket.EstablishConnection()
	// 	restapi.StartRestAPI()
	// 	utils.WaitForInterrupt()
	// } else {
	// 	listeners.StartTextListener()
	// }
}
