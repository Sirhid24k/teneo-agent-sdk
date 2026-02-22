package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/TeneoProtocolAI/teneo-agent-sdk/examples/openai-agent/internal/config"
	"github.com/TeneoProtocolAI/teneo-agent-sdk/examples/openai-agent/internal/engine"
	"github.com/TeneoProtocolAI/teneo-agent-sdk/examples/openai-agent/internal/helius"
	"github.com/TeneoProtocolAI/teneo-agent-sdk/examples/openai-agent/internal/parser"
	"github.com/TeneoProtocolAI/teneo-agent-sdk/examples/openai-agent/internal/ranking"
	"github.com/TeneoProtocolAI/teneo-agent-sdk/examples/openai-agent/internal/validator"
	"github.com/TeneoProtocolAI/teneo-agent-sdk/pkg/agent"
	"github.com/joho/godotenv"
)

// AlphaHandler implements types.AgentHandler.
// It processes "analyze <contract_address> <network> [limit]" commands and
// returns the top profitable wallets by realized PnL from swap activity.
type AlphaHandler struct {
	heliusClient *helius.Client
}

func (h *AlphaHandler) ProcessTask(ctx context.Context, task string) (string, error) {
	// 1. Parse and validate the input command
	req, err := validator.ParseCommand(task)
	if err != nil {
		return fmt.Sprintf("Invalid command: %v\n\nUsage: analyze <contract_address> sol [limit]", err), nil
	}

	log.Printf("🔍 Analyzing token %s on %s (limit: %d)", req.ContractAddress, req.Network, req.Limit)

	// 2. Fetch swap transactions from Helius
	txns, err := h.heliusClient.FetchSwapTransactions(ctx, req.ContractAddress)
	if err != nil {
		return fmt.Sprintf("Error fetching swap data: %v", err), nil
	}

	if len(txns) == 0 {
		return fmt.Sprintf("No swap transactions found for %s", req.ContractAddress), nil
	}

	log.Printf("📊 Processing %d swap transactions...", len(txns))

	// 3. Normalize swap events into buy/sell records
	swaps := parser.NormalizeSwaps(txns, req.ContractAddress)
	if len(swaps) == 0 {
		return fmt.Sprintf("No buy/sell swaps found for token %s", req.ContractAddress), nil
	}

	log.Printf("🔄 Normalized %d buy/sell records", len(swaps))

	// 4. Compute PnL per wallet using FIFO cost basis
	walletPnLs := engine.ComputePnL(swaps)

	// 5. Rank wallets and format output
	ranked := ranking.RankWallets(walletPnLs, req.Limit)
	return ranking.FormatOutput(ranked, req.ContractAddress), nil
}

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		log.Fatal("PRIVATE_KEY environment variable is required")
	}

	// Load Helius config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	heliusClient := helius.NewClient(cfg.HeliusAPIKey, cfg.HeliusBaseURL)

	// Agent Configuration
	agentConfig := agent.DefaultConfig()
	agentConfig.Name = "Alpha Wallet Finder"
	agentConfig.Description = "Analyzes a Solana token contract and returns the top profitable wallets by realized PnL from swap activity. Usage: analyze <contract_address> sol [limit]"
	agentConfig.Capabilities = []string{"analyze_address"}
	agentConfig.PrivateKey = privateKey

	// Enhanced Agent Config
	enhancedConfig := &agent.EnhancedAgentConfig{
		Config:       agentConfig,
		AgentHandler: &AlphaHandler{heliusClient: heliusClient},
	}

	// NFT Configuration Logic
	nftTokenIDStr := os.Getenv("NFT_TOKEN_ID")
	if nftTokenIDStr != "" {
		var tokenID uint64
		if _, err := fmt.Sscanf(nftTokenIDStr, "%d", &tokenID); err == nil && tokenID > 0 {
			enhancedConfig.TokenID = tokenID
			enhancedConfig.Mint = false
			log.Printf("Using existing NFT Token ID: %d", tokenID)
			agentConfig.NFTTokenID = nftTokenIDStr
		} else {
			enhancedConfig.Mint = true
		}
	} else {
		enhancedConfig.Mint = true
	}

	myAgent, err := agent.NewEnhancedAgent(enhancedConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("🚀 Starting Alpha Wallet Finder agent...")
	if err := myAgent.Run(); err != nil {
		log.Fatal(err)
	}
}