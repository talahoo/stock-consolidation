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

func closeListener(t *testing.T, listener *postgres.StockListener) {
	if err := listener.Close(); err != nil {
		t.Errorf("Failed to close listener: %v", err)
	}
}

func createTestStock() (domain.Stock, string) {
	testTime := time.Date(2025, 7, 29, 0, 0, 0, 0, time.UTC)
	stock := domain.Stock{
		ProductID: 1,
		BranchID:  1,
		Quantity:  10,
		CreatedAt: testTime,
		UpdatedAt: testTime,
	}

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

	return stock, stockJSON
}

func assertReceivedStock(t *testing.T, received, expected domain.Stock) {
	if received.ProductID != expected.ProductID ||
		received.BranchID != expected.BranchID ||
		received.Quantity != expected.Quantity ||
		!received.CreatedAt.Equal(expected.CreatedAt) ||
		!received.UpdatedAt.Equal(expected.UpdatedAt) {
		t.Error("Received stock data does not match sent data")
	}
}

func TestPostgresListener(t *testing.T) {
	t.Run("success connection with stock changes", func(t *testing.T) {
		mock := &mockPGListener{
			notifications: make(chan *pq.Notification),
		}
		expectedStock, stockJSON := createTestStock()

		listener := postgres.NewListenerWithPG(mock)
		defer closeListener(t, listener)

		stockChan, err := listener.ListenForChanges(context.Background())
		if err != nil {
			t.Fatalf("Failed to start listening: %v", err)
		}

		go func() {
			mock.notifications <- &pq.Notification{
				Channel: "stock_changes",
				Extra:   stockJSON,
			}
		}()

		select {
		case received := <-stockChan:
			assertReceivedStock(t, received, expectedStock)
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for stock notification")
		}
	})

	t.Run("notification processing error", func(t *testing.T) {
		mock := &mockPGListener{
			notifications: make(chan *pq.Notification),
		}
		listener := postgres.NewListenerWithPG(mock)
		defer closeListener(t, listener)

		stockChan, err := listener.ListenForChanges(context.Background())
		if err != nil {
			t.Fatalf("Failed to start listening: %v", err)
		}

		go func() {
			mock.notifications <- &pq.Notification{
				Channel: "stock_changes",
				Extra:   "invalid json",
			}
		}()

		select {
		case <-stockChan:
			t.Error("Should not receive stock for invalid JSON")
		case <-time.After(100 * time.Millisecond):
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		mock := &mockPGListener{
			notifications: make(chan *pq.Notification),
		}
		listener := postgres.NewListenerWithPG(mock)
		defer closeListener(t, listener)

		ctx, cancel := context.WithCancel(context.Background())
		stockChan, err := listener.ListenForChanges(ctx)
		if err != nil {
			t.Fatalf("Failed to start listening: %v", err)
		}
		cancel()

		select {
		case _, ok := <-stockChan:
			if ok {
				t.Error("Channel should be closed after context cancellation")
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("Timeout waiting for channel to close")
		}
	})

	t.Run("real connection", func(t *testing.T) {
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

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		listener, err := postgres.NewListener(cfg)
		if err != nil {
			t.Skipf("Skipping real connection test due to connection error: %v", err)
		}
		defer closeListener(t, listener)

		stockChan, err := listener.ListenForChanges(ctx)
		if err != nil {
			t.Fatalf("Failed to start listening: %v", err)
		}

		select {
		case <-stockChan:
			t.Error("Unexpected stock change received")
		case <-time.After(100 * time.Millisecond):
		}
	})
}
