package services

import (
	"context"
	"fmt"
	"os"

	genai "github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func InitGeminiClient() *genai.Client {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("GEMINI_API_KEY is missing")
		os.Exit(1)
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		fmt.Println("Failed to create Gemini client:", err)
		os.Exit(1)
	}

	// fmt.Println("Listing available Gemini models...")
	// modelIterator := client.ListModels(ctx)

	// // Iterate over available models
	// for {
	// 	model, err := modelIterator.Next()
	// 	if err != nil {
	// 		break // Exit loop when there are no more models
	// 	}
	// 	fmt.Printf("Available Model: %s\n", model.Name)
	// }

	return client
}
