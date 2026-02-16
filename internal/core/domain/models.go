package domain

import "time"

// AgentInput represents the input parameters for the agent.
type AgentInput struct {
	Chain        string `json:"chain"`        // "ethereum" or "solana"
	TokenAddress string `json:"tokenAddress"` // Contract address of the token
	Limit        int    `json:"limit"`        // Number of top wallets to return
}

// AgentOutput represents the structured output of the agent.
type AgentOutput struct {
	TokenSymbol  string      `json:"token_symbol"`
	CurrentPrice float64     `json:"current_price_usd"`
	TopWallets   []WalletPnL `json:"top_wallets"`
}

// WalletPnL contains the Profit and Loss data for a specific wallet.
type WalletPnL struct {
	Address         string  `json:"wallet_address"`
	TotalBought     float64 `json:"total_bought_tokens"`
	TotalSold       float64 `json:"total_sold_tokens"`
	CurrentBalance  float64 `json:"current_balance_tokens"`
	AverageBuyPrice float64 `json:"avg_buy_price_usd"`
	AverageSellPrice float64 `json:"avg_sell_price_usd"`
	
	RealizedPnL   float64 `json:"realized_pnl_usd"`
	UnrealizedPnL float64 `json:"unrealized_pnl_usd"`
	TotalPnL      float64 `json:"total_pnl_usd"`
	ROI           float64 `json:"roi_percentage"`
}

// Trade represents a single buy or sell event.
type Trade struct {
	Type      string    `json:"type"` // "buy" or "sell"
	Amount    float64   `json:"amount"`
	PriceUSD  float64   `json:"price_usd"`
	Timestamp time.Time `json:"timestamp"`
	TxHash    string    `json:"tx_hash"`
}

// TokenMetadata holds basic information about a token.
type TokenMetadata struct {
	Symbol   string
	Decimals int
	Name     string
}
