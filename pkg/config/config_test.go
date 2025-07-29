package config_test

import (
	"os"
	"testing"

	"stock-consolidation/pkg/config"
)

// setEnv sets an environment variable and fails the test if it fails
func setEnv(t *testing.T, key, value string) {
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("Failed to set environment variable %s: %v", key, err)
	}
}

func TestLoadConfig(t *testing.T) {
	t.Run("success load config", func(t *testing.T) {
		// Set environment variables
		setEnv(t, "DB_HOST", "localhost")
		setEnv(t, "DB_PORT", "5432")
		setEnv(t, "DB_USER", "admin")
		setEnv(t, "DB_PASSWORD", "admin")
		setEnv(t, "DB_NAME", "stockdb")
		setEnv(t, "SERVICE_PORT", "3000")
		setEnv(t, "HQ_END_POINT", "http://localhost:8080")
		setEnv(t, "HQ_BASIC_AUTHORIZATION", "Basic dXNlcjpwYXNz")

		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}

		if cfg.DBHost != "localhost" {
			t.Errorf("LoadConfig() DBHost = %v, want %v", cfg.DBHost, "localhost")
		}
		if cfg.DBPort != "5432" {
			t.Errorf("LoadConfig() DBPort = %v, want %v", cfg.DBPort, "5432")
		}
		if cfg.DBUser != "admin" {
			t.Errorf("LoadConfig() DBUser = %v, want %v", cfg.DBUser, "admin")
		}
		if cfg.DBPassword != "admin" {
			t.Errorf("LoadConfig() DBPassword = %v, want %v", cfg.DBPassword, "admin")
		}
		if cfg.DBName != "stockdb" {
			t.Errorf("LoadConfig() DBName = %v, want %v", cfg.DBName, "stockdb")
		}
		if cfg.ServicePort != "3000" {
			t.Errorf("LoadConfig() ServicePort = %v, want %v", cfg.ServicePort, "3000")
		}
		if cfg.HQEndPoint != "http://localhost:8080" {
			t.Errorf("LoadConfig() HQEndPoint = %v, want %v", cfg.HQEndPoint, "http://localhost:8080")
		}
		if cfg.HQBasicAuthorization != "Basic dXNlcjpwYXNz" {
			t.Errorf("LoadConfig() HQBasicAuthorization = %v, want %v", cfg.HQBasicAuthorization, "Basic dXNlcjpwYXNz")
		}
	})

	t.Run("missing required DB_HOST", func(t *testing.T) {
		// Clear environment variables
		os.Clearenv()

		// Set other variables except DB_HOST
		setEnv(t, "DB_PORT", "5432")
		setEnv(t, "DB_USER", "admin")
		setEnv(t, "DB_PASSWORD", "admin")
		setEnv(t, "DB_NAME", "stockdb")
		setEnv(t, "SERVICE_PORT", "3000")
		setEnv(t, "HQ_END_POINT", "http://localhost:8080")
		setEnv(t, "HQ_BASIC_AUTHORIZATION", "Basic dXNlcjpwYXNz")

		_, err := config.Load()
		if err == nil {
			t.Error("LoadConfig() expected error for missing DB_HOST, got nil")
		}
		if err.Error() != "DB_HOST is required" {
			t.Errorf("LoadConfig() error = %v, want %v", err, "DB_HOST is required")
		}
	})

	t.Run("missing required SERVICE_PORT", func(t *testing.T) {
		// Clear environment variables
		os.Clearenv()

		// Set required variables except SERVICE_PORT
		setEnv(t, "DB_HOST", "localhost")
		setEnv(t, "DB_PORT", "5432")
		setEnv(t, "DB_USER", "admin")
		setEnv(t, "DB_PASSWORD", "admin")
		setEnv(t, "DB_NAME", "stockdb")
		setEnv(t, "HQ_END_POINT", "http://localhost:8080")
		setEnv(t, "HQ_BASIC_AUTHORIZATION", "Basic dXNlcjpwYXNz")

		_, err := config.Load()
		if err == nil {
			t.Error("LoadConfig() expected error for missing SERVICE_PORT, got nil")
		}
		if err.Error() != "SERVICE_PORT is required" {
			t.Errorf("LoadConfig() error = %v, want %v", err, "SERVICE_PORT is required")
		}
	})

	t.Run("missing required DB_NAME", func(t *testing.T) {
		// Clear environment variables
		os.Clearenv()

		// Set other required variables except DB_NAME
		setEnv(t, "DB_HOST", "localhost")
		setEnv(t, "DB_PORT", "5432")
		setEnv(t, "DB_USER", "admin")
		setEnv(t, "DB_PASSWORD", "admin")
		setEnv(t, "SERVICE_PORT", "3000")
		setEnv(t, "HQ_END_POINT", "http://localhost:8080")
		setEnv(t, "HQ_BASIC_AUTHORIZATION", "Basic dXNlcjpwYXNz")

		_, err := config.Load()
		if err == nil {
			t.Error("LoadConfig() expected error for missing DB_NAME, got nil")
		}
		if err.Error() != "DB_NAME is required" {
			t.Errorf("LoadConfig() error = %v, want %v", err, "DB_NAME is required")
		}
	})
}
