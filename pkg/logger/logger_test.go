package logger_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"stock-consolidation/pkg/logger"
)

func setupLogDir() func() {
	// Set testing environment variable
	os.Setenv("TESTING", "1")

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
		os.Unsetenv("TESTING")
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
	t.Run("info logging", func(t *testing.T) {
		logger.Info("test info message")

		// Verify log file was created and contains the message
		files, err := filepath.Glob("logs/stock-consolidation-*.log")
		if err != nil {
			t.Fatalf("Failed to list log files: %v", err)
		}
		if len(files) == 0 {
			t.Fatal("No log file was created")
		}

		content, err := os.ReadFile(files[0])
		if err != nil {
			t.Fatalf("Failed to read log file: %v", err)
		}
		if !strings.Contains(string(content), "test info message") {
			t.Error("Log file does not contain expected info message")
		}
	})

	t.Run("error logging", func(t *testing.T) {
		logger.Error("test error message")

		files, err := filepath.Glob("logs/stock-consolidation-*.log")
		if err != nil {
			t.Fatalf("Failed to list log files: %v", err)
		}
		if len(files) == 0 {
			t.Fatal("No log file was created")
		}

		content, err := os.ReadFile(files[0])
		if err != nil {
			t.Fatalf("Failed to read log file: %v", err)
		}
		if !strings.Contains(string(content), "ERROR: test error message") {
			t.Error("Log file does not contain expected error message")
		}
	})

	t.Run("fatal logging doesn't exit in tests", func(t *testing.T) {
		logger.Fatal("test fatal message")

		files, err := filepath.Glob("logs/stock-consolidation-*.log")
		if err != nil {
			t.Fatalf("Failed to list log files: %v", err)
		}
		if len(files) == 0 {
			t.Fatal("No log file was created")
		}

		content, err := os.ReadFile(files[0])
		if err != nil {
			t.Fatalf("Failed to read log file: %v", err)
		}
		if !strings.Contains(string(content), "FATAL: test fatal message") {
			t.Error("Log file does not contain expected fatal message")
		}
	})
}
