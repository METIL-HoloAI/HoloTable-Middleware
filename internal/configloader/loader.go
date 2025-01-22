package configloader

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/configloader/structs"
	"gopkg.in/yaml.v3"
)

// This is a workaround to make sure filepaths are always pulled relative
// to loader.go
func getConfigPath(relativePath string) (string, error) {
	// the caller will return back the file that called this function
	// the way this is set up is that it'll always be a function within
	// loader.go, making sure that filepaths remain consistent
	_, loaderfile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to retrieve caller information")
	}

	// Get the directory where loader.go is located
	loaderDir := filepath.Dir(loaderfile)

	configPath := filepath.Join(loaderDir, "../../config", relativePath)
	configPath = filepath.Clean(configPath)

	return configPath, nil
}

func GetGeneral() (structs.GeneralSettings, error) {
	configPath, err := getConfigPath("general.yaml")
	if err != nil {
		return structs.GeneralSettings{}, err
	}

	file, err := os.ReadFile(configPath)
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
	configPath, err := getConfigPath("intentdetection.yaml")
	if err != nil {
		return structs.IntentDetectionSettings{}, err
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return structs.IntentDetectionSettings{}, err
	}

	var settings structs.IntentDetectionSettings
	if err := yaml.Unmarshal(file, &settings); err != nil {
		return structs.IntentDetectionSettings{}, err
	}

	return settings, nil
}

