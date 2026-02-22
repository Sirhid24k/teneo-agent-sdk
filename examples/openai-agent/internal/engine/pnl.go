package engine

import (
	"sort"

	"github.com/TeneoProtocolAI/teneo-agent-sdk/examples/openai-agent/internal/parser"
)

// WalletPnL holds the realized PnL metrics for a single wallet.
type WalletPnL struct {
	Wallet          string
	RealizedPnL     float64 // total realized PnL in SOL
	CompletedTrades int     // number of buy→sell cycles fully closed
	TotalBuys       int
	TotalSells      int
	WinningTrades   int     // sells with positive PnL
	WinRate         float64 // WinningTrades / TotalSells * 100
}

// buyLot represents a single buy that has not been fully consumed by sells.
type buyLot struct {
	tokenRemaining float64 // tokens still available from this buy
	costPerToken   float64 // SOL cost per token for this buy
}

// ComputePnL calculates realized PnL for each wallet using FIFO cost basis.
//
// For each wallet the swaps are sorted chronologically. Each sell is matched
// against the earliest unconsumed buy lots until the sell quantity is exhausted.
//
// Only wallets with at least one completed buy→sell cycle are returned.
func ComputePnL(swaps []parser.NormalizedSwap) []WalletPnL {
	// Group swaps by wallet
	grouped := make(map[string][]parser.NormalizedSwap)
	for _, s := range swaps {
		grouped[s.Wallet] = append(grouped[s.Wallet], s)
	}

	var results []WalletPnL

	for wallet, walletSwaps := range grouped {
		// Sort by timestamp ascending for FIFO
		sort.Slice(walletSwaps, func(i, j int) bool {
			return walletSwaps[i].Timestamp < walletSwaps[j].Timestamp
		})

		pnl := computeWalletPnL(wallet, walletSwaps)
		if pnl.CompletedTrades > 0 {
			results = append(results, pnl)
		}
	}

	return results
}

// computeWalletPnL runs the FIFO cost basis algorithm for a single wallet.
func computeWalletPnL(wallet string, swaps []parser.NormalizedSwap) WalletPnL {
	var lots []buyLot
	var realizedPnL float64
	var completedTrades int
	var totalBuys, totalSells, winningTrades int

	for _, s := range swaps {
		switch s.Type {
		case "buy":
			totalBuys++
			if s.TokenAmount > 0 {
				lots = append(lots, buyLot{
					tokenRemaining: s.TokenAmount,
					costPerToken:   s.SolAmount / s.TokenAmount,
				})
			}

		case "sell":
			totalSells++
			if len(lots) == 0 || s.TokenAmount <= 0 {
				continue
			}

			tokensToSell := s.TokenAmount
			solReceived := s.SolAmount
			sellPricePerToken := solReceived / tokensToSell
			var sellPnL float64
			traded := false

			for tokensToSell > 0 && len(lots) > 0 {
				lot := &lots[0]

				consumed := lot.tokenRemaining
				if consumed > tokensToSell {
					consumed = tokensToSell
				}

				costBasis := consumed * lot.costPerToken
				revenue := consumed * sellPricePerToken
				sellPnL += revenue - costBasis
				traded = true

				lot.tokenRemaining -= consumed
				tokensToSell -= consumed

				if lot.tokenRemaining <= 1e-12 {
					// Lot fully consumed, remove it
					lots = lots[1:]
				}
			}

			if traded {
				completedTrades++
				realizedPnL += sellPnL
				if sellPnL > 0 {
					winningTrades++
				}
			}
		}
	}

	winRate := 0.0
	if totalSells > 0 {
		winRate = float64(winningTrades) / float64(totalSells) * 100
	}

	return WalletPnL{
		Wallet:          wallet,
		RealizedPnL:     realizedPnL,
		CompletedTrades: completedTrades,
		TotalBuys:       totalBuys,
		TotalSells:      totalSells,
		WinningTrades:   winningTrades,
		WinRate:         winRate,
	}
}
