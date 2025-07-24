package main

import (
	"context"
	"fmt"
	"log"
	
	"github.com/ollama/ollama/api"
)

func main() {
	// Test Chat API instrumentation
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Printf("Creating default client: %v", err)
		client = &api.Client{}
	}
	
	ctx := context.Background()
	req := &api.ChatRequest{
		Model: "llama3:8b",
		Messages: []api.Message{
			{Role: "user", Content: "Hello!"},
		},
	}
	
	fmt.Println("Testing Chat API instrumentation...")
	
	// Track responses
	var responses []api.ChatResponse
	
	// This will trigger our instrumentation
	err = client.Chat(ctx, req, func(resp api.ChatResponse) error {
		responses = append(responses, resp)
		if resp.Done {
			fmt.Printf("Final response: %s\n", resp.Message.Content)
			fmt.Printf("Token counts - Input: %d, Output: %d\n", 
				resp.PromptEvalCount, resp.EvalCount)
		}
		return nil
	})
	
	if err != nil {
		fmt.Printf("Chat error (expected if no server): %v\n", err)
	}
	
	fmt.Println("Chat test completed!")
}