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
var IntentDetection structs.IntentDetection
var ImageGen structs.APIConfig
var VideoGen structs.APIConfig
var GifGen structs.APIConfig
var ModelGen structs.APIConfig
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
	//contentgen yamls
	ImageGen, err = getImage()
	if err != nil {
		log.Fatal("Error parsing imagegen.yaml: ", err)
	}

	VideoGen, err = getVideo()
	if err != nil {
		log.Fatal("Error parsing videogen.yaml: ", err)
	}

	GifGen, err = getGif()
	if err != nil {
		log.Fatal("Error parsing gifgen.yaml: ", err)
	}

	ModelGen, err = get3d()
	if err != nil {
		log.Fatal("Error parsing 3dgen.yaml: ", err)
	}

	//workflows
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

func getIntentDetection() (structs.IntentDetection, error) {
	configPath, err := getConfigPath("intentdetection.yaml")
	if err != nil {
		return structs.IntentDetection{}, err
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return structs.IntentDetection{}, err
	}

	var settings structs.IntentDetection
	if err := yaml.Unmarshal(file, &settings); err != nil {
		return structs.IntentDetection{}, err
	}

	return settings, nil
}

func getImage() (structs.APIConfig, error) {
	configPath, err := getConfigPath("/contentgen_yamls/imagegen.yaml")
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

func getVideo() (structs.APIConfig, error) {
	configPath, err := getConfigPath("/contentgen_yamls/videogen.yaml")
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

func getGif() (structs.APIConfig, error) {
	configPath, err := getConfigPath("/contentgen_yamls/gifgen.yaml")
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

func get3d() (structs.APIConfig, error) {
	configPath, err := getConfigPath("/contentgen_yamls/3dgen.yaml")
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

	// Find all YAML files in the directory
	files, err := filepath.Glob(filepath.Join(workflowDir, "*.yaml"))
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
		var raw map[string]structs.Workflow // Directly unmarshalling into the new struct
		err = yaml.Unmarshal(yamlData, &raw)
		if err != nil {
			log.Printf("Error parsing YAML file %s: %v", file, err)
			continue
		}

		// Step 2: Store workflows using the top-level key
		for key, workflow := range raw {
			workflows[key] = workflow // Use the YAML key directly
		}
	}

	return workflows, nil
}
