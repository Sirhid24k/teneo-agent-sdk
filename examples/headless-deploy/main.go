package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TeneoProtocolAI/teneo-agent-sdk/pkg/deploy"
)

func main() {
	// Load configuration from environment variables
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://localhost:8080"
	}

	rpcEndpoint := os.Getenv("RPC_ENDPOINT")
	if rpcEndpoint == "" {
		log.Fatal("RPC_ENDPOINT environment variable is required")
	}

	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		log.Fatal("PRIVATE_KEY environment variable is required")
	}

	// Define agent capabilities
	capabilities := []map[string]interface{}{
		{
			"name":        "text_analysis",
			"description": "Analyze and process text content",
		},
		{
			"name":        "data_processing",
			"description": "Process structured data",
		},
	}

	capabilitiesJSON, err := json.Marshal(capabilities)
	if err != nil {
		log.Fatalf("Failed to marshal capabilities: %v", err)
	}

	// Create deployment configuration
	cfg := deploy.DeployConfig{
		BackendURL:  backendURL,
		RPCEndpoint: rpcEndpoint,
		PrivateKey:  privateKey,

		// Agent configuration
		AgentID:      "my-headless-agent",
		AgentName:    "My Headless Agent",
		Description:  "An agent deployed headlessly via the SDK",
		AgentType:    "command",
		Capabilities: capabilitiesJSON,

		// Optional: custom state file path
		// StateFilePath: "/custom/path/to/state.json",
	}

	// Create deployer
	deployer, err := deploy.NewDeployer(&cfg)
	if err != nil {
		log.Fatalf("Failed to create deployer: %v", err)
	}

	// Execute deployment with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	fmt.Println("Starting headless agent deployment...")
	fmt.Println("==========================================")

	result, err := deployer.Deploy(ctx)
	if err != nil {
		log.Fatalf("Deployment failed: %v", err)
	}

	fmt.Println("\n==========================================")
	fmt.Println("Deployment Result:")
	fmt.Println("==========================================")
	fmt.Printf("  Agent ID:        %s\n", result.AgentID)
	fmt.Printf("  Token ID:        %d\n", result.TokenID)
	fmt.Printf("  Transaction:     %s\n", result.TxHash)
	fmt.Printf("  Contract:        %s\n", result.ContractAddress)
	fmt.Printf("  Metadata URI:    %s\n", result.MetadataURI)
	fmt.Printf("  Already Minted:  %v\n", result.AlreadyMinted)
	if result.DatabaseID != "" {
		fmt.Printf("  Database ID:     %s\n", result.DatabaseID)
	}
	fmt.Println("==========================================")

	// The agent is now:
	// 1. Minted on-chain with the NFT
	// 2. Saved to the backend database
	// 3. Ready to connect via WebSocket

	fmt.Println("\nAgent deployed successfully!")
	fmt.Printf("You can now use Token ID %d to connect to the Teneo network.\n", result.TokenID)
}
