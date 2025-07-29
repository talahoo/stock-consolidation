package http_test

import (
	"io"
	"net/http/httptest"
	"testing"

	"stock-consolidation/internal/adapter/http"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	app := fiber.New()
	http.SetupRoutes(app)

	t.Run("health endpoint returns 200", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, `{"status":"healthy"}`, string(body))
	})

	t.Run("not found endpoint returns 404", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/not-exist", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}
