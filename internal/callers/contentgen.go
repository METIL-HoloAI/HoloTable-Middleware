package callers

import (
	"fmt"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/configloader"
)

func LoadJSON() {

	settings, err := configloader.GetImageGeneration()
	if err != nil {
		fmt.Println("Error loading intent detection JSON")
		fmt.Println(err)
		return
	}
}

func LoadYAML() {

	settings, err := configloader.GetImageGeneration()
	if err != nil {
		fmt.Println("Error loading general settings")
		fmt.Println(err)
		return
	}
}
