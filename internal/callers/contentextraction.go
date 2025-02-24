package callers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
)

// ContentExtraction extracts content from the response based on the data type.
// It returns the extracted value along with an error if something goes wrong.
func ContentExtraction(response string, dataType string) (string, error) {
	var responseFormat, responsePath string

	// Select configuration parameters based on the provided data type.
	switch dataType {
	case "image":
		if val, ok := config.Workflows["image"].Steps[0].Response["response_format"].(string); ok {
			responseFormat = val
		} else {
			return "", errors.New("response_format is nil or not a string for data type: image")
		}
		if val, ok := config.Workflows["image"].Steps[0].Response["response_path"].(string); ok {
			responsePath = val
		} else {
			return "", errors.New("response_path is nil or not a string for data type: image")
		}
	case "video":
		if val, ok := config.Workflows["video"].Steps[0].Response["response_format"].(string); ok {
			responseFormat = val
		} else {
			return "", errors.New("response_format is nil or not a string for data type: video")
		}
		if val, ok := config.Workflows["video"].Steps[0].Response["response_path"].(string); ok {
			responsePath = val
		} else {
			return "", errors.New("response_path is nil or not a string for data type: video")
		}
	case "gif":
		if val, ok := config.Workflows["gif"].Steps[0].Response["response_format"].(string); ok {
			responseFormat = val
		} else {
			return "", errors.New("response_format is nil or not a string for data type: gif")
		}
		if val, ok := config.Workflows["gif"].Steps[0].Response["response_path"].(string); ok {
			responsePath = val
		} else {
			return "", errors.New("response_path is nil or not a string for data type: gif")
		}
	case "3d":
		// Access the last step for the 3D workflow.
		lastStep := len(config.Workflows["3d"].Steps) - 1
		if val, ok := config.Workflows["3d"].Steps[lastStep].Response["response_format"].(string); ok {
			responseFormat = val
		} else {
			return "", errors.New("response_format is nil or not a string for data type: 3d")
		}
		if val, ok := config.Workflows["3d"].Steps[lastStep].Response["response_path"].(string); ok {
			responsePath = val
		} else {
			return "", errors.New("response_path is nil or not a string for data type: 3d")
		}
	default:
		return "", errors.New("unknown data type: " + dataType)
	}

	// If the expected response format is "url", perform JSON extraction.
	if responseFormat == "url" {
		return extractURL(response, responsePath)
	}
	// For other response formats, simply return the full response.
	return response, nil
}

// extractURL extracts the URL from the JSON response using the provided responsePath.
// The responsePath is treated as a JSON path, e.g.: "data.model_urls.glb"
func extractURL(response, responsePath string) (string, error) {
	var jsonData interface{}
	if err := json.Unmarshal([]byte(response), &jsonData); err != nil {
		return "", err
	}

	// If the path starts with "response.", remove it.
	responsePath = strings.TrimPrefix(responsePath, "response.")

	// Split the JSON path into parts.
	parts := strings.Split(responsePath, ".")
	current := jsonData
	var err error
	for _, part := range parts {
		current, err = navigateJSON(current, part)
		if err != nil {
			return "", err
		}
	}

	// Ensure the final value is a string.
	str, ok := current.(string)
	if !ok {
		return "", errors.New("final value is not a string")
	}
	return str, nil
}

// navigateJSON navigates through the JSON data based on a path segment.
// It supports keys and (if no explicit index is provided) automatically takes the first element when encountering an array.
func navigateJSON(current interface{}, part string) (interface{}, error) {
	// If the part does not explicitly start with an array index and current is an array,
	// automatically take the first element.
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
		// Process the key portion before the '[', if any.
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

		// Process each array index in the part.
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

	// Otherwise, part is a simple key. Expect current to be an object.
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
