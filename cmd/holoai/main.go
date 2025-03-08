package main

import (
	"database/sql"
	"os"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/database"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/listeners"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func main() {
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

	// Check how user wants to listen for input
	// and start that listener
	if config.General.Listener == "mic" {
		logrus.Info("Microphone Listener")
	} else if config.General.Listener == "text" {
		listeners.StartTextListener()
	} else {
		logrus.Error("Invalid listener option in general.yaml")
	}
}
