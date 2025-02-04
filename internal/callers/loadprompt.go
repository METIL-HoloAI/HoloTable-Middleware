package callers

// use jsonData, err := callers.LoadPrompt(prompt) to call this function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
)

// LoadPrompt sends the prompt to chat ai, then saves and returns the JSON response
func LoadPrompt(prompt string) ([]byte, error) {

	// Read the contents of the YAML files and concatenate them in a string
	//I want opnions on if I should update this to match Mitchells approach of reading in the entire directory
	yamlFiles := []string{"3dgen.yaml", "gifgen.yaml", "imagegen.yaml", "videogen.yaml"}
	var yamlContents string
	for _, file := range yamlFiles {
		content, err := os.ReadFile(filepath.Join("config/contentgen", file))
		if err != nil {
			return nil, fmt.Errorf("error reading YAML file %s: %w", file, err)
		}
		yamlContents += fmt.Sprintf("\n---\n%s:\n%s", file, content)
	}

	// Build the initial prompt with the concatenated YAML contents
	initPrompt := fmt.Sprintf(config.IntentDetection.InitialPrompt, yamlContents)

	// Build the payload
	payload, err := BuildPayload(initPrompt, prompt)
	if err != nil {
		return nil, fmt.Errorf("error building Intent Detection payload: %w", err)
	}

	//DEBUG
	prettyPrint(payload)
	//DEBUG

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Intent Detection payload: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest(config.IntentDetection.Method, config.IntentDetection.Endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating Intent Detection request: %w", err)
	}

	// Add headers
	for key, value := range config.IntentDetection.Headers {
		req.Header.Set(key, value)
	}

	// Create HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making Intent Detection API call: %w", err)
	}
	defer resp.Body.Close()

	// Read and handle the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Parse the JSON response
	var jsonResponse map[string]interface{}

	//DEBUG
	fmt.Println("Raw Response Body:\n", string(body))
	//DEBUG

	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	//extract the message from the intent detection response
	extractedText := extractByPath(jsonResponse, config.IntentDetection.ResponsePath)
	if extractedText == "" {
		return nil, fmt.Errorf("error extracting response using path: %s", config.IntentDetection.ResponsePath)
	}

	//DEBUG
	fmt.Println("Extracted Content:\n", extractedText)
	//DEBUG

	//Return the json data as a byteslice to be used for content generation api calling
	return []byte(extractedText), nil
}

// ////
func BuildPayload(initPrompt, userPrompt string) (map[string]interface{}, error) {
	payloadConfig := config.IntentDetection.Payload
	payload := deepReplace(payloadConfig, initPrompt, userPrompt).(map[string]interface{})
	return payload, nil
}

// Recursively searches the payload for the initialPronpt and userPrompt placeholders within our intent detection struct and replaces them with the actual values
// can handle any data structure
func deepReplace(data interface{}, initPrompt, userPrompt string) interface{} {
	switch v := data.(type) {
	case string: //if the data is a string, it it searches and replaces the placeholders with the actual values
		return replacePlaceholders(v, initPrompt, userPrompt)
	case map[string]interface{}: //creates a new map to store the modified key-value pairs then traverses the map in search of placeholders
		newMap := make(map[string]interface{})
		for key, value := range v {
			newMap[key] = deepReplace(value, initPrompt, userPrompt)
		}
		return newMap
	case []interface{}: // create a new list to store the modified values then traverses the list in search of placeholders
		newList := make([]interface{}, len(v))
		for i, item := range v {
			newList[i] = deepReplace(item, initPrompt, userPrompt)
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

// Helper function to pretty-print JSON output
func prettyPrint(data interface{}) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}
	fmt.Println("Built Payload:\n", string(jsonBytes))
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
