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

type StockListener struct {
	listener *pq.Listener
	channel  string
}

func NewListener(cfg *config.Config) (*StockListener, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBUser,
		cfg.DBPassword,
	)

	reportProblem := func(ev pq.ListenerEventType, err error) {
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

func (l *StockListener) ListenForChanges(ctx context.Context) (<-chan domain.Stock, error) {
	stockChan := make(chan domain.Stock)

	go func() {
		defer close(stockChan)
		logger.Info("Starting to listen for PostgreSQL notifications on channel: %s", l.channel)

		for {
			select {
			case n := <-l.listener.Notify:
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

func (l *StockListener) Close() error {
	return l.listener.Close()
}
