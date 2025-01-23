package callers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"

    "github.com/METIL-HoloAI/HoloTable-Middleware/internal/configloader"
    "github.com/METIL-HoloAI/HoloTable-Middleware/internal/configloader/structs"
)

// LoadPrompt continues the chat with ChatGPT and sends the prompt, then saves and returns the JSON response
func LoadPrompt(prompt string) (map[string]interface{}, error) {
    // Load chat settings
    chatSettings, err := configloader.GetChat()
    if err != nil {
        return nil, fmt.Errorf("error loading chat settings: %w", err)
    }

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
    initPrompt := fmt.Sprintf(`Await for my next response once you read this and with every next text I send do only send json files with no additional explanation or text. Here are a few yaml files. Remember these and you should await for my next response which will have a prompt, use the prompt to decide the contentType the user wants between a gif, image, 3d or video. Then use the appropriate yaml to create a json file according to the information in the yaml. Here is which yaml is connected to which contentType: gif: gifgen.yaml, image:imagegen.yaml, video:videogen.yaml, 3d:3dgen.yaml. Please ignore the following parameters from the yaml and don't include them in the json file: {endpoint: "https://api.openai.com/v1/images/generations" 
method: "POST"
headers:
  Authorization: "Bearer $IMAGEGEN_API_KEY"
  Content-Type: "application/json"
}. If the user did not indicate an intent, return the following default: fakeJSONData := []byte({ 
    "ContentType": "none", 
    "requiredParameters": {
        "prompt": "A futuristic cityscape at sunset"
    },
    "optionalParameters": {
        "model": "dall-e-2",
        "n": 3,
        "quality": "standard",
        "response_format": "url",
        "size": "1024x1024",
        "style": "vivid",
        "user": "user1234"
    }
})
Remember, from now on you will only send me json files.%s`, yamlContents)

    // Create the payload for the chat API
    payload := map[string]interface{}{
        "model": "gpt-4o", // Use the appropriate model
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
    req, err := http.NewRequest(chatSettings.Method, chatSettings.Endpoint, bytes.NewBuffer(payloadBytes))
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    // Add headers
    for key, value := range chatSettings.Headers {
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

    // Save the JSON response to a file
    file, err := os.Create("response.json")
    if err != nil {
        return nil, fmt.Errorf("error creating file: %w", err)
    }
    defer file.Close()

    if _, err := file.Write(body); err != nil {
        return nil, fmt.Errorf("error writing to file: %w", err)
    }

    return jsonResponse, nil
}