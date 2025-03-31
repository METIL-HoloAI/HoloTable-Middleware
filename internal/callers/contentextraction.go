package callers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
)

// ContentExtraction extracts content from the response input based on the data type.
// The response can be either a JSON string or a mapped input (already parsed JSON).
func ContentExtraction(response interface{}, dataType string) (string, string, string, error) {
	var jsonData interface{}
	switch v := response.(type) {
	case string:
		// If the response is a string, assume it is a JSON string and unmarshal it.
		if err := json.Unmarshal([]byte(v), &jsonData); err != nil {
			return "", "", "", err
		}
	default:
		// Otherwise, assume it's already a parsed map/slice.
		jsonData = response
	}

	var responseFormat, responsePath, fileType string

	// Select configuration parameters based on the provided data type.
	lastStep := len(config.Workflows[dataType].Steps) - 1
	switch dataType {
	case "image", "video", "gif", "model":
		responseFormat, responsePath, fileType = getConfigParams(dataType, lastStep)
	default:
		return "", "", "", errors.New("unknown data type: " + dataType)
	}

	// Extract the data (URL or raw data string).
	dataExtracted, err := extractValueFromData(jsonData, responsePath)
	if err != nil {
		return "", responseFormat, "", err
	}

	// Extract file ID if a file_id_path is provided.

	return dataExtracted, responseFormat, fileType, nil
}

// getConfigParams retrieves configuration parameters for the given data type and step index.
func getConfigParams(dataType string, stepIndex int) (string, string, string) {
	workflow := config.Workflows[dataType].Steps[stepIndex].ContentExtraction
	return getStringFromMap(workflow, "response_format"),
		getStringFromMap(workflow, "response_path"),
		getStringFromMap(workflow, "file_extention")
}

// getStringFromMap safely retrieves a string value from a map.
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

// extractValueFromData traverses the JSON data using the provided JSON path (dot-separated)
// and returns the final value as a string. If the final value is not a string, it converts it.
func extractValueFromData(data interface{}, responsePath string) (string, error) {
	// If the input data is already a map and doesn't contain a "response" key,
	// remove the "response." prefix from the responsePath.
	if m, ok := data.(map[string]interface{}); ok {
		if _, exists := m["response"]; !exists {
			responsePath = strings.TrimPrefix(responsePath, "response.")
		}
	}

	// Handle misconfigured response paths.
	// If the responsePath equals "Extracted URL:" (as seen in your error),
	// override it with the expected key path for your mapped input.
	if responsePath == "Extracted URL:" {
		responsePath = "data[0].url"
	}

	parts := strings.Split(responsePath, ".")
	current := data
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
