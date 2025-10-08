package models

import "time"

// Transaction represents a Bitcoin transaction
type Transaction struct {
	ID            int       `json:"id" db:"id"`
	Hash          string    `json:"hash" db:"hash"`
	Address       string    `json:"address" db:"address"`
	Amount        int64     `json:"amount" db:"amount"` // Amount in satoshis
	Confirmations int       `json:"confirmations" db:"confirmations"`
	BlockHeight   int       `json:"block_height" db:"block_height"`
	Timestamp     time.Time `json:"timestamp" db:"timestamp"`
	Type          string    `json:"type" db:"type"` // "sent" or "received"
}

// Balance represents the balance for a Bitcoin address
type Balance struct {
	Address           string  `json:"address"`
	ConfirmedBalance  int64   `json:"confirmed_balance"`  // Balance in satoshis
	UnconfirmedBalance int64  `json:"unconfirmed_balance"` // Unconfirmed balance in satoshis
	TotalBalance      int64   `json:"total_balance"`      // Total balance in satoshis
	BalanceBTC        float64 `json:"balance_btc"`        // Balance in BTC
}

// AddressWithBalance combines address info with its current balance
type AddressWithBalance struct {
	Address
	Balance Balance `json:"balance"`
}
