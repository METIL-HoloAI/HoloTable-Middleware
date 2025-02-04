package listeners

import (
	"fmt"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/callers"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
)

func StartTextListener() {
	for {
		fmt.Println("Please type your prompt or one of the following options:")
		fmt.Println("(q)uit, (r)eload yaml")

		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Println(err)
		}

		if input == "q" {
			return
		}

		if input == "r" {
			config.LoadYaml()
			fmt.Println("Reloaded yaml...")
		} else { // Call intent detection
			jsonData, err := callers.LoadPrompt(input)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			//TODO: update loadprompt.go to NOT return data, instead call LoadIntentDetectionResponse directly
			// Pass JSON data from intent detection to contentget.go for the call
			callers.LoadIntentDetectionResponse(jsonData) //
		}

		fmt.Println()
	}
}
