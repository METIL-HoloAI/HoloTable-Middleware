package configloader

import (
	"os"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/configloader/structs"
	"gopkg.in/yaml.v3"
)

func GetGeneral() (structs.GeneralSettings, error) {

	file, err := os.ReadFile("config/general.yaml")
	if err != nil {
		return structs.GeneralSettings{}, err
	}

	var settings structs.GeneralSettings
	if err := yaml.Unmarshal(file, &settings); err != nil {
		return structs.GeneralSettings{}, err
	}

	return settings, nil
}

func GetIntentDetection() (structs.IntentDetectionSettings, error) {

	file, err := os.ReadFile("config/intentdetection.yaml")
	if err != nil {
		return structs.IntentDetectionSettings{}, err
	}

	var settings structs.IntentDetectionSettings
	if err := yaml.Unmarshal(file, &settings); err != nil {
		return structs.IntentDetectionSettings{}, err
	}

	return settings, nil
}
