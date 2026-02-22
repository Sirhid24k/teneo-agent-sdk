package ranking

import (
	"fmt"
	"sort"
	"strings"

	"github.com/TeneoProtocolAI/teneo-agent-sdk/examples/openai-agent/internal/engine"
)

// RankWallets filters to profitable wallets, sorts by PnL descending
// (with tie-breakers on win rate then trade count), and truncates to limit.
func RankWallets(wallets []engine.WalletPnL, limit int) []engine.WalletPnL {
	// Filter: only wallets with positive realized PnL
	var profitable []engine.WalletPnL
	for _, w := range wallets {
		if w.RealizedPnL > 0 {
			profitable = append(profitable, w)
		}
	}

	// Sort: PnL desc → WinRate desc → CompletedTrades desc
	sort.Slice(profitable, func(i, j int) bool {
		if profitable[i].RealizedPnL != profitable[j].RealizedPnL {
			return profitable[i].RealizedPnL > profitable[j].RealizedPnL
		}
		if profitable[i].WinRate != profitable[j].WinRate {
			return profitable[i].WinRate > profitable[j].WinRate
		}
		return profitable[i].CompletedTrades > profitable[j].CompletedTrades
	})

	if len(profitable) > limit {
		profitable = profitable[:limit]
	}

	return profitable
}

// FormatOutput builds the human-readable output string per the spec.
func FormatOutput(wallets []engine.WalletPnL, contractAddress string) string {
	if len(wallets) == 0 {
		return fmt.Sprintf("No profitable Alpha Wallets found for %s", contractAddress)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d Alpha Wallets for %s\n\n", len(wallets), contractAddress))

	for i, w := range wallets {
		pnlSign := "+"
		if w.RealizedPnL < 0 {
			pnlSign = ""
		}
		sb.WriteString(fmt.Sprintf("%d. %s (Realized PnL: %s%.4f SOL, Trades: %d, WinRate: %.0f%%)\n",
			i+1,
			w.Wallet,
			pnlSign,
			w.RealizedPnL,
			w.CompletedTrades,
			w.WinRate,
		))
	}

	return strings.TrimRight(sb.String(), "\n")
}
