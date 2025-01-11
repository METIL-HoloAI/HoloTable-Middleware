package main

import (
	"fmt"
	
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/configloader"
)

func main() {
	// Load yaml
	settings, err := configloader.LoadGeneral()
	if err != nil {
		fmt.Println("Error loading general settings")
		fmt.Println(err)
		return
	}

	// Check how user wants to listen for input
	// and start that listener
	if(settings.Listener == "cli") {
		fmt.Println("CLI Listener")
	} else if(settings.Listener == "text") {
		fmt.Println("Text Listener")
	} else {
		fmt.Println("Invalid listener option in general.yaml")
	}

}
