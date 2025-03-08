package callers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
		log.Println("Error unmarshalling intent detection response:", err)
		return
	}

	// Lookup workflow based on content type
	workflow, exists := config.Workflows[intentDetectionResponse.ContentType]
	if !exists {
		log.Println("Workflow not found for content type:", intentDetectionResponse.ContentType)
		return
	}

	// Pass intent detection response & workflow to HandleWorkflow
	HandleWorkflow(intentDetectionResponse, workflow)
}

func HandleWorkflow(intentDetectionResponse structs.IntentDetectionResponse, workflow structs.Workflow) {
	// Storage for previous step results (ensures placeholders are accessible)
	dataStore := make(map[string]interface{})

	// Loop through workflow steps
	for i, step := range workflow.Steps {
		log.Printf("\n🔹 Executing Step %d: %s\n", i+1, step.Name)

		// if its the first step store url as is, if not check for and replace placeholders in the URL
		var apiURL string
		if i == 0 {
			apiURL = step.URL
		} else {
			apiURL = deepReplace(step.URL, dataStore).(string)
		}
		log.Printf("🔄 Updated API URL: %s\n", apiURL)

		// put together API request configuration for sending to makeAPICall()
		workflowConfig := structs.APIConfig{
			Endpoint: apiURL,
			Method:   step.Method,
			Headers:  step.Headers,
		}

		// Build the payload if its an intent detection step this payload is given by their code else create based on workflow while replacing placeholders
		var payload map[string]interface{}
		if _, exists := step.Body["intent_detection_step"]; exists {
			payload = buildPayload(intentDetectionResponse) // First step: Use intent detection parameters
		} else {
			payload = deepReplace(step.Body, dataStore).(map[string]interface{}) // Replace placeholders for later steps
		}
		log.Printf("📦 Final Payload for API Call: %+v\n", payload)

		// Make the API call passing what we jsut created above
		responseData, err := makeAPICall(workflowConfig, payload)
		if err != nil {
			log.Printf("❌ Error in step '%s': %v\n", step.Name, err)
			return
		}

		log.Printf("✅ API Response for '%s': %+v\n", step.Name, responseData)

		//TODO
		//if(this is the final step){
		//  Call database function (send response data which is a map[string]interface{}, variable type (.glb in workflow), and the file path bs from the meshy docs that i need to add to workflow)
		//}

		// **Extract & Store Response Data for Future Steps**
		for placeholder, responseKey := range step.Response {
			// Ensure responseKey is a string before using it as a map key
			if responseKeyStr, ok := responseKey.(string); ok {
				if val, exists := responseData[responseKeyStr]; exists {
					dataStore[placeholder] = val
					log.Printf("🔑 Stored '%s' = %v for future use\n", placeholder, val)
				} else {
					log.Printf("⚠️ Warning: Expected response key '%s' not found in API response for step '%s'\n", responseKeyStr, step.Name)
				}
			} else {
				log.Printf("❌ Error: Response key for placeholder '%s' is not a string in step '%s'\n", placeholder, step.Name)
			}
		}

		// Handle polling if required
		if step.Poll != nil {
			log.Printf("🔍 DEBUG: Stored Task ID for Polling: %v\n", dataStore["preview_task_id"])
			err = pollForCompletion(step, dataStore)
			if err != nil {
				log.Printf("polling error in step '%s': %v\n", step.Name, err)
				return
			}
		}
		if i == len(workflow.Steps)-1 {
			extractedURL, extractedFormat, fileID, fileExtention, err := ContentExtraction(responseData, intentDetectionResponse.ContentType)
			if err != nil {
				fmt.Printf("Extraction failed: %v", err)
				return
			}
			fmt.Println("Extracted URL:", extractedURL)

			_, filePath, err := ContentStorage(intentDetectionResponse.ContentType, extractedFormat, fileID, fileExtention, []byte(extractedURL))
			if err != nil {
				fmt.Printf("Storage failed: %v", err)
				return
			}
			fmt.Println("File ID:", filePath)

			fmt.Printf("Content successfully stored at: %s\n", filePath)

			fmt.Println("🎉 Workflow execution completed successfully.")
		}
	}

	log.Printf("🎉 Workflow execution completed successfully.")
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

	// Handle non-200 and non-202 status codes
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("non-200/202 status: %d, body: %s", resp.StatusCode, body)
	}

	// Parse JSON response
	var responseData map[string]interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	// ✅ Handle 202 Accepted: Return response for polling
	if resp.StatusCode == http.StatusAccepted {
		log.Printf("🔄 Received 202 Accepted: Task is processing... Storing response.\n")
		return responseData, nil // Let the caller handle polling
	}

	// ✅ Handle 200 OK: Normal successful response
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

	targetValue, ok := step.Poll["until"].(string)
	if !ok {
		return fmt.Errorf("❌ Error: Polling target value is missing or not a string in step '%s'", step.Name)
	}

	intervalRaw, ok := step.Poll["interval"]
	if !ok {
		return fmt.Errorf("❌ Error: Polling interval is missing in step '%s'", step.Name)
	}

	var interval float64
	switch v := intervalRaw.(type) {
	case float64:
		interval = v
	case int:
		interval = float64(v)
	default:
		return fmt.Errorf("❌ Error: Polling interval is not a valid number (got %T) in step '%s'", v, step.Name)
	}

	log.Printf("🔄 Starting Polling for Step: %s | Target Status: %s | Interval: %.0f seconds\n", step.Name, targetValue, interval)

	// Create a timeout timer of 2.5 minutes
	timeout := time.After(500 * time.Second)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("⏳ Timeout: Polling for step '%s' exceeded 2.5 minutes", step.Name)
		default:
			// Replace placeholders in the polling URL
			pollURL := deepReplace(step.URL, dataStore).(string)
			log.Printf("🔄 Polling URL: %s\n", pollURL)

			req, err := http.NewRequest("GET", pollURL, nil)
			if err != nil {
				return fmt.Errorf("❌ Error creating polling request: %v", err)
			}

			for key, value := range step.Headers {
				req.Header.Set(key, value)
			}

			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("❌ Error making polling request: %v", err)
			}
			defer resp.Body.Close()

			var responseData map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
				return fmt.Errorf("❌ Error decoding polling response: %v", err)
			}

			log.Printf("📊 Polling Response: %+v\n", responseData)

			if msg, exists := responseData["message"]; exists {
				log.Printf("❌ API Error: %v\n", msg)
				return fmt.Errorf("❌ Polling failed: API error '%v'", msg)
			}

			currentStatus, exists := responseData["status"]
			if !exists {
				log.Printf("⚠️ Warning: Expected polling key 'status' not found in response\n")
				log.Printf("⏳ Retrying in %.0f seconds...\n", interval)
				time.Sleep(time.Duration(interval) * time.Second)
				continue
			}

			log.Printf("🔍 Current Status: %v | Target: %s\n", currentStatus, targetValue)

			if currentStatus == targetValue {
				log.Printf("✅ Polling complete! Step '%s' reached status '%s'\n", step.Name, targetValue)
				return nil
			}

			log.Printf("⏳ Status: %v | Retrying in %.0f seconds...\n", currentStatus, interval)
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}
}
