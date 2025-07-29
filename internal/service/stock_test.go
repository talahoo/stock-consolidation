package service_test

import (
	"context"
	"fmt"
	"os"
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

func setupTestEnv() func() {
	// Save original environment variables
	envVars := []string{
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"SERVICE_PORT",
		"HQ_END_POINT",
		"HQ_BASIC_AUTHORIZATION",
	}

	originalEnvVars := make(map[string]string)
	for _, env := range envVars {
		originalEnvVars[env] = os.Getenv(env)
	}

	// Set test environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "admin")
	os.Setenv("DB_PASSWORD", "admin123")
	os.Setenv("DB_NAME", "stockdb")
	os.Setenv("SERVICE_PORT", "3000")
	os.Setenv("HQ_END_POINT", "http://localhost:8085/stock")
	os.Setenv("HQ_BASIC_AUTHORIZATION", "Basic dXNlcjpwYXNz")

	// Return cleanup function
	return func() {
		for env, value := range originalEnvVars {
			if value != "" {
				os.Setenv(env, value)
			} else {
				os.Unsetenv(env)
			}
		}
	}
}

func TestStockService_ListenForChanges(t *testing.T) {
	cleanup := setupTestEnv()
	defer cleanup()
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
