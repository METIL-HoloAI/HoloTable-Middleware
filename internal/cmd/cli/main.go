package main

import (
	"fmt"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/configloader"
)

func main() {
	// Load yaml
	settings, err := configloader.GetGeneral()
	if err != nil {
		fmt.Println("Error loading general settings")
		fmt.Println(err)
		return
	}

	// Check how user wants to listen for input
	// and start that listener
	if settings.Listener == "mic" {
		fmt.Println("Microphone Listener")
	} else if settings.Listener == "text" {
		fmt.Println("Text Listener")
	} else {
		fmt.Println("Invalid listener option in general.yaml")
	}

}
