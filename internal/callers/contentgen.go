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

// Loads the intent detection response & selects the appropriate workflow
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
	// Storage for previous step results (ensures placeholders are accessible)
	dataStore := make(map[string]interface{})

	// Loop through workflow steps
	for i, step := range workflow.Steps {
		fmt.Printf("\n🔹 Executing Step %d: %s\n", i+1, step.Name)

		// Determine API URL (replace placeholders if it's not the first step)
		var apiURL string
		if i == 0 {
			apiURL = step.URL
		} else {
			apiURL = deepReplace(step.URL, dataStore).(string)
		}
		fmt.Printf("🔄 Updated API URL: %s\n", apiURL)

		// Create API request configuration
		apiConfig := structs.APIConfig{
			Endpoint: apiURL,
			Method:   step.Method,
			Headers:  step.Headers,
		}

		// Build the payload
		var payload map[string]interface{}
		if i == 0 {
			payload = buildPayload(intentDetectionResponse) // First step: Use intent detection parameters
		} else {
			payload = deepReplace(step.Body, dataStore).(map[string]interface{}) // Replace placeholders for later steps
		}
		fmt.Printf("📦 Final Payload for API Call: %+v\n", payload)

		// Make the API call
		responseData, err := makeAPICall(apiConfig, payload)
		if err != nil {
			fmt.Printf("❌ Error in step '%s': %v\n", step.Name, err)
			return
		}

		fmt.Printf("✅ API Response for '%s': %+v\n", step.Name, responseData)

		// **Extract & Store Response Data for Future Steps**
		for placeholder, responseKey := range step.ResponsePlaceholders {
			// Ensure responseKey is a string before using it as a map key
			if responseKeyStr, ok := responseKey.(string); ok {
				if val, exists := responseData[responseKeyStr]; exists {
					dataStore[placeholder] = val
					fmt.Printf("🔑 Stored '%s' = %v for future use\n", placeholder, val)
				} else {
					fmt.Printf("⚠️ Warning: Expected response key '%s' not found in API response for step '%s'\n", responseKeyStr, step.Name)
				}
			} else {
				fmt.Printf("❌ Error: Response key for placeholder '%s' is not a string in step '%s'\n", placeholder, step.Name)
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

	fmt.Println("🎉 Workflow execution completed successfully.")
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

// Recursively replaces placeholders **only for response placeholders**
func deepReplace(data interface{}, dataStore map[string]interface{}) interface{} {
	switch v := data.(type) {
	case string:
		// Replace response placeholders in the string
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

// Handles polling for async workflows
func pollForCompletion(step structs.Step, dataStore map[string]interface{}) error {
	client := &http.Client{}

	// Extract polling configuration dynamically
	conditionKey, ok := step.Poll["condition"].(string)
	if !ok {
		return fmt.Errorf("❌ Error: Polling condition key is missing or not a string in step '%s'", step.Name)
	}

	targetValue, ok := step.Poll["target_value"]
	if !ok {
		return fmt.Errorf("❌ Error: Polling target value is missing in step '%s'", step.Name)
	}

	interval, ok := step.Poll["interval"].(float64) // JSON numbers are parsed as float64
	if !ok {
		return fmt.Errorf("❌ Error: Polling interval is missing or not a number in step '%s'", step.Name)
	}

	// Start polling loop
	for {
		// Replace placeholders in the polling URL
		pollURL := deepReplace(step.URL, dataStore).(string)
		fmt.Printf("🔄 Polling URL: %s\n", pollURL)

		// Make a GET request to check the status
		req, err := http.NewRequest("GET", pollURL, nil)
		if err != nil {
			return fmt.Errorf("❌ Error creating polling request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("❌ Error making polling request: %v", err)
		}
		defer resp.Body.Close()

		// Extract response data
		var responseData map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
			return fmt.Errorf("❌ Error decoding polling response: %v", err)
		}

		// Debugging: Print polling response
		fmt.Printf("📊 Polling Response: %+v\n", responseData)

		// Extract condition value from response
		currentValue, exists := responseData[conditionKey]
		if !exists {
			fmt.Printf("⚠️ Warning: Expected polling key '%s' not found in response\n", conditionKey)
			continue
		}

		// Check if the condition is met
		if currentValue == targetValue {
			fmt.Printf("✅ Polling complete for '%s'!\n", step.Name)
			return nil
		}

		// Wait before polling again
		fmt.Printf("⏳ Polling '%s'... waiting %.0f seconds\n", step.Name, interval)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
