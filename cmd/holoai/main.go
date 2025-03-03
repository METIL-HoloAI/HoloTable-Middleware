package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/database"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/listeners"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load configuration from YAML.
	config.LoadYaml()

	// Create the data directory if it doesn't exist.
	if err := os.MkdirAll(config.General.DataDir, os.ModePerm); err != nil {
		log.Fatal("Failed to create data directory:", err)
	}

	// Open and initialize the database.
	dbPath := filepath.Join(config.General.DataDir, "database.db")
	db, err := sql.Open("sqlite3", dbPath+"?_mode=shared&_journal_mode=WAL")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()
	database.Init(db)

	// Start listener based on configuration.
	switch config.General.Listener {
	case "mic":
		fmt.Println("Microphone Listener")
	case "text":
		listeners.StartTextListener()
	default:
		fmt.Println("Invalid listener option in general.yaml")
	}

	// fakeJSONData := []byte(`{
	// 	"contentType": "3d",
	// 	"requiredParameters": {
	// 		"mode": "preview",
	// 		"prompt": "a monster mask"
	// 	},
	// 	"optionalParameters": {
	// 		"art_style": "realistic",
	// 		"seed": null,
	// 		"ai_model": "meshy-4",
	// 		"topology": "triangle",
	// 		"target_polycount": 30000,
	// 		"should_remesh": true,
	// 		"symmetry_mode": "auto",
	// 		"enable_pbr": false
	// 	}
	// }`)

	// callers.LoadIntentDetectionResponse(fakeJSONData)
}
