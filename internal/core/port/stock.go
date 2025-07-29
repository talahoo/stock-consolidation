// Package port defines the interfaces for external dependencies
package port

import (
	"context"
	"stock-consolidation/internal/core/domain"
)

// StockEventHandler defines the interface for handling stock events
type StockEventHandler interface {
	HandleStockChange(ctx context.Context, stock domain.Stock) error
}

// StockRepository defines the interface for stock data operations
type StockRepository interface {
	ListenForChanges(ctx context.Context) (<-chan domain.Stock, error)
	Close() error
}
