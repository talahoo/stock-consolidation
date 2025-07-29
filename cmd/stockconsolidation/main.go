package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"

	"stock-consolidation/internal/adapter/db/postgres"
	"stock-consolidation/internal/adapter/http"
	"stock-consolidation/internal/service"
	"stock-consolidation/pkg/config"
	"stock-consolidation/pkg/logger"
)

func main() {
	// Initialize logger

	if err := logger.Init(); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Close()

	logger.Info("Starting Stock Consolidation Service...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config: %v", err)
		return
	}

	// Initialize PostgreSQL listener
	listener, err := postgres.NewListener(cfg)
	if err != nil {
		logger.Fatal("Failed to create PostgreSQL listener: %v", err)
		return
	}

	// Initialize services
	stockService := service.NewStockService(listener)

	// Initialize Fiber app with custom config
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Setup routes
	http.SetupRoutes(app)

	// Start listening for stock changes in background
	go func() {
		defer listener.Close()
		if err := stockService.ListenForChanges(); err != nil {
			logger.Error("Error listening for changes: %v", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start HTTP server in a goroutine
	go func() {
		addr := fmt.Sprintf("0.0.0.0:%s", cfg.ServicePort)
		logger.Info("Starting HTTP server on %s", addr)
		if err := app.Listen(addr); err != nil {
			logger.Fatal("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit

	// Graceful shutdown
	logger.Info("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		logger.Fatal("Error shutting down server: %v", err)
	}

}
