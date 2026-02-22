package parser

import (
	"math"
	"strconv"
	"strings"

	"github.com/TeneoProtocolAI/teneo-agent-sdk/examples/openai-agent/internal/helius"
)

// NormalizedSwap represents a single buy or sell action by a wallet.
type NormalizedSwap struct {
	Wallet      string  // the wallet that initiated the swap (feePayer)
	Type        string  // "buy" or "sell"
	TokenAmount float64 // token units (decimal-adjusted)
	SolAmount   float64 // SOL units (lamports / 1e9)
	Timestamp   int64   // unix timestamp
	Signature   string  // transaction signature
}

const lamportsPerSOL = 1_000_000_000.0

// NormalizeSwaps extracts buy/sell records from enhanced transactions for a
// specific token mint.
//
// Classification rules:
//   - Buy  = nativeInput has SOL AND tokenOutputs contains the target mint
//   - Sell = tokenInputs contains the target mint AND nativeOutput has SOL
func NormalizeSwaps(txns []helius.EnhancedTransaction, tokenMint string) []NormalizedSwap {
	mintLower := strings.ToLower(tokenMint)
	var swaps []NormalizedSwap

	for _, tx := range txns {
		swap := tx.Events.Swap
		if swap == nil {
			continue
		}

		wallet := tx.FeePayer

		// Check for a BUY: SOL in → token out
		solIn := parseLamports(swap.NativeInput)
		tokenOut, tokenOutAmount := findTokenAmount(swap.TokenOutputs, mintLower)
		if solIn > 0 && tokenOut {
			swaps = append(swaps, NormalizedSwap{
				Wallet:      wallet,
				Type:        "buy",
				TokenAmount: tokenOutAmount,
				SolAmount:   solIn / lamportsPerSOL,
				Timestamp:   tx.Timestamp,
				Signature:   tx.Signature,
			})
			continue
		}

		// Check for a SELL: token in → SOL out
		tokenIn, tokenInAmount := findTokenAmount(swap.TokenInputs, mintLower)
		solOut := parseLamports(swap.NativeOutput)
		if tokenIn && solOut > 0 {
			swaps = append(swaps, NormalizedSwap{
				Wallet:      wallet,
				Type:        "sell",
				TokenAmount: tokenInAmount,
				SolAmount:   solOut / lamportsPerSOL,
				Timestamp:   tx.Timestamp,
				Signature:   tx.Signature,
			})
		}
	}

	return swaps
}

// parseLamports extracts the lamport amount from a NativeAmount, returning 0 if nil.
func parseLamports(na *helius.NativeAmount) float64 {
	if na == nil {
		return 0
	}
	val, err := strconv.ParseFloat(na.Amount, 64)
	if err != nil {
		return 0
	}
	return val
}

// findTokenAmount checks whether any SwapToken matches the target mint,
// and returns the decimal-adjusted token amount.
func findTokenAmount(tokens []helius.SwapToken, mintLower string) (bool, float64) {
	for _, t := range tokens {
		if strings.ToLower(t.Mint) == mintLower {
			amount := parseRawTokenAmount(t.RawTokenAmount)
			if amount > 0 {
				return true, amount
			}
		}
	}
	return false, 0
}

// parseRawTokenAmount converts a raw token amount string and its decimals into
// a float64 in human-readable units.
func parseRawTokenAmount(raw helius.RawTokenAmount) float64 {
	val, err := strconv.ParseFloat(raw.TokenAmount, 64)
	if err != nil {
		return 0
	}
	if raw.Decimals > 0 {
		val = val / math.Pow(10, float64(raw.Decimals))
	}
	return val
}
