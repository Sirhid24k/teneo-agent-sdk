package validator

import (
	"fmt"
	"strconv"
	"strings"
)

// AnalyzeRequest represents a validated user command.
type AnalyzeRequest struct {
	ContractAddress string
	Network         string
	Limit           int
}

// ParseCommand parses the CLI-style input: analyze <contract_address> <network> [limit]
func ParseCommand(task string) (*AnalyzeRequest, error) {
	parts := strings.Fields(strings.TrimSpace(task))
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	// The command word must be "analyze"
	if strings.ToLower(parts[0]) != "analyze" {
		return nil, fmt.Errorf("unknown command %q, expected \"analyze\"", parts[0])
	}

	if len(parts) < 3 {
		return nil, fmt.Errorf("usage: analyze <contract_address> <network> [limit]")
	}

	address := parts[1]
	network := strings.ToLower(parts[2])

	if network != "sol" {
		return nil, fmt.Errorf("unsupported network %q, only \"sol\" is supported", network)
	}

	if err := validateBase58(address); err != nil {
		return nil, fmt.Errorf("invalid contract address: %w", err)
	}

	limit := 5 // default
	if len(parts) >= 4 {
		parsed, err := strconv.Atoi(parts[3])
		if err != nil || parsed <= 0 {
			return nil, fmt.Errorf("limit must be a positive integer, got %q", parts[3])
		}
		limit = parsed
	}

	return &AnalyzeRequest{
		ContractAddress: address,
		Network:         network,
		Limit:           limit,
	}, nil
}

// validateBase58 checks that s is a plausible Solana base58 address.
func validateBase58(s string) error {
	if len(s) < 32 || len(s) > 44 {
		return fmt.Errorf("address length %d out of range [32, 44]", len(s))
	}
	const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	for _, c := range s {
		if !strings.ContainsRune(base58Alphabet, c) {
			return fmt.Errorf("invalid character %q in base58 address", c)
		}
	}
	return nil
}
