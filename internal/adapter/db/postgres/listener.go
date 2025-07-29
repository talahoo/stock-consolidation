// Package postgres provides PostgreSQL database adapter functionality
package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"stock-consolidation/internal/core/domain"
	"stock-consolidation/pkg/config"
	"stock-consolidation/pkg/logger"

	"github.com/lib/pq"
)

// PGListener defines the interface for PostgreSQL listener operations
type PGListener interface {
	Listen(channel string) error
	Ping() error
	Close() error
	NotificationChannel() <-chan *pq.Notification
}

// StockListener handles PostgreSQL notifications for stock changes
type StockListener struct {
	listener PGListener
	channel  string
}

// NewListenerWithPG creates a new StockListener with a custom PGListener
func NewListenerWithPG(listener PGListener) *StockListener {
	return &StockListener{
		listener: listener,
		channel:  "stock_changes",
	}
}

// NewListener creates a new StockListener with PostgreSQL connection
func NewListener(cfg *config.Config) (*StockListener, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBUser,
		cfg.DBPassword,
	)

	reportProblem := func(_ pq.ListenerEventType, err error) {
		if err != nil {
			logger.Error("Postgres listener error: %v", err)
		}
	}

	listener := pq.NewListener(connStr, 10, 0, reportProblem)
	if err := listener.Listen("stock_changes"); err != nil {
		return nil, fmt.Errorf("failed to start listening: %v", err)
	}

	// Verify the connection
	if err := listener.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %v", err)
	}

	logger.Info("Successfully connected to PostgreSQL and listening on channel: stock_changes")

	return &StockListener{
		listener: listener,
		channel:  "stock_changes",
	}, nil
}

// ListenForChanges starts listening for stock change notifications
func (l *StockListener) ListenForChanges(ctx context.Context) (<-chan domain.Stock, error) {
	stockChan := make(chan domain.Stock)

	go func() {
		defer close(stockChan)
		logger.Info("Starting to listen for PostgreSQL notifications on channel: %s", l.channel)

		for {
			select {
			case n := <-l.listener.NotificationChannel():
				logger.Info("Received notification: %+v", n)
				if n == nil {
					logger.Info("Received empty notification")
					continue
				}

				var stock domain.Stock
				if err := json.Unmarshal([]byte(n.Extra), &stock); err != nil {
					logger.Error("Error unmarshaling notification: %v", err)
					continue
				}
				logger.Info("Received stock change notification for product %d in branch %d", stock.ProductID, stock.BranchID)

				select {
				case stockChan <- stock:
				case <-ctx.Done():
					return
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	return stockChan, nil
}

// Close closes the PostgreSQL listener
func (l *StockListener) Close() error {
	return l.listener.Close()
}
