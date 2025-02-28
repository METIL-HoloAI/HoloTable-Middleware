package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/database"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/endpoints/websocket"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/listeners"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/utils"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config.LoadYaml()

	if err := os.MkdirAll(config.General.DataDir, os.ModePerm); err != nil {
		log.Fatal("Failed to create data directory:", err)
	}

	db, err := sql.Open("sqlite3", config.General.DataDir+"database.db?_mode=shared&_journal_mode=WAL")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	database.Init(db)

	if config.General.OpenWebsocket {
		go websocket.EstablishConnection()
		utils.WaitForInterrupt()
	} else {
		listeners.StartTextListener()
	}
}
