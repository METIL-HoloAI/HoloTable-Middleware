package config

import (
	"fmt"
	"log"
	"os"
	"strings"
    "reflect"
	"path/filepath"
	"runtime"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config/structs"
	"gopkg.in/yaml.v3"
)

var General structs.GeneralSettings
var IntentDetection structs.APIConfig
var ImageGen structs.APIConfig
var VideoGen structs.APIConfig
var GifGen structs.APIConfig
var ModelGen structs.APIConfig

func LoadYaml() {
	var err error
	General, err = getGeneral()
	if err != nil {
		log.Fatal("Error parsing general.yaml: ", err)
	}

	IntentDetection, err = getIntentDetection()
	if err != nil {
		log.Fatal("Error parsing chatgen.yaml: ", err)
	}

	ImageGen, err = getImage()
	if err != nil {
		log.Fatal("Error parsing imagegen.yaml: ", err)
	}

	// NOTE: These are commented out for the time being while the yaml files
	// are written to match the struct they are meant to be. As the yaml files
	// are written, uncomment the call, variable corresponding function to correctly load them
	//
	// VideoGen, err = getVideo()
	// if err != nil {
	// 	log.Fatal("Error parsing videogen.yaml: ", err)
	// }
	//
	// GifGen, err = getGif()
	// if err != nil {
	// 	log.Fatal("Error parsing gifgen.yaml: ", err)
	// }
	//
	// ModelGen, err = get3d()
	// if err != nil {
	// 	log.Fatal("Error parsing 3dgen.yaml: ", err)
	// }
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
	configPath, err := getConfigPath("chatgen.yaml")
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

	// Recursive function to process fields
	var processFields func(reflect.Value)
	processFields = func(v reflect.Value) {
		switch v.Kind() {
		case reflect.Ptr:
			if !v.IsNil() {
				processFields(v.Elem()) // Dereference and process
			}
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				processFields(v.Field(i))
			}
		case reflect.Map:
			for _, key := range v.MapKeys() {
				val := v.MapIndex(key)
				if val.Kind() == reflect.String && strings.Contains(val.String(), "$CHATGEN_API_KEY") {
					// Update map value if it contains $
					newValue := os.Getenv("INTENT_DETECTION_API_KEY")

					// Test print
					fmt.Println("newValue:", newValue)
					//

					v.SetMapIndex(key, reflect.ValueOf(strings.ReplaceAll(val.String(), "$CHATGEN_API_KEY", newValue)))

					// Test print
					fmt.Printf("Updated map value %s: %s\n", key, v.MapIndex(key))
					//

				} else {
					processFields(val) // Process nested values
				}
			}
		case reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				processFields(v.Index(i))
			}
		case reflect.String:
			if strings.Contains(v.String(), "$CHATGEN_API_KEY") {
				// Update string if it contains $
				newValue := os.Getenv("INTENT_DETECTION_API_KEY")
				v.SetString(strings.ReplaceAll(v.String(), "$CHATGEN_API_KEY", newValue))
			}
		}
	}

	// Start processing the fields of the settings struct
	processFields(reflect.ValueOf(&settings).Elem())

	return settings, nil
}

func getImage() (structs.APIConfig, error) {
	configPath, err := getConfigPath("/contentgen/imagegen.yaml")
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

	// Recursive function to process fields
	var processFields func(reflect.Value)
	processFields = func(v reflect.Value) {
		switch v.Kind() {
		case reflect.Ptr:
			if !v.IsNil() {
				processFields(v.Elem()) // Dereference and process
			}
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				processFields(v.Field(i))
			}
		case reflect.Map:
			for _, key := range v.MapKeys() {
				val := v.MapIndex(key)
				if val.Kind() == reflect.String && strings.Contains(val.String(), "$IMAGEGEN_API_KEY") {
					// Update map value if it contains $
					newValue := os.Getenv("IMAGE_API_KEY")

					// Test print
					fmt.Println("newValue:", newValue)
					//

					v.SetMapIndex(key, reflect.ValueOf(strings.ReplaceAll(val.String(), "$IMAGEGEN_API_KEY", newValue)))

					// Test print
					fmt.Printf("Updated map value %s: %s\n", key, v.MapIndex(key))
					//

				} else {
					processFields(val) // Process nested values
				}
			}
		case reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				processFields(v.Index(i))
			}
		case reflect.String:
			if strings.Contains(v.String(), "$IMAGEGEN_API_KEY") {
				// Update string if it contains $
				newValue := os.Getenv("IMAGE_API_KEY")
				v.SetString(strings.ReplaceAll(v.String(), "$IMAGEGEN_API_KEY", newValue))
			}
		}
	}

	// Start processing the fields of the settings struct
	processFields(reflect.ValueOf(&settings).Elem())

	return settings, nil
}

	/*reflectSettings := reflect.ValueOf(&settings).Elem()
	fmt.Printf("test2")
	for i := 0; i < reflectSettings.NumField(); i++ {
		fmt.Printf("test3")
		field := reflectSettings.Field(i)
		fmt.Printf("test4")
		if field.Kind() == reflect.String && strings.Contains(field.String(), "$") {
			fmt.Printf("test5")
			newValue := os.Getenv("IMAGE_API_KEY")
			field.SetString("$" + newValue)
			
			// Test
			fmt.Printf("Updated field %s: %s\n", reflectSettings.Type().Field(i).Name, field.String())
		}
	}*/

// func getVideo() (structs.APIConfig, error) {
// 	configPath, err := getConfigPath("/contentgen/videogen.yaml")
// 	if err != nil {
// 		return structs.APIConfig{}, err
// 	}
//
// 	file, err := os.ReadFile(configPath)
// 	if err != nil {
// 		return structs.APIConfig{}, err
// 	}
//
// 	var settings structs.APIConfig
// 	if err := yaml.Unmarshal(file, &settings); err != nil {
// 		return structs.APIConfig{}, err
// 	}
//
// 	return settings, nil
// }
//
// func getGif() (structs.APIConfig, error) {
// 	configPath, err := getConfigPath("/contentgen/gifgen.yaml")
// 	if err != nil {
// 		return structs.APIConfig{}, err
// 	}
//
// 	file, err := os.ReadFile(configPath)
// 	if err != nil {
// 		return structs.APIConfig{}, err
// 	}
//
// 	var settings structs.APIConfig
// 	if err := yaml.Unmarshal(file, &settings); err != nil {
// 		return structs.APIConfig{}, err
// 	}
//
// 	return settings, nil
// }
//
// func get3d() (structs.APIConfig, error) {
// 	configPath, err := getConfigPath("/contentgen/3dgen.yaml")
// 	if err != nil {
// 		return structs.APIConfig{}, err
// 	}
//
// 	file, err := os.ReadFile(configPath)
// 	if err != nil {
// 		return structs.APIConfig{}, err
// 	}
//
// 	var settings structs.APIConfig
// 	if err := yaml.Unmarshal(file, &settings); err != nil {
// 		return structs.APIConfig{}, err
// 	}
//
// 	return settings, nil
// }
