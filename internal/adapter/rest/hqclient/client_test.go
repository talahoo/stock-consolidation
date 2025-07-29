package hqclient_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

			// Read and check request body as raw JSON
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Failed to read request body: %v", err)
			}

			// Check that the body contains expected fields
			bodyStr := string(body)
			if !strings.Contains(bodyStr, `"product_id":1`) ||
				!strings.Contains(bodyStr, `"branch_id":1`) ||
				!strings.Contains(bodyStr, `"quantity":10`) {
				t.Error("Request body does not contain expected fields")
			}

			// Send response
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(map[string]string{"status": "success"}); err != nil {
				t.Errorf("Failed to encode response: %v", err)
			}
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
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"}); err != nil {
				t.Errorf("Failed to encode error response: %v", err)
			}
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
