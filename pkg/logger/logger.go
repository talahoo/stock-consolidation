package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	LogFile *os.File
	Logger  *log.Logger
)

func Init() error {
	logDir := "/app/logs"
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

func Close() {
	if LogFile != nil {
		LogFile.Close()
	}
}

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

func Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log.Print("ERROR: " + msg) // Console output
	if Logger != nil {
		Logger.Print("ERROR: " + msg) // File output
		LogFile.Sync()                // Force write to disk
	}
}

func Fatal(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if Logger != nil {
		Logger.Print("FATAL: " + msg) // File output
		LogFile.Sync()                // Force write to disk
	}
	log.Fatal("FATAL: " + msg) // Console output and exit
}
