package postgres_test

import (
	"context"
	"testing"
	"time"

	"stock-consolidation/internal/adapter/db/postgres"
	"stock-consolidation/pkg/config"
)

func TestPostgresListener(t *testing.T) {
	cfg := &config.Config{
		DBHost:               "localhost",
		DBPort:               "5432",
		DBUser:               "admin",
		DBPassword:           "admin",
		DBName:               "stockdb",
		ServicePort:          "3000",
		HQEndPoint:           "http://localhost:8080",
		HQBasicAuthorization: "Basic dXNlcjpwYXNz",
	}

	t.Run("success connection", func(t *testing.T) {
		listener, err := postgres.NewListener(cfg)
		if err != nil {
			t.Fatalf("Failed to create listener: %v", err)
		}
		defer listener.Close()

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
	})

	t.Run("invalid configuration", func(t *testing.T) {
		invalidCfg := &config.Config{
			DBHost:               "invalid-host",
			DBPort:               "invalid-port",
			DBUser:               "invalid-user",
			DBPassword:           "invalid-password",
			DBName:               "invalid-db",
			ServicePort:          "3000",
			HQEndPoint:           "http://localhost:8080",
			HQBasicAuthorization: "Basic dXNlcjpwYXNz",
		}

		_, err := postgres.NewListener(invalidCfg)
		if err == nil {
			t.Error("Expected error with invalid configuration, got nil")
		}
	})
}
