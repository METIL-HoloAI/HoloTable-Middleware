package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/callers"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/configloader"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/listeners"
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

	// Load database
	db, err := sql.Open("sqlite3", settings.DataDir+"/db.db")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Check how user wants to listen for input
	// and start that listener
	if settings.Listener == "mic" {
		fmt.Println("Microphone Listener")
	} else if settings.Listener == "text" {
		listeners.StartTextListener()
	} else {
		fmt.Println("Invalid listener option in general.yaml")
	}

	// callers.InitializeChat()

	prompt := "A futuristic cityscape at sunset"
	jsonData, err := callers.LoadPrompt(prompt) // for testing purposes
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Generated JSON Data:", string(jsonData))

	// // call contentgen.go (this is a shortcut for testing prior to intent detection being completed)
	// fakeJSONData := []byte(`{
	// 	"ContentType": "image",
	// 	"requiredParameters": {
	// 		"prompt": "A futuristic cityscape at sunset"
	// 	},
	// 	"optionalParameters": {
	// 		"model": "dall-e-2",
	// 		"n": 3,
	// 		"quality": "standard",
	// 		"response_format": "url",
	// 		"size": "1024x1024",
	// 		"style": "vivid",
	// 		"user": "user1234"
	// 	}
	// }`)

	// callers.LoadIntentDetectionResponse(fakeJSONData) // testing
	// // callers.LoadIntentDetectionResponse(jsonData) // actual
}
