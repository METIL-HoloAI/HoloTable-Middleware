package callers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
	"github.com/sirupsen/logrus"
)

// ContentExtraction extracts content from the response input based on the data type.
func ContentExtraction(response interface{}, dataType string) (string, string, string, error) {
	// Parse the response into a JSON-compatible structure.
	jsonData, err := parseResponse(response)
	if err != nil {
		logrus.Errorf("Failed to parse response: %v", err)
		return "", "", "", err
	}

	// Retrieve configuration parameters for the given data type.
	configParams, err := getConfigParams(dataType)
	if err != nil {
		logrus.Errorf("Failed to retrieve config params for data type '%s': %v", dataType, err)
		return "", "", "", err
	}

	// Extract the data (URL or raw data string) using the response path.
	dataExtracted, err := extractValueFromData(jsonData, configParams.responsePath)
	if err != nil {
		logrus.Errorf("Failed to extract value from data using path '%s': %v", configParams.responsePath, err)
		return "", configParams.responseFormat, "", err
	}
	println("data:", dataExtracted)

	return dataExtracted, configParams.responseFormat, configParams.fileType, nil
}

// parseResponse parses the response into a JSON-compatible structure.
func parseResponse(response interface{}) (interface{}, error) {
	switch v := response.(type) {
	case string:
		var jsonData interface{}
		if err := json.Unmarshal([]byte(v), &jsonData); err != nil {
			logrus.Errorf("Failed to unmarshal JSON string: %v", err)
			return nil, err
		}
		return jsonData, nil
	default:
		return response, nil
	}
}

// ConfigParams holds configuration parameters for content extraction.
type ConfigParams struct {
	responseFormat string
	responsePath   string
	fileType       string
}

// getConfigParams retrieves configuration parameters for the given data type.
func getConfigParams(dataType string) (*ConfigParams, error) {
	workflow, exists := config.Workflows[dataType]
	if !exists || len(workflow.Steps) == 0 {
		err := fmt.Errorf("unknown or invalid data type: %s", dataType)
		logrus.Error(err)
		return nil, err
	}

	lastStep := workflow.Steps[len(workflow.Steps)-1].ContentExtraction
	return &ConfigParams{
		responseFormat: getStringFromMap(lastStep, "response_format"),
		responsePath:   getStringFromMap(lastStep, "response_path"),
		fileType:       getStringFromMap(lastStep, "file_extention"),
	}, nil
}

// getStringFromMap safely retrieves a string value from a map.
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	logrus.Warnf("Key '%s' not found or not a string in map", key)
	return ""
}

// extractValueFromData traverses the JSON data using the provided JSON path.
func extractValueFromData(data interface{}, responsePath string) (string, error) {
	parts := strings.Split(responsePath, ".")
	current := data

	for _, part := range parts {
		var err error
		current, err = navigateJSON(current, part)
		if err != nil {
			logrus.Errorf("Failed to navigate JSON for part '%s': %v", part, err)
			return "", err
		}
	}
	// Convert the final value to a string if necessary.
	return fmt.Sprintf("%v", current), nil
}

// navigateJSON navigates through the JSON data based on a path segment.
func navigateJSON(data interface{}, path string) (interface{}, error) {
	keys := strings.Split(path, ".") // Split "choices.0.message.content" into ["choices", "0", "message", "content"]
	var current interface{} = data   // Used to keep track of where we are in the JSON structure

	for _, key := range keys {
		switch v := current.(type) {
		case map[string]interface{}: // Check if key exists in the map
			if val, exists := v[key]; exists {
				current = val // If so, update current to the value of the key
			} else {
				err := fmt.Errorf("key '%s' not found in map", key)
				logrus.Error(err)
				return nil, err
			}
		case []interface{}: // Handle array indexing
			index, err := parseIndex(key) // Convert string to int
			if err != nil || index < 0 || index >= len(v) {
				err := fmt.Errorf("invalid array index '%s': %v", key, err)
				logrus.Error(err)
				return nil, err
			}
			current = v[index] // Update current to array element
		default:
			err := fmt.Errorf("unexpected type encountered while navigating JSON at key '%s'", key)
			logrus.Error(err)
			return nil, err
		}
	}
	return current, nil
}
