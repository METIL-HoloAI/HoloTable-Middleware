package callers

// use jsonData, err := callers.LoadPrompt(prompt) to call this function

//TODD:
//figure out how to store the intent detection prompt best in the yaml file
//figure out how to store the users prompt and send it regardless of the name of the body paramater (prolly similar to env functionality)
//figure out best way to send yaml files with content gen
//make calling work
//handle response
//improve intent detection initial prompting

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

	//MAYBE USEFULL, DO RESEARCH ON THE BEST WAY TO PASS THESE FILES TO GPT
	// Read the contents of the YAML files
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
	initialPrompt := config.IntentDetection.InitialPrompt
	initPrompt := fmt.Sprintf(initialPrompt, yamlContents)

	// Create the payload for the chat API
	model := config.IntentDetection.Body["model"]

	// UPDATE THIS PART
	payload := map[string]interface{}{
		"model": model, // Use the appropriate model
		"messages": []map[string]string{
			{"role": "system", "content": initPrompt},
			{"role": "user", "content": prompt},
		},
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

	//CODE TO HANDLE RESPONSE, MAYBE USEFULL WHO KNOWS
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
	choices := jsonResponse["choices"].([]interface{})
	message := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	content := message["content"].(string)

	// Clean up the content string
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	// Convert the cleaned-up content string to a byte slice
	JSONData := []byte(content)

	return JSONData, nil
}
