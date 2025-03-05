package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/callers"
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

	// Sample DALL-E JSON response with an actual working image URL.
	// Note: The JSON also contains an "id" field for demonstration.
	sampleJSON :=
		`{
		"art_style": "realistic",
		"created_at": 1.741027287777e+12,
		"finished_at": 1.741027394517e+12,
		"id": "01955d50-221d-7920-bcb3-fdac7cac8b41",
		"mode": "refine",
		"model_urls": {
			"fbx": "https://assets.meshy.ai/d78131f0-3ed1-42ee-9ff4-92299dd3ea96/tasks/output/model.fbx",
			"glb": "https://assets.meshy.ai/d78131f0-3ed1-42ee-9ff4-92299dd3ea96/tasks/output/model.glb"
		},
		"status": "SUCCEEDED",
		"thumbnail_url": "https://assets.meshy.ai/d78131f0-3ed1-42ee-9ff4-92299dd3ea96/tasks/output/preview.png",
		"video_url": "https://assets.meshy.ai/d78131f0-3ed1-42ee-9ff4-92299dd3ea96/tasks/output/output.mp4",
		"progress": 100,
		"prompt": "A 3D model of a detailed dog with realistic features.",
		"seed": 1.741027248e+09,
		"started_at": 1.741027287911e+12
	}`

	// Extract the content from the JSON.
	// This function returns the extracted URL (or data), the response format, and the file ID.
	extractedURL, extractedFormat, fileID, fileExtention, err := callers.ContentExtraction(sampleJSON, "model")
	if err != nil {
		log.Fatalf("Extraction failed: %v", err)
	}
	fmt.Println("Extracted URL:", extractedURL)

	// Determine a filename.
	// If fileID is empty, use a temp filename.
	if fileID == "" {
		fileID = "temp"
	}

	// The storage function will detect if the content is a URL and download it if necessary.
	// It returns the content (as bytes) and the local file path.
	_, filePath, err := callers.ContentStorage("model", extractedFormat, fileID, fileExtention, []byte(extractedURL))
	if err != nil {
		log.Fatalf("Storage failed: %v", err)
	}
	fmt.Println("File ID:", filePath)

	// Verify the file was stored.
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("File was not stored at expected location: %s", filePath)
	}
	fmt.Printf("Content successfully stored at: %s\n", filePath)
}
