package listeners

import (
	"fmt"

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
		} else {
			// CallIntentDetection(input)
		}

		fmt.Println()
	}
}
