package main

import (
	"fmt"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/callers"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/database"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/listeners"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config.LoadYaml()

	// here I'm calling StartTextListener from listeners/text.go
	// listeners.StartTextListener()

	//TESTING START
	prompt := "Create one image of a dog jumping in a chair" // for testing purposes
	jsonData, err := callers.LoadPrompt(prompt)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonData))

	callers.LoadIntentDetectionResponse(jsonData)
	//TESTING END

	database.Init()

	// Check how user wants to listen for input
	// and start that listener
	if config.General.Listener == "mic" {
		fmt.Println("Microphone Listener")
	} else if config.General.Listener == "text" {
		listeners.StartTextListener()
	} else {
		fmt.Println("Invalid listener option in general.yaml")
	}

}
