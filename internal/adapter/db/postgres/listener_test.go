package postgres_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"stock-consolidation/internal/adapter/db/postgres"
	"stock-consolidation/internal/core/domain"
	"stock-consolidation/pkg/config"

	"github.com/lib/pq"
)

type mockPGListener struct {
	notifications chan *pq.Notification
	listenError   error
	pingError     error
	closeError    error
	listenCalled  bool
	pingCalled    bool
	closeCalled   bool
}

func (m *mockPGListener) Listen(_ string) error {
	m.listenCalled = true
	return m.listenError
}

func (m *mockPGListener) Ping() error {
	m.pingCalled = true
	return m.pingError
}

func (m *mockPGListener) Close() error {
	m.closeCalled = true
	return m.closeError
}

func (m *mockPGListener) NotificationChannel() <-chan *pq.Notification {
	return m.notifications
}

// closeListener is a helper function to close the listener and check for errors
func closeListener(t *testing.T, listener *postgres.StockListener) {
	if err := listener.Close(); err != nil {
		t.Errorf("Failed to close listener: %v", err)
	}
}

func TestPostgresListener(t *testing.T) {
	testTime := time.Date(2025, 7, 29, 0, 0, 0, 0, time.UTC)

	t.Run("success connection with stock changes", func(t *testing.T) {
		mock := &mockPGListener{
			notifications: make(chan *pq.Notification),
		}

		// Create test stock
		stock := domain.Stock{
			ProductID: 1,
			BranchID:  1,
			Quantity:  10,
			CreatedAt: testTime,
			UpdatedAt: testTime,
		}

		stockData, err := json.Marshal(stock)
		if err != nil {
			t.Fatalf("Failed to marshal stock data: %v", err)
		}

		// Create listener with mock
		listener := postgres.NewListenerWithPG(mock)
		defer closeListener(t, listener)

		// Start listening
		stockChan, err := listener.ListenForChanges(context.Background())
		if err != nil {
			t.Fatalf("Failed to start listening: %v", err)
		}

		// Send test notification in background
		go func() {
			mock.notifications <- &pq.Notification{
				Channel: "stock_changes",
				Extra:   string(stockData),
			}
		}()

		// Verify received stock change with timeout
		select {
		case received := <-stockChan:
			if received.ProductID != stock.ProductID ||
				received.BranchID != stock.BranchID ||
				received.Quantity != stock.Quantity ||
				!received.CreatedAt.Equal(stock.CreatedAt) ||
				!received.UpdatedAt.Equal(stock.UpdatedAt) {
				t.Error("Received stock data does not match sent data")
			}
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for stock notification")
		}
	})

	t.Run("invalid notification data", func(t *testing.T) {
		mock := &mockPGListener{
			notifications: make(chan *pq.Notification),
		}

		// Send invalid JSON data
		go func() {
			mock.notifications <- &pq.Notification{
				Channel: "stock_changes",
				Extra:   "invalid json",
			}
			close(mock.notifications)
		}()

		// Create listener with mock
		listener := postgres.NewListenerWithPG(mock)
		defer closeListener(t, listener)

		// Start listening
		stockChan, err := listener.ListenForChanges(context.Background())
		if err != nil {
			t.Fatalf("Failed to start listening: %v", err)
		}

		// Should not receive anything due to invalid data
		select {
		case <-stockChan:
			t.Error("Should not receive stock for invalid data")
		case <-time.After(100 * time.Millisecond):
			// This is expected
		}
	})

	t.Run("listener error", func(t *testing.T) {
		mock := &mockPGListener{
			listenError: fmt.Errorf("mock listen error"),
		}

		// Create listener with mock
		listener := postgres.NewListenerWithPG(mock)
		defer closeListener(t, listener)

		// Start listening - should get error
		_, err := listener.ListenForChanges(context.Background())
		if err == nil {
			t.Error("Expected error from ListenForChanges, got nil")
		}
	})

	t.Run("nil notification handling", func(t *testing.T) {
		mock := &mockPGListener{
			notifications: make(chan *pq.Notification),
		}

		// Create listener with mock
		listener := postgres.NewListenerWithPG(mock)
		defer closeListener(t, listener)

		// Start listening
		stockChan, err := listener.ListenForChanges(context.Background())
		if err != nil {
			t.Fatalf("Failed to start listening: %v", err)
		}

		// Send nil notification
		go func() {
			mock.notifications <- nil
		}()

		// Should not receive anything due to nil notification
		select {
		case <-stockChan:
			t.Error("Should not receive stock for nil notification")
		case <-time.After(100 * time.Millisecond):
			// This is expected
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		mock := &mockPGListener{
			notifications: make(chan *pq.Notification),
		}

		// Create listener with mock
		listener := postgres.NewListenerWithPG(mock)
		defer closeListener(t, listener)

		ctx, cancel := context.WithCancel(context.Background())

		// Start listening
		stockChan, err := listener.ListenForChanges(ctx)
		if err != nil {
			t.Fatalf("Failed to start listening: %v", err)
		}

		// Cancel context
		cancel()

		// Channel should be closed
		select {
		case _, ok := <-stockChan:
			if ok {
				t.Error("Channel should be closed after context cancellation")
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("Timeout waiting for channel to close")
		}
	})

	// Test successful real connection
	t.Run("real connection", func(t *testing.T) {
		cfg := &config.Config{
			DBHost:     "localhost",
			DBPort:     "5432",
			DBUser:     "admin",
			DBPassword: "admin123",
			DBName:     "stockdb",
		}

		listener, err := postgres.NewListener(cfg)
		if err == nil {
			defer closeListener(t, listener)
			stockChan, err := listener.ListenForChanges(context.Background())
			if err != nil {
				t.Fatalf("Failed to start listening: %v", err)
			}

			// Just test that the channel is created
			select {
			case <-stockChan:
				// Should not receive anything
				t.Error("Unexpected stock change received")
			case <-time.After(100 * time.Millisecond):
				// This is expected
			}
		}
	})
}
