package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/TeneoProtocolAI/teneo-agent-sdk/internal/adapters/chain"
	"github.com/TeneoProtocolAI/teneo-agent-sdk/internal/adapters/price"
	"github.com/TeneoProtocolAI/teneo-agent-sdk/internal/core/domain"
	"github.com/TeneoProtocolAI/teneo-agent-sdk/internal/core/service"
	"github.com/TeneoProtocolAI/teneo-agent-sdk/pkg/agent"
	"github.com/joho/godotenv"
)

type AlphaWalletFinderAgent struct {
	agentService *service.AgentService
}

func (a *AlphaWalletFinderAgent) ProcessTask(ctx context.Context, task string) (string, error) {
	log.Printf("Processing task: %s", task)

	// Clean input
	task = strings.TrimSpace(task)
	task = strings.TrimPrefix(task, "/")

	// Expected Input: JSON string or simple command?
	// Requirement: "Input: { chain, tokenAddress, limit }"
	// We'll try to parse JSON first. If fails, check for command style.
	
	var input domain.AgentInput
	if err := json.Unmarshal([]byte(task), &input); err != nil {
		// Fallback: Try whitespace separated "ethereum 0x... 5"
		parts := strings.Fields(task)
		if len(parts) >= 2 {
			input.Chain = parts[0]
			input.TokenAddress = parts[1]
			if len(parts) >= 3 {
				fmt.Sscanf(parts[2], "%d", &input.Limit)
			}
			if input.Limit == 0 {
				input.Limit = 10 // default
			}
		} else {
			return "", fmt.Errorf("invalid input format: expected JSON or 'chain address [limit]'")
		}
	}

	// Validate (basic)
	if input.Chain == "" || input.TokenAddress == "" {
		return "", fmt.Errorf("missing chain or token address")
	}

	result, err := a.agentService.AnalyzeToken(ctx, input)
	if err != nil {
		return "", fmt.Errorf("analysis failed: %w", err)
	}

	outputBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(outputBytes), nil
}

func main() {
	_ = godotenv.Load() // Load .env if present

	// Initialize Services
	// Using Alchemy URLs for both if present
	ethURL := os.Getenv("ALCHEMY_ETHEREUM_URL")
	solURL := os.Getenv("ALCHEMY_SOLANA_URL")

	// Initialize chain services with Alchemy RPC URLs
	ethService := chain.NewEthereumService(ethURL)
	solService := chain.NewSolanaService(solURL)
	
	chains := []domain.ChainService{ethService, solService}
	priceService := price.NewDexScreenerService()
	pnlCalc := service.NewPnLCalculator()

	agentService := service.NewAgentService(chains, priceService, pnlCalc)

	// Configure Agent
	config := agent.DefaultConfig()
	config.Name = "Alpha Wallet Finder"
	config.Description = "Finds top performing wallets for a specific token"
	config.Capabilities = []string{"analyze_pnl", "find_alpha"}
	config.Image = "https://teneo.ai/assets/agent-avatar.png" // Added default image
	config.PrivateKey = os.Getenv("PRIVATE_KEY")
	// Handle NFT Token ID and Minting
	nftTokenID := os.Getenv("NFT_TOKEN_ID")
	var tokenID uint64
	mint := nftTokenID == ""

	if !mint {
		fmt.Sscanf(nftTokenID, "%d", &tokenID)
	}
	config.NFTTokenID = nftTokenID
	config.OwnerAddress = os.Getenv("OWNER_ADDRESS")

	// Set Backend URL and RPC Endpoint for NFT operations (auto-minting)
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "https://backend.developer.chatroom.teneo-protocol.ai"
	}
	rpcEndpoint := os.Getenv("RPC_ENDPOINT")
	if rpcEndpoint == "" {
		// Use Peaq RPC as the Teneo NFT contract is on Peaq (Chain ID 3338)
		rpcEndpoint = "https://peaq.api.onfinality.io/public"
	}

	enhancedAgent, err := agent.NewEnhancedAgent(&agent.EnhancedAgentConfig{
		Config:      config,
		Mint:        mint,
		TokenID:     tokenID,
		BackendURL:  backendURL,
		RPCEndpoint: rpcEndpoint,
		AgentHandler: &AlphaWalletFinderAgent{
			agentService: agentService,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting Alpha Wallet Finder Agent...")
	enhancedAgent.Run()
}