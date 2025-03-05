package callers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
)

// ContentExtraction extracts content from the JSON response based on the data type.
func ContentExtraction(response string, dataType string) (string, string, string, string, error) {
	var responseFormat, responsePath, fileIDPath, fileType string

	// Select configuration parameters based on the provided data type.
	lastStep := len(config.Workflows[dataType].Steps) - 1
	switch dataType {
	case "image", "video", "gif", "model":
		responseFormat, responsePath, fileIDPath, fileType = getConfigParams(dataType, lastStep)
	default:
		return "", "", "", "", errors.New("unknown data type: " + dataType)
	}

	// Extract the data (URL or raw data string).
	dataExtracted, err := extractValue(response, responsePath)
	if err != nil {
		return "", responseFormat, "", "", err
	}

	// Extract file ID if a file_id_path is provided.
	var fileID string
	if fileIDPath != "" {
		fileID, err = extractValue(response, fileIDPath)
		if err != nil {
			return "", responseFormat, "", "", err
		}
	}

	return dataExtracted, responseFormat, fileID, fileType, nil
}

// getConfigParams retrieves configuration parameters for the given data type and step index.
func getConfigParams(dataType string, stepIndex int) (string, string, string, string) {
	workflow := config.Workflows[dataType].Steps[stepIndex].Response
	return getStringFromMap(workflow, "response_format"),
		getStringFromMap(workflow, "response_path"),
		getStringFromMap(workflow, "file_id_path"),
		getStringFromMap(workflow, "file_extention")
}

// getStringFromMap safely retrieves a string value from a map.
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

// extractValue traverses the JSON response using the provided JSON path (dot-separated)
// and returns the final value as a string. If the final value is not a string, it converts it.
func extractValue(response, responsePath string) (string, error) {
	var jsonData interface{}
	if err := json.Unmarshal([]byte(response), &jsonData); err != nil {
		return "", err
	}

	// Remove the "response." prefix if present.
	responsePath = strings.TrimPrefix(responsePath, "response.")
	parts := strings.Split(responsePath, ".")
	current := jsonData
	var err error
	for _, part := range parts {
		current, err = navigateJSON(current, part)
		if err != nil {
			return "", err
		}
	}

	// If the final value is not a string, convert it to one.
	if str, ok := current.(string); ok {
		return str, nil
	}
	return fmt.Sprintf("%v", current), nil
}

// navigateJSON navigates through the JSON data based on a path segment.
// It supports simple keys and, if no explicit array index is provided,
// automatically selects the first element when encountering an array.
func navigateJSON(current interface{}, part string) (interface{}, error) {
	// If current is an array and the path segment does not start with an explicit index,
	// automatically select the first element.
	if len(part) > 0 && part[0] != '[' {
		if arr, ok := current.([]interface{}); ok {
			if len(arr) == 0 {
				return nil, errors.New("array is empty when processing key: " + part)
			}
			current = arr[0]
		}
	}

	// If the segment contains an explicit array index, process it.
	if idx := strings.Index(part, "["); idx != -1 {
		// Process the key portion before the '[' (if any).
		key := part[:idx]
		if key != "" {
			m, ok := current.(map[string]interface{})
			if !ok {
				return nil, errors.New("expected JSON object for key: " + key)
			}
			var exists bool
			current, exists = m[key]
			if !exists {
				return nil, errors.New("key not found: " + key)
			}
		}

		// Process all array indices in the part.
		for {
			start := strings.Index(part, "[")
			if start == -1 {
				break
			}
			end := strings.Index(part, "]")
			if end == -1 {
				return nil, errors.New("malformed array index in part: " + part)
			}
			indexStr := part[start+1 : end]
			arrIdx, err := strconv.Atoi(indexStr)
			if err != nil {
				return nil, errors.New("invalid array index: " + indexStr)
			}
			arr, ok := current.([]interface{})
			if !ok {
				return nil, errors.New("expected JSON array when processing index: " + indexStr)
			}
			if arrIdx < 0 || arrIdx >= len(arr) {
				return nil, errors.New("array index out of range: " + indexStr)
			}
			current = arr[arrIdx]
			part = part[end+1:]
		}
		return current, nil
	}

	// Otherwise, treat part as a simple key.
	m, ok := current.(map[string]interface{})
	if !ok {
		return nil, errors.New("expected JSON object for key: " + part)
	}
	val, exists := m[part]
	if !exists {
		return nil, errors.New("key not found: " + part)
	}
	return val, nil
}
