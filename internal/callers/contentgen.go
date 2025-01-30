package callers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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
	var workflowKey string
	switch intentDetectionResponse.ContentType {
	case "image":
		workflowKey = "imagegen_api"
	case "video":
		workflowKey = "videogen_api"
	case "gif":
		workflowKey = "gifgen_api"
	case "3d":
		workflowKey = "3dgen_api"
	default:
		fmt.Println("Intent detection provided invalid content type:", intentDetectionResponse.ContentType)
		return
	}

	// Lookup the workflow from Workflows map
	workflow, exists := config.Workflows[workflowKey]
	if !exists {
		fmt.Println("Workflow not found for content type:", intentDetectionResponse.ContentType)
		return
	}

	// Instead of calling BuildAPICall directly, call HandleWorkflow
	HandleWorkflow(intentDetectionResponse, workflow.Workflow)
}

func HandleWorkflow(intentDetectionResponse structs.IntentDetectionResponse, workflow structs.Workflow) {
	// Storage for previous step results (e.g., job_id, task_id)
	dataStore := make(map[string]interface{})

	// Loop through workflow steps
	for _, step := range workflow.Steps {
		// Replace placeholders in URL (e.g., {job_id} → 12345)
		apiURL := replacePlaceholders(step.URL, dataStore)

		// Build API config dynamically
		apiConfig := structs.APIConfig{
			Endpoint: apiURL,
			Method:   step.Method,
			Headers:  map[string]string{}, // Headers can be customized per step if needed
		}

		// Make API call (pass an empty request body if no payload is needed)
		responseData, err := BuildAPICall(intentDetectionResponse, apiConfig)
		if err != nil {
			fmt.Printf("Error in step '%s': %v\n", step.Name, err)
			return
		}

		// Store response values (e.g., "job_id") for future steps
		if step.ResponseKey != "" {
			dataStore[step.ResponseKey] = responseData[step.ResponseKey]
		}

		// Handle polling if required
		if step.Poll != nil {
			err = pollForCompletion(step, dataStore)
			if err != nil {
				fmt.Printf("Polling error in step '%s': %v\n", step.Name, err)
				return
			}
		}
	}

	fmt.Println("Workflow execution completed successfully.")
}

func BuildAPICall(intentDetectionResponse structs.IntentDetectionResponse, apiConfig structs.APIConfig) (map[string]interface{}, error) {
	// Build the payload for the API call
	payload := buildPayload(intentDetectionResponse)

	// Make the API call
	response, err := makeAPICall(apiConfig, payload)
	if err != nil {
		return nil, fmt.Errorf("error making API call: %w", err)
	}

	// Parse the JSON response into a map
	var responseData map[string]interface{}
	if err := json.Unmarshal([]byte(response), &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	return responseData, nil
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

//Helper Functions for HandleWorkflow

// this function replaces the {job_id} at the end of the url with the response form the previous step (that is stored in dataStore)
func replacePlaceholders(url string, dataStore map[string]interface{}) string {
	for key, value := range dataStore {
		placeholder := "{" + key + "}"
		url = strings.ReplaceAll(url, placeholder, fmt.Sprintf("%v", value)) // Convert interface{} to string
	}
	return url
}

// this function handles the polling
func pollForCompletion(step structs.Step, dataStore map[string]interface{}) error {
	client := &http.Client{}

	for {
		url := replacePlaceholders(step.URL, dataStore)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Extract response data
		var responseData map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&responseData)

		// Store response values for later use
		if step.ResponseKey != "" {
			dataStore[step.ResponseKey] = responseData[step.ResponseKey]
		}

		// Check if the polling condition is met
		if responseData[step.ResponseKey] == step.Poll.Until {
			fmt.Println("Polling complete for", step.Name)
			return nil
		}

		// Wait before polling again
		fmt.Printf("Polling '%s'... waiting %d seconds\n", step.Name, step.Poll.Interval)
		time.Sleep(time.Duration(step.Poll.Interval) * time.Second)
	}
}
