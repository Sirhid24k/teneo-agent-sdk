package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/TeneoProtocolAI/teneo-agent-sdk/internal/core/domain"
)

type AgentService struct {
	chainServices []domain.ChainService
	priceService  domain.PriceService
	pnlCalculator domain.PnLCalculator
}

func NewAgentService(
	chains []domain.ChainService,
	price domain.PriceService,
	calc domain.PnLCalculator,
) *AgentService {
	return &AgentService{
		chainServices: chains,
		priceService:  price,
		pnlCalculator: calc,
	}
}

func (s *AgentService) AnalyzeToken(ctx context.Context, input domain.AgentInput) (*domain.AgentOutput, error) {
	// 1. Find correct chain service
	var chainService domain.ChainService
	for _, cs := range s.chainServices {
		if cs.IsSupported(input.Chain) {
			chainService = cs
			break
		}
	}
	if chainService == nil {
		return nil, fmt.Errorf("chain %s not supported", input.Chain)
	}

	// 2. Fetch Token Metadata
	meta, err := chainService.GetTokenMetadata(ctx, input.TokenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get token metadata: %w", err)
	}

	// 3. Fetch Price
	price, err := s.priceService.GetCurrentPrice(ctx, input.Chain, input.TokenAddress)
	if err != nil {
		// Proceed with 0 price or error? Agent usually needs price.
		return nil, fmt.Errorf("failed to get price: %w", err)
	}

	// 4. Fetch Trades/Holders
	holdersMap, err := chainService.GetHoldersWithTrades(ctx, input.TokenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get trades: %w", err)
	}

	// 5. Calculate PnL for each wallet
	var results []domain.WalletPnL
	for addr, trades := range holdersMap {
		// Optional: Filter out logic here (contracts, deployer) if not done in ChainService
		// For now we assume ChainService returns relevant user wallets or we filter here if we had metadata.
		
		stats := s.pnlCalculator.Calculate(trades, price)
		stats.Address = addr
		
		// Filter out zero activity or tiny dust if needed
		if stats.TotalBought == 0 && stats.TotalSold == 0 {
			continue
		}
		
		results = append(results, *stats)
	}

	// 6. Rank by Total PnL (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].TotalPnL > results[j].TotalPnL
	})

	// 7. Limit output
	limit := input.Limit
	if limit > len(results) {
		limit = len(results)
	}
	topWallets := results[:limit]

	return &domain.AgentOutput{
		TokenSymbol:  meta.Symbol,
		CurrentPrice: price,
		TopWallets:   topWallets,
	}, nil
}
