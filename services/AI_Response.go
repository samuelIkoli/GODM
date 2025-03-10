package services

import (
	"context"
	"encoding/json"
	"fmt"

	genai "github.com/google/generative-ai-go/genai"
)

func GetAIResponse(client *genai.Client, message string) (string, error) {

	ctx := context.Background()
	model := client.GenerativeModel("models/gemini-1.5-flash-latest")

	resp, err := model.GenerateContent(ctx, genai.Text(message))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	// Convert first response part to text
	if textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		return string(textPart), nil
	}

	return "", fmt.Errorf("unexpected response format")
}

func FormatResponse(data interface{}) ([]byte, error) {
	formattedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}
	return formattedJSON, nil
}

func GetGeminiEmbedding(client *genai.Client, text string) ([]float32, error) {

    ctx := context.Background()
    model := client.EmbeddingModel("text-embedding-004") // Use an embedding model

    res, err := model.EmbedContent(ctx, genai.Text(text))
    if err != nil {
        return nil, err
    }

    return res.Embedding.Values, nil
}