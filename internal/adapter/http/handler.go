package http

import (
	"stock-consolidation/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/health", healthCheck)
}

func healthCheck(c *fiber.Ctx) error {
	logger.Info("Health check endpoint accessed")
	return c.JSON(fiber.Map{
		"status": "healthy",
	})
}
