package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// UpdateAgentVisibility sets an agent's visibility to public or private on the Teneo network.
//
// This is a standalone utility that can be called from any context — no running agent required.
// The agent must have been deployed and connected at least once before visibility can be changed.
//
// Parameters:
//   - backendURL: The Teneo backend URL (e.g. "https://backend.developer.chatroom.teneo-protocol.ai")
//   - agentName: The agent's name as registered (used to derive the agent ID)
//   - creatorWallet: The Ethereum wallet address that owns the agent's NFT
//   - public: true to make the agent publicly visible, false to make it private
//
// Example — make an agent public:
//
//	err := agent.UpdateAgentVisibility(
//	    "https://backend.developer.chatroom.teneo-protocol.ai",
//	    "Interior Architecture Advisor",
//	    "0x10aaF658FA638a1A153dD3730236088950Ab7572",
//	    true,
//	)
//
// Example — make an agent private:
//
//	err := agent.UpdateAgentVisibility(
//	    "https://backend.developer.chatroom.teneo-protocol.ai",
//	    "Interior Architecture Advisor",
//	    "0x10aaF658FA638a1A153dD3730236088950Ab7572",
//	    false,
//	)
//
// HTTP API equivalent (for non-Go clients):
//
//	POST {backendURL}/api/agents/{agent-id}/visibility
//	Content-Type: application/json
//
//	{
//	    "is_public": true,
//	    "creator_wallet": "0xYourWalletAddress"
//	}
//
// The agent ID is derived from the agent name: lowercased, spaces replaced with hyphens,
// non-alphanumeric characters removed. For example "Interior Architecture Advisor" becomes
// "interior-architecture-advisor".
func UpdateAgentVisibility(backendURL, agentName, creatorWallet string, public bool) error {
	agentID := generateAgentID(agentName)
	backendURL = strings.TrimRight(backendURL, "/")

	reqBody, err := json.Marshal(map[string]interface{}{
		"is_public":      public,
		"creator_wallet": creatorWallet,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/agents/%s/visibility", backendURL, agentID)
	resp, err := http.Post(url, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(body, &errResp) == nil && errResp.Error != "" {
			return fmt.Errorf("visibility update failed: %s", errResp.Error)
		}
		return fmt.Errorf("visibility update failed with status %d", resp.StatusCode)
	}

	return nil
}
