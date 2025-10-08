// Package clients provides external API clients for blockchain data
package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ihladush/bitcoin/internal/models"
)

// BlockchairClient interacts with Blockchair API
type BlockchairClient struct {
	baseURL    string
	httpClient *http.Client
}

// BlockchairAddressResponse represents the response from Blockchair address API
type BlockchairAddressResponse struct {
	Data map[string]BlockchairAddressData `json:"data"`
}

// BlockchairAddressData represents address data from Blockchair
type BlockchairAddressData struct {
	Address struct {
		Balance               int64  `json:"balance"`
		BalanceUsd            float64 `json:"balance_usd"`
		Received              int64  `json:"received"`
		Spent                 int64  `json:"spent"`
		OutputCount           int    `json:"output_count"`
		UnspentOutputCount    int    `json:"unspent_output_count"`
		FirstSeenReceiving    string `json:"first_seen_receiving"`
		LastSeenReceiving     string `json:"last_seen_receiving"`
		FirstSeenSpending     string `json:"first_seen_spending"`
		LastSeenSpending      string `json:"last_seen_spending"`
		TransactionCount      int    `json:"transaction_count"`
	} `json:"address"`
}

// BlockchairTransactionsResponse represents the response from Blockchair transactions API
type BlockchairTransactionsResponse struct {
	Data struct {
		Transactions []BlockchairTransaction `json:"transactions"`
	} `json:"data"`
}

// BlockchairTransaction represents a transaction from Blockchair API
type BlockchairTransaction struct {
	BlockID         int64     `json:"block_id"`
	Hash            string    `json:"hash"`
	Time            time.Time `json:"time"`
	BalanceChange   int64     `json:"balance_change"`
	InputTotalValue int64     `json:"input_total_value"`
	OutputTotalValue int64    `json:"output_total_value"`
}

// BitcoinClient interface defines the contract for Bitcoin blockchain clients
type BitcoinClient interface {
	GetBalance(address string) (*models.Balance, error)
	GetTransactions(address string, limit int) ([]models.Transaction, error)
	IsValidAddress(address string) bool
}

// NewBlockchairClient creates a new Blockchair client
func NewBlockchairClient() *BlockchairClient {
	return &BlockchairClient{
		baseURL: "https://api.blockchair.com/bitcoin",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetBalance retrieves the current balance for a Bitcoin address
func (c *BlockchairClient) GetBalance(address string) (*models.Balance, error) {
	url := fmt.Sprintf("%s/dashboards/address/%s", c.baseURL, address)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch balance: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var addressResp BlockchairAddressResponse
	if err := json.NewDecoder(resp.Body).Decode(&addressResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	addressData, exists := addressResp.Data[address]
	if !exists {
		return nil, fmt.Errorf("address data not found in response")
	}

	// Convert satoshis to BTC
	balanceBTC := float64(addressData.Address.Balance) / 100000000

	return &models.Balance{
		Address:            address,
		ConfirmedBalance:   addressData.Address.Balance,
		UnconfirmedBalance: 0, // Blockchair doesn't separate confirmed/unconfirmed in this endpoint
		TotalBalance:       addressData.Address.Balance,
		BalanceBTC:         balanceBTC,
	}, nil
}

// GetTransactions retrieves recent transactions for a Bitcoin address
func (c *BlockchairClient) GetTransactions(address string, limit int) ([]models.Transaction, error) {
	url := fmt.Sprintf("%s/dashboards/address/%s?limit=%d", c.baseURL, address, limit)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var transResp BlockchairTransactionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&transResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var transactions []models.Transaction
	for _, tx := range transResp.Data.Transactions {
		// Determine transaction type based on balance change
		txType := "received"
		if tx.BalanceChange < 0 {
			txType = "sent"
		}

		// Calculate confirmations (simplified - we assume recent blocks)
		confirmations := 6 // Default to 6 confirmations for simplicity
		if tx.BlockID == 0 {
			confirmations = 0 // Unconfirmed transaction
		}

		transaction := models.Transaction{
			Hash:          tx.Hash,
			Address:       address,
			Amount:        tx.BalanceChange,
			Confirmations: confirmations,
			BlockHeight:   int(tx.BlockID),
			Timestamp:     tx.Time,
			Type:          txType,
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// IsValidAddress checks if a Bitcoin address is valid (basic check)
func (c *BlockchairClient) IsValidAddress(address string) bool {
	// Basic validation - check length and prefixes
	if len(address) < 26 || len(address) > 62 {
		return false
	}

	// Check for valid Bitcoin address prefixes
	validPrefixes := []string{"1", "3", "bc1"}
	for _, prefix := range validPrefixes {
		if len(address) >= len(prefix) && address[:len(prefix)] == prefix {
			return true
		}
	}

	return false
}

// GetDetailedTransactions retrieves detailed transaction information for an address
func (c *BlockchairClient) GetDetailedTransactions(address string) ([]models.Transaction, error) {
	// This would require a more complex API call that gets individual transaction details
	// For now, we'll use the simpler dashboard endpoint
	return c.GetTransactions(address, 50)
}
