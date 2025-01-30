package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config/structs"
	"gopkg.in/yaml.v3"
)

var General structs.GeneralSettings
var IntentDetection structs.APIConfig
var Workflows structs.WorkflowCollection

func LoadYaml() {
	var err error
	General, err = getGeneral()
	if err != nil {
		log.Fatal("Error parsing general.yaml: ", err)
	}

	IntentDetection, err = getIntentDetection()
	if err != nil {
		log.Fatal("Error parsing intentdetection.yaml: ", err)
	}

	Workflows, err = loadWorkflowsFromDir()
	if err != nil {
		log.Fatalf("Error loading workflows: %v", err)
	}

	loadEnv()
}

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

func getGeneral() (structs.GeneralSettings, error) {
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

func getIntentDetection() (structs.APIConfig, error) {
	configPath, err := getConfigPath("intentdetection.yaml")
	if err != nil {
		return structs.APIConfig{}, err
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return structs.APIConfig{}, err
	}

	var settings structs.APIConfig
	if err := yaml.Unmarshal(file, &settings); err != nil {
		return structs.APIConfig{}, err
	}

	return settings, nil
}

func loadWorkflowsFromDir() (structs.WorkflowCollection, error) {
	workflows := make(structs.WorkflowCollection)

	workflowDir, err := getConfigPath("contentgen_workflows")
	if err != nil {
		return nil, fmt.Errorf("error getting workflow config path: %v", err)
	}

	// Ensure the directory exists
	files, err := filepath.Glob(filepath.Join(workflowDir, "*.yaml")) // Find all YAML files
	if err != nil {
		return nil, fmt.Errorf("error finding YAML files in workflow directory: %v", err)
	}

	// Process each YAML file
	for _, file := range files {
		yamlData, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Error reading file %s: %v", file, err)
			continue
		}

		// Step 1: Load YAML into a map (top-level key is retained)
		var raw map[string]structs.APIYaml
		err = yaml.Unmarshal(yamlData, &raw)
		if err != nil {
			log.Printf("Error parsing YAML file %s: %v", file, err)
			continue
		}

		// Step 2: Store workflows using the top-level key
		for key, config := range raw {
			workflows[key] = config // Use the YAML key directly
		}
	}

	return workflows, nil
}
