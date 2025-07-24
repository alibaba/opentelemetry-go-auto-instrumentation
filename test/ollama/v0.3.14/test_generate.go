package main

import (
	"context"
	"fmt"
	"log"
	
	"github.com/ollama/ollama/api"
)

func main() {
	// Test Generate API instrumentation
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Printf("Creating default client: %v", err)
		client = &api.Client{} // Create default client for testing
	}
	
	ctx := context.Background()
	req := &api.GenerateRequest{
		Model:  "llama3:8b",
		Prompt: "Hello, world!",
	}
	
	fmt.Println("Testing Generate API instrumentation...")
	
	// Track responses
	var responses []api.GenerateResponse
	
	// This will trigger our instrumentation
	err = client.Generate(ctx, req, func(resp api.GenerateResponse) error {
		responses = append(responses, resp)
		if resp.Done {
			fmt.Printf("Final response: %s\n", resp.Response)
			fmt.Printf("Token counts - Input: %d, Output: %d\n", 
				resp.PromptEvalCount, resp.EvalCount)
		}
		return nil
	})
	
	if err != nil {
		fmt.Printf("Generate error (expected if no server): %v\n", err)
	}
	
	fmt.Println("Generate test completed!")
}