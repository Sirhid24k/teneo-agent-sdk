package chain

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/TeneoProtocolAI/teneo-agent-sdk/internal/core/domain"
)

type SolanaService struct {
	rpcURL string
	client *http.Client
}

func NewSolanaService(rpcURL string) *SolanaService {
	return &SolanaService{
		rpcURL: rpcURL,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *SolanaService) IsSupported(chain string) bool {
	return chain == "solana" || chain == "sol"
}

func (s *SolanaService) GetTokenMetadata(ctx context.Context, tokenAddress string) (*domain.TokenMetadata, error) {
	if s.rpcURL == "" {
		return &domain.TokenMetadata{
			Symbol:   "MOCK-SOL",
			Decimals: 9,
			Name:     "Mock Sol Token",
		}, nil
	}
	// Real Alchemy/Solana RPC call would go here
	// e.g. "getTokenSupply" or "getAccountInfo"
	return &domain.TokenMetadata{Symbol: "SOL-TOKEN", Decimals: 9}, nil
}

func (s *SolanaService) GetHoldersWithTrades(ctx context.Context, tokenAddress string) (map[string][]domain.Trade, error) {
	if s.rpcURL == "" {
		return s.mockSolanaTrades()
	}
	
	// Real implementation using Alchemy Solana RPC
	// body := map[string]interface{}{
	// 	"jsonrpc": "2.0",
	// 	"id": 1,
	// 	"method": "getSignaturesForAddress",
	// 	"params": []interface{}{tokenAddress, map[string]interface{}{"limit": 1000}},
	// }
	// req, _ := http.NewRequest("POST", s.rpcURL, body)
	
	return s.mockSolanaTrades()
}

func (s *SolanaService) GetTrades(ctx context.Context, tokenAddress string) ([]domain.Trade, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *SolanaService) mockSolanaTrades() (map[string][]domain.Trade, error) {
	trades := make(map[string][]domain.Trade)
	
	w1 := "SolAddr1..."
	trades[w1] = []domain.Trade{
		{Type: "buy", Amount: 50, PriceUSD: 20.0, Timestamp: time.Now().Add(-10 * time.Hour)},
	}
	
	w2 := "SolAddr2..."
	trades[w2] = []domain.Trade{
		{Type: "buy", Amount: 100, PriceUSD: 18.0, Timestamp: time.Now().Add(-5 * time.Hour)},
		{Type: "sell", Amount: 100, PriceUSD: 25.0, Timestamp: time.Now().Add(-1 * time.Hour)}, // 700 profit
	}
	
	return trades, nil
}
