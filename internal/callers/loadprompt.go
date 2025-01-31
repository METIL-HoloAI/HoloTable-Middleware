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

// LoadPrompt continues the chat with ChatGPT and sends the prompt, then saves and returns the JSON response
func LoadPrompt(prompt string) ([]byte, error) {

	// Read the contents of the YAML files // agreed in group to leave these as is for now
	yamlFiles := []string{"3dgen.yaml", "gifgen.yaml", "imagegen.yaml", "videogen.yaml"}
	var yamlContents string
	for _, file := range yamlFiles {
		content, err := os.ReadFile(filepath.Join("config/contentgen", file))
		if err != nil {
			return nil, fmt.Errorf("error reading YAML file %s: %w", file, err)
		}
		yamlContents += fmt.Sprintf("\n---\n%s:\n%s", file, content)
	}

	// Create the initialization prompt
	initPrompt := fmt.Sprintf(config.IntentDetection.InitialPrompt, yamlContents)

	// Build the payload using the new function
	payload, err := BuildPayload(initPrompt, prompt)
	if err != nil {
		return nil, fmt.Errorf("error building payload: %w", err)
	}

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
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

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

// MADE FOR more MODULAR PAYLOAD using variables instead of hardcode.
// BuildPayload constructs the payload based on dynamically loaded config
func BuildPayload(initPrompt, userPrompt string) (map[string]interface{}, error) {
	payloadConfig := config.IntentDetection.Payload

	// Generate messages with the prompt replacements
	messages := generateMessages(payloadConfig.Messages, initPrompt, userPrompt)

	// Create payload with the model and messages
	payload := map[string]interface{}{
		"model":    payloadConfig.Model,
		"messages": messages,
	}

	return payload, nil
}

// Helper function to generate messages with dynamic placeholders
func generateMessages(messageTemplates []map[string]string, initPrompt, userPrompt string) []map[string]string {
	messages := make([]map[string]string, len(messageTemplates))
	for i, msgTemplate := range messageTemplates {
		message := make(map[string]string)
		for key, value := range msgTemplate {
			content := value
			content = strings.ReplaceAll(content, "$initialPrompt", initPrompt)
			content = strings.ReplaceAll(content, "$userPrompt", userPrompt)
			message[key] = content
		}
		messages[i] = message
	}
	return messages
}
