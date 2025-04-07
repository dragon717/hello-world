package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Gemini struct {
		APIKey    string `yaml:"api_key"`
		Model     string `yaml:"model"`
		MaxTokens int    `yaml:"max_tokens"`
	} `yaml:"gemini"`
	Server struct {
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
		Timeout int    `yaml:"timeout"`
	} `yaml:"server"`
	Logging struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
		Output string `yaml:"output"`
	} `yaml:"logging"`
}

// LoadConfig loads the configuration from the specified file
func LoadConfig(configPath string) (*Config, error) {
	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse the YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Validate required fields
	if config.Gemini.APIKey == "" {
		return nil, fmt.Errorf("gemini api_key is required")
	}

	return &config, nil
}
