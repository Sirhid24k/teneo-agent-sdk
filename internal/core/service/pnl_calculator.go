package service

import (
	"sort"

	"github.com/TeneoProtocolAI/teneo-agent-sdk/internal/core/domain"
)

type PnLCalculator struct{}

func NewPnLCalculator() domain.PnLCalculator {
	return &PnLCalculator{}
}

// Calculate computes PnL metrics for a set of trades.
func (p *PnLCalculator) Calculate(trades []domain.Trade, currentPrice float64) *domain.WalletPnL {
	var totalBought, totalSold float64
	var totalCost, totalRevenue float64
	var realizedPnL float64

	// Sort trades by date (ascending) to process sequentially if needed
	// For this simple aggregation, order might not strictly matter for averages, 
	// but good practice.
	sort.Slice(trades, func(i, j int) bool {
		return trades[i].Timestamp.Before(trades[j].Timestamp)
	})

	for _, t := range trades {
		if t.Type == "buy" {
			totalBought += t.Amount
			totalCost += t.Amount * t.PriceUSD
		} else if t.Type == "sell" {
			// Basic FIFO or Average Cost logic could differ. 
			// Here we use Average Buy Price logic for PnL.
			// Realized PnL on this sell = (SellPrice - AvgBuyPrice) * Amount
			
			// However, simple approach:
			totalSold += t.Amount
			totalRevenue += t.Amount * t.PriceUSD
		}
	}

	currentBalance := totalBought - totalSold
	// Sanity check for data inconsistencies
	if currentBalance < 0 {
		currentBalance = 0 
	}

	avgBuyPrice := 0.0
	if totalBought > 0 {
		avgBuyPrice = totalCost / totalBought
	}

	avgSellPrice := 0.0
	if totalSold > 0 {
		avgSellPrice = totalRevenue / totalSold
	}

	// Realized PnL: Revenue - (Cost of Sold Tokens)
	// Cost of Sold Tokens = TotalSold * AvgBuyPrice
	costOfSold := totalSold * avgBuyPrice
	realizedPnL = totalRevenue - costOfSold

	// Unrealized PnL: Value of current holdings - Cost of current holdings
	// Cost of current holdings = CurrentBalance * AvgBuyPrice
	currentValue := currentBalance * currentPrice
	costOfHeld := currentBalance * avgBuyPrice
	unrealizedPnL := currentValue - costOfHeld

	totalPnL := realizedPnL + unrealizedPnL

	// ROI = TotalPnL / Total Cost Invested (TotalCost) * 100
	roi := 0.0
	if totalCost > 0 {
		roi = (totalPnL / totalCost) * 100
	}

	return &domain.WalletPnL{
		TotalBought:     totalBought,
		TotalSold:       totalSold,
		CurrentBalance:  currentBalance,
		AverageBuyPrice: avgBuyPrice,
		AverageSellPrice: avgSellPrice,
		RealizedPnL:     realizedPnL,
		UnrealizedPnL:   unrealizedPnL,
		TotalPnL:        totalPnL,
		ROI:             roi,
	}
}
