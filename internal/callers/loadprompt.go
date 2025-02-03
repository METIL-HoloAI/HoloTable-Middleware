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

	payload, err := BuildPayload(initPrompt, prompt)
	if err != nil {
		return nil, fmt.Errorf("error building payload: %w", err)
	}

	prettyPrint(payload)

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest(config.IntentDetection.Method, config.IntentDetection.Endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
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
		return nil, fmt.Errorf("error making API call: %w", err)
	}
	defer resp.Body.Close()

	// Read and handle the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Parse the JSON response
	var jsonResponse map[string]interface{}
	fmt.Println("Raw Response Body:\n", string(body))
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	jsonFormatted, _ := json.MarshalIndent(jsonResponse, "", "  ")
	fmt.Println("Parsed JSON Response:\n", string(jsonFormatted))

	// Extract the content from the assistant's message
	choices, ok := jsonResponse["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return nil, fmt.Errorf("invalid or missing choices in response")
	}
	message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid message format in response")
	}
	content := message["content"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid content format in response")
	}

	// Convert the cleaned-up content string to a byte slice
	JSONData := []byte(content)

	return JSONData, nil
}

// BuildPayload now correctly follows your function signature
func BuildPayload(initPrompt, userPrompt string) (map[string]interface{}, error) {
	payloadConfig := config.IntentDetection.Payload
	payload := deepReplace(payloadConfig, initPrompt, userPrompt).(map[string]interface{})
	return payload, nil
}

// Recursively replace placeholders in ANY structure
func deepReplace(data interface{}, initPrompt, userPrompt string) interface{} {
	switch v := data.(type) {
	case string:
		return replacePlaceholders(v, initPrompt, userPrompt)
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for key, value := range v {
			newMap[key] = deepReplace(value, initPrompt, userPrompt)
		}
		return newMap
	case []interface{}:
		newList := make([]interface{}, len(v))
		for i, item := range v {
			newList[i] = deepReplace(item, initPrompt, userPrompt)
		}
		return newList
	default:
		return v
	}
}

// Helper function to replace placeholders in a string
func replacePlaceholders(text, initPrompt, userPrompt string) string {
	text = strings.ReplaceAll(text, "$initialPrompt", initPrompt)
	text = strings.ReplaceAll(text, "$userPrompt", userPrompt)
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
