package main

import (
	"fmt"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/callers"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	prompt := "three images of a dog jumping in a chair" // for testing purposes
	jsonData, err := callers.LoadPrompt(prompt)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	callers.LoadIntentDetectionResponse(jsonData)
}
