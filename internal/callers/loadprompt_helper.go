package callers

import (
	"fmt"
	"strings"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
)

// Build the payload for the intent detection API call
func BuildPayload(initPrompt, userPrompt string) (map[string]interface{}, error) {
	payloadConfig := config.IntentDetection.Payload
	payload := searchPayload(payloadConfig, initPrompt, userPrompt).(map[string]interface{})
	return payload, nil
}

// Recursively searches the payload for the initialPronpt and userPrompt placeholders within our intent detection struct and replaces them with the actual values
// can handle any data structure
func searchPayload(data interface{}, initPrompt, userPrompt string) interface{} {
	switch v := data.(type) {
	case string: //if the data is a string, it it searches and replaces the placeholders with the actual values
		return replacePlaceholders(v, initPrompt, userPrompt)
	case map[string]interface{}: //creates a new map to store the modified key-value pairs then traverses the map in search of placeholders
		newMap := make(map[string]interface{})
		for key, value := range v {
			newMap[key] = searchPayload(value, initPrompt, userPrompt)
		}
		return newMap
	case []interface{}: // create a new list to store the modified values then traverses the list in search of placeholders
		newList := make([]interface{}, len(v))
		for i, item := range v {
			newList[i] = searchPayload(item, initPrompt, userPrompt)
		}
		return newList
	default: //if the data is not a string, map, or list, it is returned as is
		return v
	}
}

// Replace placeholders in the text with the actual values
func replacePlaceholders(text, initPrompt, userPrompt string) string {
	text = strings.ReplaceAll(text, "%initialPrompt", initPrompt)
	text = strings.ReplaceAll(text, "%userPrompt", userPrompt)
	return text
}

// Extract intent detections response from the api return using the provided response path
func extractByPath(data interface{}, path string) string {
	keys := strings.Split(path, ".") // Split "choices.0.message.content" into ["choices", "0", "message", "content"]
	var current interface{} = data   //used to keep track of where we are in the JSON structure

	for _, key := range keys {
		switch v := current.(type) {
		case map[string]interface{}: // check if key exists in the map
			if val, exists := v[key]; exists {
				current = val //if so, update current to the value of the key
			} else {
				return "" // if not, return an empty string
			}
		case []interface{}:
			index, err := parseIndex(key) //convert string to int
			if err != nil || index < 0 || index >= len(v) {
				fmt.Println("Error parsing index of intent detection response:", err)
				return "" // return empty string if index is invalid
			}
			current = v[index] //update current to array element
		default:
			return "" // Unexpected type
		}
	}

	if str, ok := current.(string); ok {
		return str // return the intent detection response
	}
	return "" // uh oh! No valid text found :(
}

// Convert string to int
func parseIndex(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}
