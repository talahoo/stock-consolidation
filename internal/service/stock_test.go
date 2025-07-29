package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"stock-consolidation/internal/core/domain"
	"stock-consolidation/internal/service"
)

type mockStockRepository struct {
	ListenForChangesFunc func(ctx context.Context) (<-chan domain.Stock, error)
	CloseFunc            func() error
}

func (m *mockStockRepository) ListenForChanges(ctx context.Context) (<-chan domain.Stock, error) {
	if m.ListenForChangesFunc != nil {
		return m.ListenForChangesFunc(ctx)
	}
	return nil, fmt.Errorf("mock error")
}

func (m *mockStockRepository) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func TestStockService_ListenForChanges(t *testing.T) {
	t.Run("success receive stock changes", func(t *testing.T) {
		stockChan := make(chan domain.Stock)
		mockRepo := &mockStockRepository{
			ListenForChangesFunc: func(ctx context.Context) (<-chan domain.Stock, error) {
				return stockChan, nil
			},
		}

		service := service.NewStockService(mockRepo)

		// Start listening in background
		go func() {
			err := service.ListenForChanges()
			if err != nil {
				t.Errorf("ListenForChanges() error = %v", err)
			}
		}()

		// Send test stock changes
		testStocks := []domain.Stock{
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

		for _, stock := range testStocks {
			stockChan <- stock
			time.Sleep(50 * time.Millisecond) // Allow time for processing
		}

		close(stockChan)
	})

	t.Run("repository close", func(t *testing.T) {
		mockRepo := &mockStockRepository{
			ListenForChangesFunc: func(ctx context.Context) (<-chan domain.Stock, error) {
				ch := make(chan domain.Stock)
				close(ch)
				return ch, nil
			},
		}

		service := service.NewStockService(mockRepo)

		err := service.ListenForChanges()
		if err != nil {
			t.Errorf("ListenForChanges() error = %v", err)
		}
	})
}
