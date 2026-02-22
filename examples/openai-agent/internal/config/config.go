package config

import (
	"fmt"
	"os"
)

// Config holds application-level configuration loaded from environment variables.
type Config struct {
	HeliusAPIKey  string
	HeliusBaseURL string
}

// Load reads configuration from the environment.
// HELIUS_API_KEY is required; HeliusBaseURL has a sensible default.
func Load() (*Config, error) {
	apiKey := os.Getenv("HELIUS_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("HELIUS_API_KEY environment variable is required")
	}

	baseURL := os.Getenv("HELIUS_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api-mainnet.helius-rpc.com"
	}

	return &Config{
		HeliusAPIKey:  apiKey,
		HeliusBaseURL: baseURL,
	}, nil
}
