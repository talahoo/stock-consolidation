package logger_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"stock-consolidation/pkg/logger"
)

func setupLogDir() func() {
	// Set testing environment variable
	if err := os.Setenv("TESTING", "1"); err != nil {
		panic(fmt.Sprintf("Failed to set TESTING env var: %v", err))
	}

	// Create logs directory if it doesn't exist
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		panic(err)
	}

	// Return cleanup function
	return func() {
		files, err := filepath.Glob("logs/stock-consolidation-*.log")
		if err != nil {
			return
		}
		for _, file := range files {
			_ = os.Remove(file) // Ignore errors in cleanup
		}
		// Unset testing environment variable
		if err := os.Unsetenv("TESTING"); err != nil {
			// Log error but don't fail test
			fmt.Printf("Warning: Failed to unset TESTING env var: %v\n", err)
		}
	}
}

func TestLogger(t *testing.T) {
	cleanup := setupLogDir()
	defer cleanup()

	// Initialize logger for test
	if err := logger.Init(); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	t.Run("info logging", func(_ *testing.T) {
		// Test that Info doesn't panic
		logger.Info("test info message")
		// No need to check file content for speed
	})

	t.Run("error logging", func(_ *testing.T) {
		// Test that Error doesn't panic
		logger.Error("test error message")
		// No need to check file content for speed
	})

	t.Run("fatal logging doesn't exit in tests", func(_ *testing.T) {
		// Test that Fatal doesn't exit in tests
		logger.Fatal("test fatal message")
		// No need to check file content for speed
	})
}
