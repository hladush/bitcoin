package repository

import (
	"fmt"

	"github.com/ihladush/bitcoin/internal/models"
)

// SaveTransaction saves a transaction to the database
func (r *SQLiteRepository) SaveTransaction(tx *models.Transaction) error {
	query := `
	INSERT OR REPLACE INTO transactions 
	(hash, address, amount, confirmations, block_height, timestamp, type) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		tx.Hash, tx.Address, tx.Amount, tx.Confirmations,
		tx.BlockHeight, tx.Timestamp, tx.Type,
	)
	if err != nil {
		return fmt.Errorf("failed to save transaction: %w", err)
	}

	return nil
}

// GetTransactionsByAddress retrieves transactions for a specific address with pagination
func (r *SQLiteRepository) GetTransactionsByAddress(address string, limit, offset int) ([]models.Transaction, error) {
	query := `
	SELECT id, hash, address, amount, confirmations, block_height, timestamp, type 
	FROM transactions 
	WHERE address = ? 
	ORDER BY timestamp DESC 
	LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, address, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		err := rows.Scan(
			&tx.ID, &tx.Hash, &tx.Address, &tx.Amount,
			&tx.Confirmations, &tx.BlockHeight, &tx.Timestamp, &tx.Type,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

// TransactionExists checks if a transaction already exists for an address
func (r *SQLiteRepository) TransactionExists(hash, address string) (bool, error) {
	query := `SELECT COUNT(*) FROM transactions WHERE hash = ? AND address = ?`
	
	var count int
	err := r.db.QueryRow(query, hash, address).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check transaction existence: %w", err)
	}

	return count > 0, nil
}

// GetBalance retrieves the calculated balance for an address
func (r *SQLiteRepository) GetBalance(address string) (*models.Balance, error) {
	return r.CalculateBalance(address)
}

// CalculateBalance calculates the balance based on transactions
func (r *SQLiteRepository) CalculateBalance(address string) (*models.Balance, error) {
	// Calculate confirmed balance (transactions with confirmations >= 1)
	confirmedQuery := `
	SELECT COALESCE(SUM(amount), 0) 
	FROM transactions 
	WHERE address = ? AND confirmations >= 1`

	// Calculate unconfirmed balance (transactions with confirmations = 0)
	unconfirmedQuery := `
	SELECT COALESCE(SUM(amount), 0) 
	FROM transactions 
	WHERE address = ? AND confirmations = 0`

	var confirmedBalance, unconfirmedBalance int64

	err := r.db.QueryRow(confirmedQuery, address).Scan(&confirmedBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate confirmed balance: %w", err)
	}

	err = r.db.QueryRow(unconfirmedQuery, address).Scan(&unconfirmedBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate unconfirmed balance: %w", err)
	}

	totalBalance := confirmedBalance + unconfirmedBalance
	balanceBTC := float64(totalBalance) / 100000000 // Convert satoshis to BTC

	return &models.Balance{
		Address:            address,
		ConfirmedBalance:   confirmedBalance,
		UnconfirmedBalance: unconfirmedBalance,
		TotalBalance:       totalBalance,
		BalanceBTC:         balanceBTC,
	}, nil
}

// Close closes the database connection
func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}
