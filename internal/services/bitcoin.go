// Package services provides business logic for the Bitcoin tracker application
package services

import (
	"fmt"
	"time"

	"github.com/ihladush/bitcoin/internal/clients"
	"github.com/ihladush/bitcoin/internal/models"
	"github.com/ihladush/bitcoin/internal/repository"
)

// BitcoinService handles business logic for Bitcoin tracking
type BitcoinService struct {
	repo   repository.Repository
	client clients.BitcoinClient
}

// NewBitcoinService creates a new Bitcoin service
func NewBitcoinService(repo repository.Repository, client clients.BitcoinClient) *BitcoinService {
	return &BitcoinService{
		repo:   repo,
		client: client,
	}
}

// AddAddress adds a new Bitcoin address for tracking
func (s *BitcoinService) AddAddress(address, label string) (*models.Address, error) {
	// Validate address format
	if !s.client.IsValidAddress(address) {
		return nil, fmt.Errorf("invalid Bitcoin address: %s", address)
	}

	// Check if address already exists
	existingAddr, err := s.repo.GetAddress(address)
	if err == nil && existingAddr != nil {
		return nil, fmt.Errorf("address already being tracked: %s", address)
	}

	// Add address to repository
	addr, err := s.repo.AddAddress(address, label)
	if err != nil {
		return nil, fmt.Errorf("failed to add address: %w", err)
	}

	// Perform initial sync
	if err := s.SyncAddress(address); err != nil {
		// Log the error but don't fail the add operation
		fmt.Printf("Warning: initial sync failed for address %s: %v\n", address, err)
	}

	return addr, nil
}

// RemoveAddress removes a Bitcoin address from tracking
func (s *BitcoinService) RemoveAddress(address string) error {
	return s.repo.RemoveAddress(address)
}

// GetAllAddresses returns all tracked addresses with their balances
func (s *BitcoinService) GetAllAddresses() ([]models.AddressWithBalance, error) {
	addresses, err := s.repo.GetAllAddresses()
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses: %w", err)
	}

	var addressesWithBalance []models.AddressWithBalance
	for _, addr := range addresses {
		balance, err := s.repo.GetBalance(addr.Address)
		if err != nil {
			// Return zero balance if calculation fails
			balance = &models.Balance{
				Address:            addr.Address,
				ConfirmedBalance:   0,
				UnconfirmedBalance: 0,
				TotalBalance:       0,
				BalanceBTC:         0,
			}
		}

		addressWithBalance := models.AddressWithBalance{
			Address: addr,
			Balance: *balance,
		}
		addressesWithBalance = append(addressesWithBalance, addressWithBalance)
	}

	return addressesWithBalance, nil
}

// GetAddress returns a specific address with its balance
func (s *BitcoinService) GetAddress(address string) (*models.AddressWithBalance, error) {
	addr, err := s.repo.GetAddress(address)
	if err != nil {
		return nil, fmt.Errorf("address not found: %w", err)
	}

	balance, err := s.repo.GetBalance(address)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	return &models.AddressWithBalance{
		Address: *addr,
		Balance: *balance,
	}, nil
}

// GetBalance returns the current balance for an address
func (s *BitcoinService) GetBalance(address string) (*models.Balance, error) {
	// Verify address exists in our tracking
	_, err := s.repo.GetAddress(address)
	if err != nil {
		return nil, fmt.Errorf("address not being tracked: %w", err)
	}

	return s.repo.GetBalance(address)
}

// GetTransactions returns transactions for an address with pagination
func (s *BitcoinService) GetTransactions(address string, limit, offset int) ([]models.Transaction, error) {
	// Verify address exists in our tracking
	_, err := s.repo.GetAddress(address)
	if err != nil {
		return nil, fmt.Errorf("address not being tracked: %w", err)
	}

	// Set default limit if not provided
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100 // Maximum limit
	}

	return s.repo.GetTransactionsByAddress(address, limit, offset)
}

// SyncAddress synchronizes transaction data for a specific address
func (s *BitcoinService) SyncAddress(address string) error {
	// Verify address exists in our tracking
	_, err := s.repo.GetAddress(address)
	if err != nil {
		return fmt.Errorf("address not being tracked: %w", err)
	}

	// Fetch transactions from blockchain API
	transactions, err := s.client.GetTransactions(address, 100)
	if err != nil {
		return fmt.Errorf("failed to fetch transactions from API: %w", err)
	}

	// Save new transactions to database
	var savedCount int
	for _, tx := range transactions {
		// Check if transaction already exists
		exists, err := s.repo.TransactionExists(tx.Hash, address)
		if err != nil {
			return fmt.Errorf("failed to check transaction existence: %w", err)
		}

		if !exists {
			if err := s.repo.SaveTransaction(&tx); err != nil {
				return fmt.Errorf("failed to save transaction: %w", err)
			}
			savedCount++
		}
	}

	// Update last synced time
	if err := s.repo.UpdateLastSynced(address, time.Now()); err != nil {
		return fmt.Errorf("failed to update last synced time: %w", err)
	}

	fmt.Printf("Synced %d new transactions for address %s\n", savedCount, address)
	return nil
}

// SyncAllAddresses synchronizes all tracked addresses
func (s *BitcoinService) SyncAllAddresses() error {
	addresses, err := s.repo.GetAllAddresses()
	if err != nil {
		return fmt.Errorf("failed to get addresses for sync: %w", err)
	}

	var errors []error
	for _, addr := range addresses {
		if err := s.SyncAddress(addr.Address); err != nil {
			errors = append(errors, fmt.Errorf("sync failed for %s: %w", addr.Address, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("sync completed with %d errors", len(errors))
	}

	return nil
}
