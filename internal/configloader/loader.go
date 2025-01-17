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

func GetImageGeneration() (structs.APIConfig, error) {
	file, err := os.ReadFile("config/contentgeneration/imagegeneration.yaml")
	if err != nil {
		return structs.APIConfig{}, err
	}

	var settings structs.APIConfig
	if err := yaml.Unmarshal(file, &settings); err != nil {
		return structs.APIConfig{}, err
	}

	return settings, nil
}

func GetVideoGeneration() (structs.APIConfig, error) {
	file, err := os.ReadFile("config/contentgeneration/videogeneration.yaml")
	if err != nil {
		return structs.APIConfig{}, err
	}

	var settings structs.APIConfig
	if err := yaml.Unmarshal(file, &settings); err != nil {
		return structs.APIConfig{}, err
	}

	return settings, nil
}

func GetGifGeneration() (structs.APIConfig, error) {
	file, err := os.ReadFile("config/contentgeneration/gifgeneration.yaml")
	if err != nil {
		return structs.APIConfig{}, err
	}

	var settings structs.APIConfig
	if err := yaml.Unmarshal(file, &settings); err != nil {
		return structs.APIConfig{}, err
	}

	return settings, nil
}

func Get3dGeneration() (structs.APIConfig, error) {
	file, err := os.ReadFile("config/contentgeneration/3dgeneration.yaml")
	if err != nil {
		return structs.APIConfig{}, err
	}

	var settings structs.APIConfig
	if err := yaml.Unmarshal(file, &settings); err != nil {
		return structs.APIConfig{}, err
	}

	return settings, nil
}
