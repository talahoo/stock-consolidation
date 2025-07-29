package hqclient_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"stock-consolidation/internal/adapter/rest/hqclient"
	"stock-consolidation/internal/core/domain"
	"stock-consolidation/pkg/config"
)

func TestHQClient_SendStockChange(t *testing.T) {
	testTime := time.Date(2025, 7, 29, 0, 0, 0, 0, time.UTC)
	t.Run("successful stock update", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check request method
			if r.Method != http.MethodPost {
				t.Errorf("Expected POST request, got %s", r.Method)
			}

			// Check authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "Basic dXNlcjpwYXNz" {
				t.Errorf("Expected Authorization header Basic dXNlcjpwYXNz, got %s", authHeader)
			}

			// Check content type
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}

			// Decode request body
			var stock domain.Stock
			if err := json.NewDecoder(r.Body).Decode(&stock); err != nil {
				t.Errorf("Failed to decode request body: %v", err)
			}

			// Check times
			if !stock.CreatedAt.Equal(testTime) || !stock.UpdatedAt.Equal(testTime) {
				t.Error("Time fields do not match")
			}

			// Send response
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		}))
		defer server.Close()

		// Create client with test server URL
		cfg := &config.Config{
			HQEndPoint:           server.URL,
			HQBasicAuthorization: "Basic dXNlcjpwYXNz",
		}
		client := hqclient.NewHQClient(cfg)

		// Create test stock data
		stock := domain.Stock{
			ProductID: 1,
			BranchID:  1,
			Quantity:  10,
			CreatedAt: testTime,
			UpdatedAt: testTime,
		}

		// Send update
		err := client.SendStockChange(context.Background(), stock)
		if err != nil {
			t.Errorf("SendStockChange() error = %v", err)
		}
	})

	t.Run("server error response", func(t *testing.T) {
		// Create test server that returns an error
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		}))
		defer server.Close()

		// Create client with test server URL
		cfg := &config.Config{
			HQEndPoint:           server.URL,
			HQBasicAuthorization: "Basic dXNlcjpwYXNz",
		}
		client := hqclient.NewHQClient(cfg)

		// Create test stock data
		stock := domain.Stock{
			ProductID: 1,
			BranchID:  1,
			Quantity:  10,
			CreatedAt: testTime,
			UpdatedAt: testTime,
		}

		// Send update and expect error
		err := client.SendStockChange(context.Background(), stock)
		if err == nil {
			t.Error("SendStockChange() expected error for server error, got nil")
		}
	})

	t.Run("connection error", func(t *testing.T) {
		// Create client with invalid URL
		cfg := &config.Config{
			HQEndPoint:           "http://invalid-url",
			HQBasicAuthorization: "Basic dXNlcjpwYXNz",
		}
		client := hqclient.NewHQClient(cfg)

		// Create test stock data
		stock := domain.Stock{
			ProductID: 1,
			BranchID:  1,
			Quantity:  10,
			CreatedAt: testTime,
			UpdatedAt: testTime,
		}

		// Send update and expect error
		err := client.SendStockChange(context.Background(), stock)
		if err == nil {
			t.Error("SendStockChange() expected error for invalid URL, got nil")
		}
	})
}
