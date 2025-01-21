package listeners

import (
	"fmt"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/callers"
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
			fmt.Println("TODO")
		}

		//fmt.Print("\033[H\033[2J")

		callers.LoadIntentDetectionResponse()
	}
}
