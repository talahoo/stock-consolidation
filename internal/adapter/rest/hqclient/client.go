// Package hqclient provides HTTP client functionality for communicating with HQ endpoints
package hqclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"stock-consolidation/internal/core/domain"
	"stock-consolidation/pkg/config"
	"stock-consolidation/pkg/logger"
)

// HQClient handles communication with the HQ endpoint
type HQClient struct {
	endpoint   string
	authHeader string
	httpClient *http.Client
}

// NewHQClient creates a new HQClient instance
func NewHQClient(cfg *config.Config) *HQClient {
	return &HQClient{
		endpoint:   cfg.HQEndPoint,
		authHeader: cfg.HQBasicAuthorization,
		httpClient: &http.Client{},
	}
}

// SendStockChange sends a stock change notification to the HQ endpoint
func (c *HQClient) SendStockChange(ctx context.Context, stock domain.Stock) error {
	payload, err := json.Marshal(stock)
	if err != nil {
		return fmt.Errorf("failed to marshal stock: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.authHeader)

	logger.Info("Sending stock update to HQ endpoint %s for product %d in branch %d", c.endpoint, stock.ProductID, stock.BranchID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Error("Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HQ endpoint returned error status: %d", resp.StatusCode)
	}

	return nil
}
