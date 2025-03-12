package listeners

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/callers"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
	"github.com/sirupsen/logrus"
)

func StartTextListener() {
	reader := bufio.NewReader(os.Stdin) // Create a buffered reader

	for {
		fmt.Println("Please type your prompt or one of the following options:")
		fmt.Println("(q)uit, (r)eload yaml")

		input, err := reader.ReadString('\n') // Read full line, including spaces
		if err != nil {
			logrus.Warn("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input) // Remove any trailing newline or spaces

		if input == "q" {
			return
		}

		if input == "r" {
			config.LoadYaml()
			fmt.Println("Reloaded yaml...")
		} else { // Call intent detection
			go callers.StartIntentDetection(input)
		}

		fmt.Println()
	}
}
