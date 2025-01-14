package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)


func callGPTAPI(client *openai.Client, prompt string) string {
	// gpt api call
	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are an assistant that identifies the user's intent as one of the following: video, gif, 3d asset, zoom in, zoom out or image. If the intent is unclear, respond with 'unknown'.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	})

	// if there is an error in the api call return unknown
	if err != nil {
		log.Printf("GPT API call error: %v", err)
		return "unknown"
	}

	// get gpts response
	responseContent := resp.Choices[0].Message.Content

	// if it containts the key words return either the string video, gif 3d assset or image else unknown.
	if strings.Contains(strings.ToLower(responseContent), "video") {
		return "create video"
	} else if strings.Contains(strings.ToLower(responseContent), "gif") {
		return "create gif"
	} else if strings.Contains(strings.ToLower(responseContent), "3d asset") {
		return "create 3d asset"
	} else if strings.Contains(strings.ToLower(responseContent), "image") {
		return "create image"
	} else if strings.Contains(strings.ToLower(responseContent), "zoom in") {
		return "zoom in"
	} else if strings.Contains(strings.ToLower(responseContent), "zoom out") {
		return "zoom out"
	}
	return "unknown"
}

func intentDetection(client *openai.Client, prompt string) string{
	// if propmt is null return string no prompt
	if prompt == "" {
		return "no prompt"
	}

	// call Gpt api for intent detection
	return (callGPTAPI(client, prompt))

}

// temo main for testing
func main() {
	// Use Gpt API KEY
	apiKey := ""
	client := openai.NewClient(apiKey)

	// temp prompt and call intentDetection
	prompt := "zoom into a dog flying"
	intent := intentDetection(client, prompt)

	// print intent that is returned
	fmt.Println("Detected intent:", intent)
}
