package callers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config/structs"
)

func LoadIntentDetectionResponse(JSONData []byte) {
	//read in JSON data from intent detection
	var intentDetectionResponse structs.IntentDetectionResponse
	if err := json.Unmarshal(JSONData, &intentDetectionResponse); err != nil {
		fmt.Println("Error unmarshalling intent detection response")
		fmt.Println(err)
		return
	}

	//load content gen yaml based off JSON data
	var apiConfig structs.APIConfig
	switch intentDetectionResponse.ContentType {
	case "image":
		apiConfig = config.ImageGen
	case "video":
		apiConfig = config.VideoGen
	case "gif":
		apiConfig = config.GifGen
	case "3d":
		apiConfig = config.ModelGen
	default:
		fmt.Println("Intent detection provided invalid content type")
		return
	}

	BuildAPICall(intentDetectionResponse, apiConfig)
}

func BuildAPICall(intentDetectionResponse structs.IntentDetectionResponse, apiConfig structs.APIConfig) {
	// Build the payload for the API call
	payload := buildPayload(intentDetectionResponse)

	// Make the API call
	response, err := makeAPICall(apiConfig, payload)
	if err != nil {
		fmt.Println("Error making API call")
		fmt.Println(err)
		return
	}

	fmt.Println("API Response:", response)
}

func buildPayload(intentDetectionResponse structs.IntentDetectionResponse) map[string]interface{} {
	// Combine required and optional parameters into a single payload
	payload := make(map[string]interface{})
	for key, value := range intentDetectionResponse.RequiredParameters {
		payload[key] = value
	}
	for key, value := range intentDetectionResponse.OptionalParameters {
		payload[key] = value
	}
	return payload
}

func makeAPICall(apiConfig structs.APIConfig, payload map[string]interface{}) (string, error) {
	//create HTTP client using Go's std library. this client handles network communication
	client := &http.Client{}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest(apiConfig.Method, apiConfig.Endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range apiConfig.Headers {
		req.Header.Set(key, value)
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read and return the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 status: %d, body: %s", resp.StatusCode, body)
	}

	return string(body), nil
}
