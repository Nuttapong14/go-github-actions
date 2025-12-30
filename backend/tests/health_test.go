package tests

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nuttapong14/go-github-actions/backend/handlers"
)

func TestLivenessHandler(t *testing.T) {
	app := fiber.New()
	app.Get("/health/live", handlers.Liveness)

	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "returns ok status",
			expectedStatus: 200,
			expectedBody: map[string]interface{}{
				"status": "ok",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/health/live", nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedBody["status"], result["status"])
			assert.NotEmpty(t, result["timestamp"])
		})
	}
}

func TestReadinessHandler(t *testing.T) {
	app := fiber.New()

	// Create a mock DB checker for testing
	mockDBCheck := func() error {
		return nil // Simulate healthy DB
	}

	app.Get("/health/ready", handlers.ReadinessWithChecker(mockDBCheck))

	t.Run("returns ok when database is healthy", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health/ready", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, "ok", result["status"])
		assert.NotEmpty(t, result["timestamp"])

		checks, ok := result["checks"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "ok", checks["database"])
	})
}

func TestReadinessHandler_DatabaseFailure(t *testing.T) {
	app := fiber.New()

	// Create a mock DB checker that fails
	mockDBCheck := func() error {
		return assert.AnError
	}

	app.Get("/health/ready", handlers.ReadinessWithChecker(mockDBCheck))

	t.Run("returns error when database is unhealthy", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health/ready", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 503, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, "error", result["status"])

		checks, ok := result["checks"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "error", checks["database"])
	})
}
