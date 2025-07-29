package postgres_test

import (
	"context"
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
	// Use time format that matches the domain parser
	testTime := time.Date(2025, 7, 29, 0, 0, 0, 0, time.UTC)

	t.Run("success connection with stock changes", func(t *testing.T) {
		mock := &mockPGListener{
			notifications: make(chan *pq.Notification),
		}

		// Create test stock with correct time format
		stock := domain.Stock{
			ProductID: 1,
			BranchID:  1,
			Quantity:  10,
			CreatedAt: testTime,
			UpdatedAt: testTime,
		}

		// Create JSON manually to match the expected format
		stockJSON := fmt.Sprintf(`{
			"product_id": %d,
			"branch_id": %d,
			"quantity": %d,
			"reserved": 0,
			"created_at": "%s",
			"updated_at": "%s"
		}`, stock.ProductID, stock.BranchID, stock.Quantity,
			testTime.Format("2006-01-02T15:04:05.999999"),
			testTime.Format("2006-01-02T15:04:05.999999"))

		stockData := []byte(stockJSON)

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

	t.Run("notification processing error", func(t *testing.T) {
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

		// Send notification with invalid JSON
		go func() {
			mock.notifications <- &pq.Notification{
				Channel: "stock_changes",
				Extra:   "invalid json",
			}
		}()

		// Should not receive anything due to invalid JSON
		select {
		case <-stockChan:
			t.Error("Should not receive stock for invalid JSON")
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

	// Test successful real connection (only run if DB is available)
	t.Run("real connection", func(t *testing.T) {
		// Skip this test in CI environments where DB might not be available
		if testing.Short() {
			t.Skip("Skipping real connection test in short mode")
		}

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
