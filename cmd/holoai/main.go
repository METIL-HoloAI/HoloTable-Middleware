package main

import (
	"database/sql"
	"fmt"
	"os"

	"log"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/database"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/endpoints/restapi"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/endpoints/websocket"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/listeners"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/utils"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/unityserver"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func main() {

	// Load configuration
	config.LoadYaml()

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

	database.Init(db)

	// Start WebSocket server
	go unityserver.StartWebSocketServer()
	<-unityserver.ClientReady
	log.Println("Unity client connected")

	if config.General.OpenWebsocket {
		go websocket.EstablishConnection()
		restapi.StartRestAPI()
		utils.WaitForInterrupt()
	} else {
		listeners.StartTextListener()
	}
}

func ReadFileBytes(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	fileSize := fileInfo.Size()
	data := make([]byte, fileSize)

	_, err = file.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}
