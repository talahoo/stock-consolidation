package service

import (
	"context"
	"fmt"
	"stock-consolidation/internal/adapter/rest/hqclient"
	"stock-consolidation/internal/core/port"
	"stock-consolidation/pkg/config"
	"stock-consolidation/pkg/logger"
)

type StockService struct {
	repo   port.StockRepository
	client *hqclient.HQClient
}

func NewStockService(repo port.StockRepository) *StockService {
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config: %v", err)
	}

	return &StockService{
		repo:   repo,
		client: hqclient.NewHQClient(cfg),
	}
}

func (s *StockService) ListenForChanges() error {
	logger.Info("Starting StockService.ListenForChanges()")
	ctx := context.Background()
	stockChan, err := s.repo.ListenForChanges(ctx)
	if err != nil {
		logger.Error("Failed to start listening for changes: %v", err)
		return fmt.Errorf("failed to start listening: %v", err)
	}

	logger.Info("Successfully started listening for stock changes")
	for stock := range stockChan {
		logger.Info("Processing stock change notification: ProductID=%d, BranchID=%d", stock.ProductID, stock.BranchID)

		if err := s.client.SendStockChange(ctx, stock); err != nil {
			logger.Error("Failed to send stock change to HQ: %v", err)
			continue
		}

		logger.Info("Successfully sent stock change for product %d in branch %d", stock.ProductID, stock.BranchID)
	}
	logger.Info("Stopped listening for stock changes")
	return nil
}
