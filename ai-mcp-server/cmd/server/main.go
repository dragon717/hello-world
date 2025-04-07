package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bugod/ai-mcp-server/internal/config"
	"github.com/bugod/ai-mcp-server/internal/service"
)

func main() {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// Go up to project root (from cmd/server to project root)
	projectRoot := filepath.Join(wd, "..", "..")

	// Load configuration
	configPath := filepath.Join(projectRoot, "config.yaml")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create Gemini service
	geminiService, err := service.NewGeminiService(cfg.Gemini.APIKey, cfg.Gemini.Model)
	if err != nil {
		log.Fatalf("Failed to create Gemini service: %v", err)
	}
	defer geminiService.Close()

	// Test the service with a simple prompt
	ctx := context.Background()
	prompt := "Hello! Please introduce yourself."

	response, err := geminiService.GenerateResponse(ctx, prompt)
	if err != nil {
		log.Fatalf("Failed to generate response: %v", err)
	}

	fmt.Printf("Prompt: %s\n", prompt)
	fmt.Printf("Response: %s\n", response)
}
