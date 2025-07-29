// Package logger provides logging functionality for the application
package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	// LogFile is the file handle for logging
	LogFile *os.File
	// Logger is the logger instance
	Logger *log.Logger
)

// Init initializes the logger with file output
func Init() error {
	// Use relative path for tests, absolute path for production
	logDir := "logs"
	if os.Getenv("TESTING") == "" {
		logDir = "/app/logs"
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	currentTime := time.Now().Format("2006-01-02")
	logPath := filepath.Join(logDir, fmt.Sprintf("stock-consolidation-%s.log", currentTime))

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	LogFile = file
	Logger = log.New(file, "", log.LstdFlags)

	return nil
}

// Close closes the log file
func Close() {
	if LogFile != nil {
		if err := LogFile.Close(); err != nil {
			log.Printf("Warning: Failed to close log file: %v", err)
		}
	}
}

// Info logs an informational message
func Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log.Print(msg) // Console output
	if Logger != nil {
		Logger.Print(msg)                      // File output
		if err := LogFile.Sync(); err != nil { // Force write to disk
			log.Printf("Warning: Failed to sync log file: %v", err)
		}
	}
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log.Print("ERROR: " + msg) // Console output
	if Logger != nil {
		Logger.Print("ERROR: " + msg)          // File output
		if err := LogFile.Sync(); err != nil { // Force write to disk
			log.Printf("Warning: Failed to sync log file: %v", err)
		}
	}
}

// Fatal logs a fatal message and exits the application
func Fatal(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if Logger != nil {
		Logger.Print("FATAL: " + msg)          // File output
		if err := LogFile.Sync(); err != nil { // Force write to disk
			log.Printf("Warning: Failed to sync log file: %v", err)
		}
	}
	log.Fatal("FATAL: " + msg) // Console output and exit
}
