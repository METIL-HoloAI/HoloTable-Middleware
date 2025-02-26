package callers

// use jsonData, err := callers.LoadPrompt(prompt) to call this function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
)

func StartIntentDetection(input string) {
	jsonData, err := LoadPrompt(input)
	if err != nil {
		fmt.Println("Error running intent detection:", err)
		return
	}
	// Pass JSON data from intent detection to contentget.go for the call
	LoadIntentDetectionResponse(jsonData)
}

// LoadPrompt sends the prompt to chat ai, then saves and returns the JSON response
func LoadPrompt(prompt string) ([]byte, error) {

	videoString, err := json.Marshal(config.VideoGen)
	if err != nil {
		return nil, fmt.Errorf("error marshalling VideoGen config: %w", err)
	}

	gifString, err := json.Marshal(config.GifGen)
	if err != nil {
		return nil, fmt.Errorf("error marshalling GifGen config: %w", err)
	}

	modelString, err := json.Marshal(config.ModelGen)
	if err != nil {
		return nil, fmt.Errorf("error marshalling ModelGen config: %w", err)
	}

	imageString, err := json.Marshal(config.ImageGen)
	if err != nil {
		return nil, fmt.Errorf("error marshalling ImageGen config: %w", err)
	}

	yamlContents := bytes.Join([][]byte{
		[]byte("video: " + string(videoString) + "\n"),
		[]byte("gif: " + string(gifString) + "\n"),
		[]byte("model: " + string(modelString) + "\n"),
		[]byte("image: " + string(imageString) + "\n"),
	}, []byte{})

	// Build the initial prompt with the concatenated YAML contents
	initPrompt := fmt.Sprintf(config.IntentDetection.InitialPrompt, yamlContents)

	// Build the payload
	payload, err := BuildPayload(initPrompt, prompt)
	if err != nil {
		return nil, fmt.Errorf("error building Intent Detection payload: %w", err)
	}

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

	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	//extract the message from the intent detection response
	extractedText := extractByPath(jsonResponse, config.IntentDetection.ResponsePath)
	if extractedText == "" {
		return nil, fmt.Errorf("error extracting response using path: %s", config.IntentDetection.ResponsePath)
	}

	//Return the json data as a byteslice to be used for content generation api calling
	return []byte(extractedText), nil
}
