// Package repository provides data access layer for the Bitcoin tracker application
package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ihladush/bitcoin/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

// Repository interface defines the contract for data access
type Repository interface {
	// Address operations
	AddAddress(address, label string) (*models.Address, error)
	RemoveAddress(address string) error
	GetAddress(address string) (*models.Address, error)
	GetAllAddresses() ([]models.Address, error)
	UpdateLastSynced(address string, syncTime time.Time) error

	// Transaction operations
	SaveTransaction(tx *models.Transaction) error
	GetTransactionsByAddress(address string, limit, offset int) ([]models.Transaction, error)
	TransactionExists(hash, address string) (bool, error)

	// Balance operations
	GetBalance(address string) (*models.Balance, error)
	CalculateBalance(address string) (*models.Balance, error)
}

// SQLiteRepository implements Repository interface using SQLite
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates a new SQLite repository
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	repo := &SQLiteRepository{db: db}
	if err := repo.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return repo, nil
}

// createTables creates the necessary database tables
func (r *SQLiteRepository) createTables() error {
	// Create addresses table
	addressTable := `
	CREATE TABLE IF NOT EXISTS addresses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		address TEXT UNIQUE NOT NULL,
		label TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_synced DATETIME
	);`

	// Create transactions table
	transactionTable := `
	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		hash TEXT NOT NULL,
		address TEXT NOT NULL,
		amount INTEGER NOT NULL,
		confirmations INTEGER NOT NULL,
		block_height INTEGER NOT NULL,
		timestamp DATETIME NOT NULL,
		type TEXT NOT NULL,
		UNIQUE(hash, address),
		FOREIGN KEY(address) REFERENCES addresses(address) ON DELETE CASCADE
	);`

	// Create indexes for better performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_transactions_address ON transactions(address);",
		"CREATE INDEX IF NOT EXISTS idx_transactions_timestamp ON transactions(timestamp);",
		"CREATE INDEX IF NOT EXISTS idx_transactions_hash ON transactions(hash);",
	}

	// Execute table creation
	if _, err := r.db.Exec(addressTable); err != nil {
		return fmt.Errorf("failed to create addresses table: %w", err)
	}

	if _, err := r.db.Exec(transactionTable); err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}

	// Create indexes
	for _, index := range indexes {
		if _, err := r.db.Exec(index); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}

// AddAddress adds a new address to track
func (r *SQLiteRepository) AddAddress(address, label string) (*models.Address, error) {
	query := `INSERT INTO addresses (address, label) VALUES (?, ?) RETURNING id, created_at`
	
	var addr models.Address
	addr.Address = address
	addr.Label = label
	
	err := r.db.QueryRow(query, address, label).Scan(&addr.ID, &addr.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to add address: %w", err)
	}

	return &addr, nil
}

// RemoveAddress removes an address from tracking
func (r *SQLiteRepository) RemoveAddress(address string) error {
	query := `DELETE FROM addresses WHERE address = ?`
	result, err := r.db.Exec(query, address)
	if err != nil {
		return fmt.Errorf("failed to remove address: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("address not found: %s", address)
	}

	return nil
}

// GetAddress retrieves a specific address
func (r *SQLiteRepository) GetAddress(address string) (*models.Address, error) {
	query := `SELECT id, address, label, created_at, last_synced FROM addresses WHERE address = ?`
	
	var addr models.Address
	var lastSynced sql.NullTime
	
	err := r.db.QueryRow(query, address).Scan(
		&addr.ID, &addr.Address, &addr.Label, &addr.CreatedAt, &lastSynced,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("address not found: %s", address)
		}
		return nil, fmt.Errorf("failed to get address: %w", err)
	}

	if lastSynced.Valid {
		addr.LastSynced = &lastSynced.Time
	}

	return &addr, nil
}

// GetAllAddresses retrieves all tracked addresses
func (r *SQLiteRepository) GetAllAddresses() ([]models.Address, error) {
	query := `SELECT id, address, label, created_at, last_synced FROM addresses ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses: %w", err)
	}
	defer rows.Close()

	var addresses []models.Address
	for rows.Next() {
		var addr models.Address
		var lastSynced sql.NullTime
		
		err := rows.Scan(&addr.ID, &addr.Address, &addr.Label, &addr.CreatedAt, &lastSynced)
		if err != nil {
			return nil, fmt.Errorf("failed to scan address: %w", err)
		}

		if lastSynced.Valid {
			addr.LastSynced = &lastSynced.Time
		}

		addresses = append(addresses, addr)
	}

	return addresses, nil
}

// UpdateLastSynced updates the last sync time for an address
func (r *SQLiteRepository) UpdateLastSynced(address string, syncTime time.Time) error {
	query := `UPDATE addresses SET last_synced = ? WHERE address = ?`
	_, err := r.db.Exec(query, syncTime, address)
	if err != nil {
		return fmt.Errorf("failed to update last synced: %w", err)
	}
	return nil
}
