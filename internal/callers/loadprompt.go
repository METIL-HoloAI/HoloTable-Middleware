package callers

// use jsonData, err := callers.LoadPrompt(prompt) to call this function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
	"github.com/sirupsen/logrus"
)

func StartIntentDetection(input string, numRetrys int) {
	jsonData, err := LoadPrompt(input)
	if err != nil {
		logrus.Warn("Error running intent detection:", err)
		return
	}
	// Pass JSON data from intent detection to contentget.go for the call
	// Also pass user text input through for potential retrying in makeAPICall in contentgen.go
	LoadIntentDetectionResponse(jsonData, input, numRetrys)
}

// LoadPrompt sends the prompt to chat ai, then saves and returns the JSON response
func LoadPrompt(prompt string) ([]byte, error) {
	// Assume config.VideoGen, config.GifGen, etc., are populated
	videoString, err := MarshalandPretty(config.VideoGen)
	if err != nil {
		logrus.WithError(err).Error("\nError marshalling VideoGen config")
		return nil, fmt.Errorf("error marshalling VideoGen config: %w", err)
	}

	gifString, err := MarshalandPretty(config.GifGen)
	if err != nil {
		logrus.WithError(err).Error("\nError marshalling GifGen config")
		return nil, fmt.Errorf("error marshalling GifGen config: %w", err)
	}

	modelString, err := MarshalandPretty(config.ModelGen)
	if err != nil {
		logrus.WithError(err).Error("\nError marshalling ModelGen config")
		return nil, fmt.Errorf("error marshalling ModelGen config: %w", err)
	}

	imageString, err := MarshalandPretty(config.ImageGen)
	if err != nil {
		logrus.WithError(err).Error("\nError marshalling ImageGen config")
		return nil, fmt.Errorf("error marshalling ImageGen config: %w", err)
	}

	// Format as readable YAML-style output
	yamlContents := fmt.Sprintf(
		"video:\n%s\n\ngif:\n%s\n\nmodel:\n%s\n\nimage:\n%s\n",
		videoString, gifString, modelString, imageString,
	)

	// Log final output (structured log)
	//logrus.Trace("\nFinal YAML Output:\n" + yamlContents)

	// Use formatted YAML string in initial prompt
	initPrompt := fmt.Sprintf(config.IntentDetection.InitialPrompt, yamlContents)
	logrus.Trace("\n\nFormatted Initial Prompt:\n", initPrompt) // User-friendly print

	// Build the payload
	payload, err := BuildPayload(initPrompt, prompt)
	if err != nil {
		logrus.WithError(err).Error("\nError building Intent Detection payload:")
		return nil, fmt.Errorf("error building Intent Detection payload: %w", err)
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logrus.WithError(err).Error("\nError marshalling Intent Detection payload:")
		return nil, fmt.Errorf("error marshalling Intent Detection payload: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest(config.IntentDetection.Method, config.IntentDetection.Endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		logrus.WithError(err).Error("\nError creating Intent Detection request:")
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
		logrus.WithError(err).Error("\nError making Intent Detection API call:")
		return nil, fmt.Errorf("error making Intent Detection API call: %w", err)
	}
	defer resp.Body.Close()

	// Read and handle the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("\nError reading response body:")
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Parse the JSON response
	var jsonResponse map[string]interface{}

	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		logrus.WithError(err).Error("\nError unmarshalling response:")
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	// Pretty print jsonResponse
	prettyJSON, err := json.MarshalIndent(jsonResponse, "", "    ")
	if err != nil {
		logrus.WithError(err).Error("\nError formatting JSON:")
		fmt.Printf("Error formatting JSON: %v\n", err)
	}

	logrus.Debug("\nIntent Detection output: \n", string(prettyJSON))

	//extract the message from the intent detection response
	extractedText := ExtractByPath(jsonResponse, config.IntentDetection.ResponsePath)
	if extractedText == "" {
		logrus.Error("Error extracting response using path:", config.IntentDetection.ResponsePath)
		return nil, fmt.Errorf("error extracting response using path: %s", config.IntentDetection.ResponsePath)
	}

	//Return the json data as a byteslice to be used for content generation api calling
	return []byte(extractedText), nil
}
