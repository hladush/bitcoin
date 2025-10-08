// Package models contains the core data structures for the Bitcoin tracker application
package models

import "time"

// Address represents a Bitcoin address being tracked
type Address struct {
	ID         int       `json:"id" db:"id"`
	Address    string    `json:"address" db:"address"`
	Label      string    `json:"label" db:"label"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	LastSynced *time.Time `json:"last_synced" db:"last_synced"`
}

// AddAddressRequest represents the request payload for adding an address
type AddAddressRequest struct {
	Address string `json:"address"`
	Label   string `json:"label,omitempty"`
}
