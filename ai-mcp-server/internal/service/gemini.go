package service

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiService handles interactions with the Gemini API
type GeminiService struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

// NewGeminiService creates a new GeminiService instance
func NewGeminiService(apiKey string, modelName string) (*GeminiService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating Gemini client: %w", err)
	}

	model := client.GenerativeModel(modelName)
	return &GeminiService{
		client: client,
		model:  model,
	}, nil
}

// GenerateResponse generates a response from the Gemini model
func (s *GeminiService) GenerateResponse(ctx context.Context, prompt string) (string, error) {
	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("error generating content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	return fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0]), nil
}

// Close closes the Gemini client
func (s *GeminiService) Close() error {
	return s.client.Close()
}
