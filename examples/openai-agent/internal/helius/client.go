package helius

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	// maxPages caps pagination to keep latency within the Teneo task timeout.
	maxPages = 10
	// pageSize is the maximum number of transactions per Helius API call.
	pageSize = 100
)

// Client communicates with the Helius Enhanced Transactions API.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Helius API client.
func NewClient(apiKey, baseURL string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchSwapTransactions retrieves all SWAP-type transactions for a given token
// mint address, paginating up to maxPages (1,000 transactions).
func (c *Client) FetchSwapTransactions(ctx context.Context, tokenMint string) ([]EnhancedTransaction, error) {
	var allTxns []EnhancedTransaction
	var beforeSig string

	for page := 0; page < maxPages; page++ {
		txns, err := c.fetchPage(ctx, tokenMint, beforeSig)
		if err != nil {
			return nil, fmt.Errorf("page %d: %w", page, err)
		}

		if len(txns) == 0 {
			break
		}

		// Filter: keep only transactions with a swap event and no error
		for i := range txns {
			if txns[i].TransactionError != nil {
				continue
			}
			if txns[i].Events.Swap == nil {
				continue
			}
			allTxns = append(allTxns, txns[i])
		}

		// Set cursor for next page
		beforeSig = txns[len(txns)-1].Signature

		log.Printf("📡 Fetched page %d: %d transactions (%d swaps total)", page+1, len(txns), len(allTxns))

		// If we got fewer than a full page, there are no more transactions
		if len(txns) < pageSize {
			break
		}
	}

	log.Printf("✅ Total swap transactions fetched: %d", len(allTxns))
	return allTxns, nil
}

// fetchPage retrieves a single page of enhanced transactions.
func (c *Client) fetchPage(ctx context.Context, address, beforeSig string) ([]EnhancedTransaction, error) {
	endpoint := fmt.Sprintf("%s/v0/addresses/%s/transactions", c.baseURL, address)

	params := url.Values{}
	params.Set("api-key", c.apiKey)
	params.Set("type", "SWAP")
	params.Set("limit", fmt.Sprintf("%d", pageSize))
	if beforeSig != "" {
		params.Set("before", beforeSig)
	}

	fullURL := endpoint + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Helius API returned status %d: %s", resp.StatusCode, string(body))
	}

	var txns []EnhancedTransaction
	if err := json.NewDecoder(resp.Body).Decode(&txns); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return txns, nil
}
