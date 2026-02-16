package domain

import "context"

// ChainService defines the operations required to interact with a blockchain.
type ChainService interface {
	// IsSupported checks if the chain is supported by this service.
	IsSupported(chain string) bool
	
	// GetTokenMetadata fetches the decimals and symbol for a token.
	GetTokenMetadata(ctx context.Context, tokenAddress string) (*TokenMetadata, error)
	
	// GetTrades fetches all relevant trades for a token.
	// In a real scenario, this might need pagination or streaming, 
	// but for this agent we'll assume we fetch a reasonable batch.
	GetTrades(ctx context.Context, tokenAddress string) ([]Trade, error)
	
	// GetWalletTrades fetches trades specifically for a wallet if the API structure requires it,
	// or this logic might be internal to GetTrades depending on the provider.
	// For this abstraction, we will assume GetTrades returns a flat list of trades 
	// mapped by wallet internally or we process them. 
	// Actually, a better abstraction for the "Finder" agent might be:
	// GetHoldersWithTrades returns a map of wallet -> trades.
	GetHoldersWithTrades(ctx context.Context, tokenAddress string) (map[string][]Trade, error)
}

// PriceService defines how to get token price data.
type PriceService interface {
	// GetCurrentPrice returns the current USD price of the token.
	GetCurrentPrice(ctx context.Context, chain, tokenAddress string) (float64, error)
}

// PnLCalculator defines the logic to compute PnL.
type PnLCalculator interface {
	Calculate(trades []Trade, currentPrice float64) *WalletPnL
}
