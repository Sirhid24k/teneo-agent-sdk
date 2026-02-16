package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/TeneoProtocolAI/teneo-agent-sdk/pkg/agent"
	"github.com/joho/godotenv"
)

// StaticHandler implements types.AgentHandler
type StaticHandler struct{}

func (h *StaticHandler) ProcessTask(ctx context.Context, task string) (string, error) {
	return "Thank you for sending a request, the logic for finding Alpha Wallets will soon be implemented, stay tuned.", nil
}

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		log.Fatal("PRIVATE_KEY environment variable is required")
	}

	// Agent Configuration
	agentConfig := agent.DefaultConfig()
	agentConfig.Name = "Smart Wallet Finder V2"
	agentConfig.Description = "When prompted with the smart contract address of a memecoin token, Agent returns a list of profitable wallets (Limit, 5) supports only Solana network"
	agentConfig.Capabilities = []string{"analyze_address"}
	agentConfig.PrivateKey = privateKey

	// Enhanced Agent Config
	enhancedConfig := &agent.EnhancedAgentConfig{
		Config:       agentConfig,
		AgentHandler: &StaticHandler{},
	}

	// NFT Configuration Logic
	nftTokenIDStr := os.Getenv("NFT_TOKEN_ID")
	if nftTokenIDStr != "" {
		var tokenID uint64
		if _, err := fmt.Sscanf(nftTokenIDStr, "%d", &tokenID); err == nil && tokenID > 0 {
			enhancedConfig.TokenID = tokenID
			enhancedConfig.Mint = false
			log.Printf("Using existing NFT Token ID: %d", tokenID)
			// Also update the inner config so it knows about the ID for protocol handler
			agentConfig.NFTTokenID = nftTokenIDStr
		} else {
			// Invalid ID, default to minting
			enhancedConfig.Mint = true
		}
	} else {
		// No ID, default to minting
		enhancedConfig.Mint = true
	}

	myAgent, err := agent.NewEnhancedAgent(enhancedConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting Smart Wallet Finder V2 with Static Handler...")
	if err := myAgent.Run(); err != nil {
		log.Fatal(err)
	}
}