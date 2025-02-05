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

// Loads intent detection response & selects the appropriate workflow
func LoadIntentDetectionResponse(JSONData []byte) {
	// Read JSON data from intent detection
	var intentDetectionResponse structs.IntentDetectionResponse
	if err := json.Unmarshal(JSONData, &intentDetectionResponse); err != nil {
		fmt.Println("Error unmarshalling intent detection response:", err)
		return
	}

	// Lookup workflow based on content type
	workflow, exists := config.Workflows[intentDetectionResponse.ContentType]
	if !exists {
		fmt.Println("Workflow not found for content type:", intentDetectionResponse.ContentType)
		return
	}

	// Pass intent detection response & workflow to HandleWorkflow
	HandleWorkflow(intentDetectionResponse, workflow)
}

// Handles workflow execution dynamically
func HandleWorkflow(intentDetectionResponse structs.IntentDetectionResponse, workflow structs.Workflow) {
	// Global storage for all response values (ensures placeholders are accessible)
	dataStore := make(map[string]interface{})

	// Loop through workflow steps
	for _, step := range workflow.Steps {
		fmt.Printf("\nExecuting Step: %s\n", step.Name)

		// Replace placeholders dynamically before making the API call
		apiURL := deepReplace(step.URL, dataStore).(string)
		fmt.Printf("Updated API URL: %s\n", apiURL)

		// Create API request configuration dynamically
		apiConfig := structs.APIConfig{
			Endpoint: apiURL,
			Method:   step.Method,
			Headers:  step.Headers,
		}

		// Build API payload dynamically with recursive placeholder replacement
		payload := deepReplace(step.Body, dataStore).(map[string]interface{})
		fmt.Printf("Updated Payload: %+v\n", payload)

		// Make API call
		responseData, err := makeAPICall(apiConfig, payload)
		if err != nil {
			fmt.Printf("Error in step '%s': %v\n", step.Name, err)
			return
		}

		// Print full API response for debugging
		fmt.Printf("API Response for '%s': %+v\n", step.Name, responseData)

		// Store response values for future steps using response_placeholders mapping
		for placeholder, responseKey := range step.ResponsePlaceholders {
			if val, exists := responseData[responseKey]; exists {
				dataStore[placeholder] = val
				fmt.Printf("Stored '%s' = %v for future use\n", placeholder, val)
			} else {
				fmt.Printf("Warning: Expected response key '%s' not found in step '%s'\n", responseKey, step.Name)
			}
		}

		// Handle polling if required
		if step.Poll != nil {
			err = pollForCompletion(step, dataStore)
			if err != nil {
				fmt.Printf("polling error in step '%s': %v\n", step.Name, err)
				return
			}
		}
	}

	fmt.Println("Workflow execution completed successfully.")
}

// Makes the API request & returns the response
func makeAPICall(apiConfig structs.APIConfig, payload map[string]interface{}) (map[string]interface{}, error) {
	client := &http.Client{}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest(apiConfig.Method, apiConfig.Endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range apiConfig.Headers {
		req.Header.Set(key, value)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status: %d, body: %s", resp.StatusCode, body)
	}

	// Parse JSON response
	var responseData map[string]interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	return responseData, nil
}

// Handles polling for async workflows
func pollForCompletion(step structs.Step, dataStore map[string]interface{}) error {
	client := &http.Client{}

	for {
		url := deepReplace(step.URL, dataStore).(string)
		fmt.Printf("Polling URL: %s\n", url)

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

		// Print polling response for debugging
		fmt.Printf("Polling Response: %+v\n", responseData)

		// Store response values for later use
		for placeholder, responseKey := range step.ResponsePlaceholders {
			if val, exists := responseData[responseKey]; exists {
				dataStore[placeholder] = val
				fmt.Printf("Updating DataStore: '%s' = %v\n", placeholder, val)
			}
		}

		// Check if polling condition is met
		if responseData[step.Poll.Condition] == step.Poll.TargetValue {
			fmt.Printf("Polling complete for '%s'!\n", step.Name)
			return nil
		}

		// Wait before polling again
		fmt.Printf("Polling '%s'... waiting %d seconds\n", step.Name, step.Poll.Interval)
		time.Sleep(time.Duration(step.Poll.Interval) * time.Second)
	}
}

// Recursively replaces placeholders within maps, lists, and strings
func deepReplace(data interface{}, dataStore map[string]interface{}) interface{} {
	switch v := data.(type) {
	case string:
		// Replace all known placeholders in the string
		for key, value := range dataStore {
			placeholder := "{" + key + "}"
			v = strings.ReplaceAll(v, placeholder, fmt.Sprintf("%v", value))
		}
		return v

	case map[string]interface{}:
		// Recursively replace placeholders in a map
		newMap := make(map[string]interface{})
		for key, value := range v {
			newMap[key] = deepReplace(value, dataStore)
		}
		return newMap

	case []interface{}:
		// Recursively replace placeholders in a list
		newList := make([]interface{}, len(v))
		for i, item := range v {
			newList[i] = deepReplace(item, dataStore)
		}
		return newList

	default:
		return v
	}
}
