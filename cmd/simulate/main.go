package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/TeneoProtocolAI/teneo-agent-sdk/internal/adapters/chain"
	"github.com/TeneoProtocolAI/teneo-agent-sdk/internal/adapters/price"
	"github.com/TeneoProtocolAI/teneo-agent-sdk/internal/core/domain"
	"github.com/TeneoProtocolAI/teneo-agent-sdk/internal/core/service"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// 1. Initialize same services as main.go
	ethURL := os.Getenv("ALCHEMY_ETHEREUM_URL")
	solURL := os.Getenv("ALCHEMY_SOLANA_URL")

	ethService := chain.NewEthereumService(ethURL)
	solService := chain.NewSolanaService(solURL)
	chains := []domain.ChainService{ethService, solService}
	priceService := price.NewDexScreenerService()
	pnlCalc := service.NewPnLCalculator()

	agentService := service.NewAgentService(chains, priceService, pnlCalc)

	// 2. Prepare the input
	input := domain.AgentInput{
		Chain:        "ethereum",
		TokenAddress: "0xbC56a8efee5871B397Fb06254D12a04546B62924",
		Limit:        5,
	}

	// 3. Execute logic
	log.Printf("Simulating analysis for Token: %s on %s...", input.TokenAddress, input.Chain)
	result, err := agentService.AnalyzeToken(context.Background(), input)
	if err != nil {
		log.Fatalf("Analysis failed: %v", err)
	}

	// 4. Print Output
	output, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(output))
}
