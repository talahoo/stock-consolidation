// Package domain contains the core domain entities and business logic
package domain

import (
	"encoding/json"
	"strings"
	"time"
)

// Stock represents the stock entity
type Stock struct {
	ID        string    `json:"id"`
	ProductID int       `json:"product_id"`
	BranchID  int       `json:"branch_id"`
	Quantity  int       `json:"quantity"`
	Reserved  int       `json:"reserved"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Custom time format for PostgreSQL timestamps
const pgTimeFormat = "2006-01-02T15:04:05.999999"

// UnmarshalJSON implements custom JSON unmarshaling for Stock
func (s *Stock) UnmarshalJSON(data []byte) error {
	// Create an auxiliary type to avoid recursion
	type Aux struct {
		ID        string `json:"id"`
		ProductID int    `json:"product_id"`
		BranchID  int    `json:"branch_id"`
		Quantity  int    `json:"quantity"`
		Reserved  int    `json:"reserved"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	var aux Aux
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse time fields with custom format
	createdAt, err := time.Parse(pgTimeFormat, strings.TrimSpace(aux.CreatedAt))
	if err != nil {
		return err
	}

	updatedAt, err := time.Parse(pgTimeFormat, strings.TrimSpace(aux.UpdatedAt))
	if err != nil {
		return err
	}

	// Assign values to the actual Stock struct
	s.ID = aux.ID
	s.ProductID = aux.ProductID
	s.BranchID = aux.BranchID
	s.Quantity = aux.Quantity
	s.Reserved = aux.Reserved
	s.CreatedAt = createdAt
	s.UpdatedAt = updatedAt

	return nil
}
