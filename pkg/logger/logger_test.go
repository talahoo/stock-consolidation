package logger_test

import (
	"testing"

	"stock-consolidation/pkg/logger"
)

func TestLogger(t *testing.T) {
	t.Run("info logging", func(t *testing.T) {
		logger.Info("test info message")
		// Since logging is side effect, we just verify it doesn't panic
	})

	t.Run("error logging", func(t *testing.T) {
		logger.Error("test error message")
		// Since logging is side effect, we just verify it doesn't panic
	})

	t.Run("fatal logging doesn't exit in tests", func(t *testing.T) {
		logger.Fatal("test fatal message")
		// Since we're in test mode, Fatal doesn't actually exit
	})
}
