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
	// call contentgen.go (this is a shortcut for testing prior to intent detection being completed)
	fakeJSONData := []byte(`{
		"ContentType": "image",
		"endpoint": "https://api.openai.com/v1/images/generations",
		"method": "POST",
		"headers": {
			"Authorization": "Bearer $IMAGEGEN_API_KEY",
			"Content-Type": "application/json"
		},
		"requiredParameters": {
			"prompt": {
				"description": "A detailed description of the desired image, such as 'a futuristic cityscape at sunset.'",
				"options": []
			}
		},
		"optionalParameters": {
			"model": {
				"default": "dall-e-2",
				"description": "The AI model used to generate the image.",
				"options": ["dall-e-2", "dall-e-3"]
			},
			"n": {
				"default": 1,
				"description": "The number of images to generate. Must be between 1 and 10. For DALL-E 3, only 1 image is supported.",
				"options": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
			},
			"quality": {
				"default": "standard",
				"description": "The quality of the image. Use 'hd' for high definition (only supported by DALL-E 3).",
				"options": ["standard", "hd"]
			},
			"response_format": {
				"default": "url",
				"description": "Specifies how the generated image is returned: as a URL or Base64-encoded JSON.",
				"options": ["url", "b64_json"]
			},
			"size": {
				"default": "1024x1024",
				"description": "Dimensions of the image, such as '256x256' or '1024x1024'.",
				"options": ["256x256", "512x512", "1024x1024", "1792x1024", "1024x1792"]
			},
			"style": {
				"default": "vivid",
				"description": "The artistic style of the generated image, such as 'vivid' or 'natural.' Only supported for DALL-E 3.",
				"options": ["vivid", "natural"]
			},
			"user": {
				"description": "A unique identifier for the user making the request, used for tracking and analytics.",
				"options": []
			}
		}
	}`)

	callers.LoadIntentDetectionResponse(fakeJSONData)
}
