package main

import (
	"fmt"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/configloader"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/database"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/listeners"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load yaml
	settings, err := configloader.GetGeneral()
	if err != nil {
		fmt.Println("Error loading general settings")
		fmt.Println(err)
		return
	}

	database.Init()

	// Load API keys
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Check how user wants to listen for input
	// and start that listener
	if settings.Listener == "mic" {
		fmt.Println("Microphone Listener")
	} else if settings.Listener == "text" {
		listeners.StartTextListener()
	} else {
		fmt.Println("Invalid listener option in general.yaml")
	}

}
