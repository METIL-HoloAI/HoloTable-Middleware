package main

import (
	"database/sql"
	"os"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/endpoints/restapi"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/endpoints/websocket"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/listeners"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/unityserver"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/utils"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func main() {
	// config.LoadYaml()

	// Initialize logger
	config.InitLogger()

	if err := os.MkdirAll(config.General.DataDir, os.ModePerm); err != nil {
		logrus.Fatal("Failed to create data directory:", err)
	}

	db, err := sql.Open("sqlite3", config.General.DataDir+"database.db?_mode=shared&_journal_mode=WAL")
	if err != nil {
		logrus.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// database.Init(db)

	// Start WebSocket server // IMPORTANT TO ADD THIS
	go unityserver.StartWebSocketServer()
	<-unityserver.ClientReady

	// unityserver.ExportAsset(fileName, extension, filePath) // HOW TO call my function

	if config.General.OpenWebsocket {
		go websocket.EstablishConnection()
		restapi.StartRestAPI()
		utils.WaitForInterrupt()
	} else {
		listeners.StartTextListener()
	}
}
