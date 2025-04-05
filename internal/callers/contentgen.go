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
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/unityserver"
	"github.com/sirupsen/logrus"
)

// Loads the intent detection response & selects the appropriate workflow
func LoadIntentDetectionResponse(JSONData []byte, originalInput string, numRetrys int) {
	// Read JSON data from intent detection
	var intentDetectionResponse structs.IntentDetectionResponse
	if err := json.Unmarshal(JSONData, &intentDetectionResponse); err != nil {
		logrus.Error("\nError unmarshalling intent detection response:", err)
		return
	}

	// Lookup workflow based on content type
	workflow, exists := config.Workflows[intentDetectionResponse.ContentType]
	if !exists {
		logrus.Error("Workflow not found for content type:", intentDetectionResponse.ContentType)
		return
	}

	// Pass intent detection response & workflow to HandleWorkflow
	HandleWorkflow(intentDetectionResponse, workflow, originalInput, numRetrys)
}

func HandleWorkflow(intentDetectionResponse structs.IntentDetectionResponse, workflow structs.Workflow, originalInput string, numRetrys int) {
	// Storage for previous step results (ensures placeholders are accessible)
	dataStore := make(map[string]interface{})

	// Loop through workflow steps
	for i, step := range workflow.Steps {
		logrus.Debugf("\n🔹 Executing Step %d: %s\n", i+1, step.Name)

		// if its the first step store url as is, if not check for and replace placeholders in the URL
		var apiURL string
		if i == 0 {
			apiURL = step.URL
		} else {
			apiURL = deepReplace(step.URL, dataStore).(string)
		}
		logrus.Debugf("\n🔄 Updated API URL: %s\n\n", apiURL)

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
		logrus.Debugf("\n📦 API Call Configuration for Step '%s': \nURL: %s,\n Method: %s,\n Headers: %+s\n\n", step.Name, workflowConfig.Endpoint, workflowConfig.Method, workflowConfig.Headers)
		logrus.Debugf("\n📦 API Call Payload for Step '%s': %+s\n\n", step.Name, PrettyPrintJSON(payload))

		// Make the API call passing what we jsut created above
		var responseData map[string]interface{}
		responseData, err := makeAPICall(workflowConfig, payload, originalInput, numRetrys)
		if err != nil {
			logrus.Errorf("\n❌ Error in step '%s': %v\n", step.Name, err)
			return
		}

		logrus.Debugf("\n✅ API Response for '%s': %+s\n\n", step.Name, PrettyPrintJSON(responseData))

		//TODO
		//if(this is the final step){
		//  Call database function (send response data which is a map[string]interface{}, variable type (.glb in workflow), and the file path bs from the meshy docs that i need to add to workflow)
		//}

		// **Extract & Store Response Data for Future Steps**
		for placeholder, responseKey := range step.ResponsePlaceholders {
			if responseKeyStr, ok := responseKey.(string); ok && step.Poll == nil {
				// Use ExtractByPath to extract nested values dynamically
				extractedValue := ExtractByPath(responseData, responseKeyStr)
				if extractedValue != "" && i != len(workflow.Steps)-1 {
					dataStore[placeholder] = extractedValue
					logrus.Tracef("🔑 Stored '%s' = %v for future use\n", placeholder, extractedValue)
				} else {
					logrus.Warnf("\n⚠️ Warning: Expected response key '%s' not found in API response for step '%s'\n", responseKeyStr, step.Name)
				}
			} else if step.Poll == nil {
				logrus.Errorf("\n❌ Error: Response key for placeholder '%s' is not a string in step '%s'\n", placeholder, step.Name)
			}
		}

		// Handle polling if required
		if step.Poll != nil {
			logrus.Debugf("\n🔍 Stored Task ID for Polling: %v\n", dataStore["preview_task_id"])
			err = pollForCompletion(step, dataStore, responseData)
			if err != nil {
				logrus.Errorf("\npolling error in step '%s': %v\n", step.Name, err)
				return
			}
		}

		if i == len(workflow.Steps)-1 {
			extractedURL, extractedFormat, fileExtention, err := ContentExtraction(responseData, intentDetectionResponse.ContentType)
			if err != nil {
				fmt.Printf("Extraction failed: %v", err)
				return
			}
			//fmt.Println("Extracted URL:", extractedURL)

			dataBytes, filePath, fileID, err := ContentStorage(intentDetectionResponse.ContentType, extractedFormat, fileExtention, []byte(extractedURL))

			if err != nil {
				fmt.Printf("Storage failed: %v", err)
				return
			}
			logrus.Tracef("Content successfully stored at: %s\n", filePath)
			logrus.Debugf("🎉 Workflow execution completed successfully.")

			if unityserver.IsUsingFilepath {
				logrus.Debug("About to send file paths for asset export.")
				unityserver.ExportAssetFile(fileID, fileExtention, filePath)
			} else {
				logrus.Debug("About to send raw data for asset export.")
				unityserver.ExportAssetData(fileID, fileExtention, dataBytes)
			}
		}
	}
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
func makeAPICall(apiConfig structs.APIConfig, payload map[string]interface{}, originalInput string, numRetrys int) (map[string]interface{}, error) {
	client := &http.Client{}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logrus.WithError(err).Error("\nFailed to marshal payload:\n")
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest(apiConfig.Method, apiConfig.Endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		logrus.WithError(err).Error("\nFailed to create request:")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range apiConfig.Headers {
		req.Header.Set(key, value)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("\nFailed to make request:")
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("\nFailed to read response body:")
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// // Handle non-200 and non-202 status codes
	// if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
	// 	logrus.WithError(err).Errorf("\nnon-200/202 status: %d, body: %s", resp.StatusCode, body)
	// 	return nil, fmt.Errorf("non-200/202 status: %d, body: %s", resp.StatusCode, body)
	// }

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated {
		logrus.Errorf("\n❌ Error in API call: non-200-202 status: %d", resp.StatusCode)

		// Try to pretty-print the response if it's in JSON format
		var jsonResponse map[string]interface{}
		if err := json.Unmarshal(body, &jsonResponse); err == nil {
			logrus.Errorf("\n🔍 API Response:\n%s\n", PrettyPrintJSON(jsonResponse))
			if numRetrys < config.General.NumIntentDetectionRetries {
				// Retry intent detection
				logrus.Warn("🔄 Unnecessary data detected. Retrying intent detection...")
				numRetrys++
				go StartIntentDetection(originalInput, numRetrys)
				return nil, fmt.Errorf("retrying due to unnecessary data")
			} else {
				logrus.Error("\n❌ Max retries reached. Unable to process the request. To Increase the number of retries, please update the general.yaml file.")
				return nil, fmt.Errorf("max retries reached")
			}
		} else {
			// If response is not JSON, print it as a raw string
			logrus.Errorf("\n🔍 API Raw Response:\n%s\n", string(body))
		}

		return nil, fmt.Errorf("non-200-202 status: %d", resp.StatusCode)
	}

	// Parse JSON response
	var responseData map[string]interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		logrus.WithError(err).Error("Failed to unmarshal API response:")
		return nil, fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	// ✅ Handle 202 Accepted: Return response for polling
	if resp.StatusCode == http.StatusAccepted {
		logrus.Info("🔄 Received 202 Accepted: Task is processing... Storing response.\n")
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
func pollForCompletion(step structs.Step, dataStore map[string]interface{}, responseData map[string]interface{}) error {
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

	logrus.Debugf("🔄 Starting Polling for Step: %s | Target Status: %s | Interval: %.0f seconds\n", step.Name, targetValue, interval)

	// Create a timeout timer of 2.5 minutes
	timeout := time.After(500 * time.Second)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("⏳ Timeout: Polling for step '%s' exceeded 2.5 minutes", step.Name)
		default:
			// Replace placeholders in the polling URL
			pollURL := deepReplace(step.URL, dataStore).(string)
			logrus.Debugf("🔄 Polling URL: %s\n", pollURL)

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

			// var responseData map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
				return fmt.Errorf("❌ Error decoding polling response: %v", err)
			}

			prettyJSON, _ := json.MarshalIndent(responseData, "", "    ")
			logrus.Debugf("📊 Polling Response: %+s\n", string(prettyJSON))

			if msg, exists := responseData["message"]; exists {
				log.Printf("❌ API Error: %v\n", msg)
				return fmt.Errorf("❌ Polling failed: API error '%v'", msg)
			}

			// 🔍 Extract and check polling status dynamically
			var currentStatus string
			for placeholder, responseKey := range step.ResponsePlaceholders {
				if placeholder == "status" { // Only extract "status" for polling
					if responseKeyStr, ok := responseKey.(string); ok {
						extractedValue := ExtractByPath(responseData, responseKeyStr)
						if extractedValue != "" {
							currentStatus = extractedValue
							logrus.Debugf("🔑 Extracted status: '%s' = %v for polling comparison\n", placeholder, extractedValue)
						} else {
							logrus.Warnf("⚠️ Warning: Expected response key '%s' not found in API response for step '%s'\n", responseKeyStr, step.Name)
						}
					} else {
						logrus.Warnf("❌ Error: Response key for placeholder '%s' is not a string in step '%s'\n", placeholder, step.Name)
					}
				}
			}

			// If no status was extracted, retry
			if currentStatus == "" {
				logrus.Warnf("⚠️ Warning: No status extracted, retrying in %.0f seconds...\n", interval)
				time.Sleep(time.Duration(interval) * time.Second)
				continue
			}

			logrus.Debugf("🔍 Current Status: %v | Target: %s\n", currentStatus, targetValue)

			// ✅ If the status matches the "until" condition, extract other placeholders
			if currentStatus == targetValue {
				logrus.Debugf("✅ Polling complete! Step '%s' reached status '%s'\n", step.Name, targetValue)

				// Extract and store additional placeholders (like image_id)
				for placeholder, responseKey := range step.ResponsePlaceholders {
					if placeholder != "status" { // Ignore status, we already checked it
						if responseKeyStr, ok := responseKey.(string); ok {
							extractedValue := ExtractByPath(responseData, responseKeyStr)
							if extractedValue != "" {
								dataStore[placeholder] = extractedValue
								logrus.Debugf("in polling step 🔑 Stored '%s' = %v for future use\n", placeholder, extractedValue)
							} else {
								logrus.Warnf("⚠️ Warning: Expected response key '%s' not found in API response for step '%s'\n", responseKeyStr, step.Name)
							}
						} else {
							logrus.Errorf("❌ Error: Response key for placeholder '%s' is not a string in step '%s'\n", placeholder, step.Name)
						}
					}
				}

				return nil // Exit polling loop
			}

			logrus.Debugf("⏳ Status: %v | Retrying in %.0f seconds...\n", currentStatus, interval)
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}
}

func PrettyPrintJSON(data map[string]interface{}) string {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false) // Prevents escaping of slashes and other HTML characters
	encoder.SetIndent("", "  ")  // 2-space indentation

	if err := encoder.Encode(data); err != nil {
		logrus.WithError(err).Error("Failed to pretty print JSON payload")
		return "{}"
	}
	return buf.String()
}
