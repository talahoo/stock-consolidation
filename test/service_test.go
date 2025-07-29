package service_test

import (
	"context"
	"testing"
	"time"

	"stock-consolidation/internal/core/domain"
	"stock-consolidation/internal/service"
)

type mockStockRepository struct {
	stockChan chan domain.Stock
}

func (m *mockStockRepository) ListenForChanges(_ context.Context) (<-chan domain.Stock, error) {
	return m.stockChan, nil
}

func (m *mockStockRepository) Close() error {
	close(m.stockChan)
	return nil
}

func TestStockService(t *testing.T) {
	t.Run("successful stock change notification", func(t *testing.T) {
		mockRepo := &mockStockRepository{
			stockChan: make(chan domain.Stock),
		}

		service := service.NewStockService(mockRepo)

		// Start listening in background
		_, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			err := service.ListenForChanges()
			if err != nil {
				t.Errorf("ListenForChanges returned error: %v", err)
			}
		}()

		// Send multiple test stock changes
		testCases := []domain.Stock{
			{
				ID:        "123e4567-e89b-12d3-a456-426614174000",
				ProductID: 1,
				BranchID:  1,
				Quantity:  10,
				Reserved:  0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "223e4567-e89b-12d3-a456-426614174001",
				ProductID: 2,
				BranchID:  2,
				Quantity:  20,
				Reserved:  5,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		for _, tc := range testCases {
			mockRepo.stockChan <- tc
			// Allow some time for processing
			time.Sleep(50 * time.Millisecond)
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		mockRepo := &mockStockRepository{
			stockChan: make(chan domain.Stock),
		}

		service := service.NewStockService(mockRepo)

		// Create a context with cancel
		_, cancel := context.WithCancel(context.Background())

		// Start listening in background
		go func() {
			err := service.ListenForChanges()
			if err != nil {
				t.Errorf("ListenForChanges returned error: %v", err)
			}
		}()

		// Cancel context after a short delay
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		// Send a test stock change
		testStock := domain.Stock{
			ID:        "323e4567-e89b-12d3-a456-426614174002",
			ProductID: 3,
			BranchID:  3,
			Quantity:  30,
			Reserved:  10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.stockChan <- testStock

		// Allow time for cleanup
		time.Sleep(200 * time.Millisecond)
	})

	t.Run("repository close", func(t *testing.T) {
		mockRepo := &mockStockRepository{
			stockChan: make(chan domain.Stock),
		}

		service := service.NewStockService(mockRepo)

		// Start listening in background
		go func() {
			err := service.ListenForChanges()
			if err != nil {
				t.Errorf("ListenForChanges returned error: %v", err)
			}
		}()

		// Close repository after sending one message
		testStock := domain.Stock{
			ID:        "423e4567-e89b-12d3-a456-426614174003",
			ProductID: 4,
			BranchID:  4,
			Quantity:  40,
			Reserved:  15,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.stockChan <- testStock
		time.Sleep(50 * time.Millisecond)

		err := mockRepo.Close()
		if err != nil {
			t.Errorf("Failed to close repository: %v", err)
		}

		// Allow time for cleanup
		time.Sleep(100 * time.Millisecond)
	})
}
