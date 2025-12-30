package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// setupTestApp creates a test Fiber app with API routes
func setupTestApp() *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Register API routes for testing
	api := app.Group("/api/v1")

	// Environment routes
	api.Get("/environments", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"data": []fiber.Map{
				{
					"id":                "550e8400-e29b-41d4-a716-446655440001",
					"name":              "alpha",
					"url":               "https://alpha.example.com",
					"approval_required": false,
					"is_production":     false,
					"soak_period_hours": 0,
				},
				{
					"id":                "550e8400-e29b-41d4-a716-446655440002",
					"name":              "beta",
					"url":               "https://beta.example.com",
					"approval_required": true,
					"is_production":     false,
					"soak_period_hours": 0,
				},
			},
		})
	})

	api.Get("/environments/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "550e8400-e29b-41d4-a716-446655440001" {
			return c.JSON(fiber.Map{
				"id":                "550e8400-e29b-41d4-a716-446655440001",
				"name":              "alpha",
				"url":               "https://alpha.example.com",
				"approval_required": false,
				"is_production":     false,
				"soak_period_hours": 0,
			})
		}
		return c.Status(404).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "Environment not found",
			},
		})
	})

	// Release routes
	api.Get("/releases", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"data": []fiber.Map{
				{
					"id":                    "550e8400-e29b-41d4-a716-446655440010",
					"version":               "1.0.0",
					"git_sha":               "abc123def456789012345678901234567890abcd",
					"git_sha_short":         "abc123d",
					"docker_image_backend":  "ghcr.io/example/backend:1.0.0",
					"docker_image_frontend": "ghcr.io/example/frontend:1.0.0",
					"created_by":            "ci-system",
				},
			},
			"pagination": fiber.Map{
				"page":        1,
				"limit":       20,
				"total":       1,
				"total_pages": 1,
			},
		})
	})

	api.Post("/releases", func(c *fiber.Ctx) error {
		var body map[string]interface{}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "INVALID_REQUEST",
					"message": "Invalid JSON body",
				},
			})
		}

		// Check required fields
		required := []string{"version", "git_sha", "docker_image_backend", "docker_image_frontend", "created_by"}
		for _, field := range required {
			if _, ok := body[field]; !ok {
				return c.Status(400).JSON(fiber.Map{
					"error": fiber.Map{
						"code":    "VALIDATION_ERROR",
						"message": "Missing required field: " + field,
					},
				})
			}
		}

		return c.Status(201).JSON(fiber.Map{
			"id":                    "550e8400-e29b-41d4-a716-446655440020",
			"version":               body["version"],
			"git_sha":               body["git_sha"],
			"git_sha_short":         body["git_sha"].(string)[:7],
			"docker_image_backend":  body["docker_image_backend"],
			"docker_image_frontend": body["docker_image_frontend"],
			"created_by":            body["created_by"],
		})
	})

	// Deployment routes
	api.Get("/deployments", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"data": []fiber.Map{
				{
					"id":             "550e8400-e29b-41d4-a716-446655440030",
					"environment_id": "550e8400-e29b-41d4-a716-446655440001",
					"release_id":     "550e8400-e29b-41d4-a716-446655440010",
					"status":         "success",
					"triggered_by":   "ci-system",
					"trigger_type":   "automatic",
				},
			},
			"pagination": fiber.Map{
				"page":        1,
				"limit":       20,
				"total":       1,
				"total_pages": 1,
			},
		})
	})

	api.Post("/deployments", func(c *fiber.Ctx) error {
		var body map[string]interface{}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "INVALID_REQUEST",
					"message": "Invalid JSON body",
				},
			})
		}

		return c.Status(201).JSON(fiber.Map{
			"id":             "550e8400-e29b-41d4-a716-446655440031",
			"environment_id": body["environment_id"],
			"release_id":     body["release_id"],
			"status":         "pending",
			"triggered_by":   body["triggered_by"],
			"trigger_type":   body["trigger_type"],
		})
	})

	api.Get("/deployments/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "550e8400-e29b-41d4-a716-446655440030" {
			return c.JSON(fiber.Map{
				"id":             "550e8400-e29b-41d4-a716-446655440030",
				"environment_id": "550e8400-e29b-41d4-a716-446655440001",
				"release_id":     "550e8400-e29b-41d4-a716-446655440010",
				"status":         "success",
				"triggered_by":   "ci-system",
				"trigger_type":   "automatic",
				"environment": fiber.Map{
					"name": "alpha",
					"url":  "https://alpha.example.com",
				},
				"release": fiber.Map{
					"version": "1.0.0",
				},
				"health_checks": []fiber.Map{
					{
						"check_type": "liveness",
						"status":     "passing",
					},
				},
			})
		}
		return c.Status(404).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "Deployment not found",
			},
		})
	})

	api.Post("/deployments/:id/rollback", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "550e8400-e29b-41d4-a716-446655440030" {
			var body map[string]interface{}
			if err := c.BodyParser(&body); err != nil {
				return c.Status(400).JSON(fiber.Map{
					"error": fiber.Map{
						"code":    "INVALID_REQUEST",
						"message": "Invalid JSON body",
					},
				})
			}

			return c.Status(201).JSON(fiber.Map{
				"id":             "550e8400-e29b-41d4-a716-446655440032",
				"environment_id": "550e8400-e29b-41d4-a716-446655440001",
				"release_id":     "550e8400-e29b-41d4-a716-446655440010",
				"status":         "pending",
				"triggered_by":   body["triggered_by"],
				"trigger_type":   "rollback",
				"rollback_of_id": id,
			})
		}
		return c.Status(404).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "Deployment not found",
			},
		})
	})

	return app
}

// TestListEnvironments tests GET /api/v1/environments
func TestListEnvironments(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/environments", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		t.Fatal("Response should contain 'data' array")
	}

	if len(data) != 2 {
		t.Errorf("Expected 2 environments, got %d", len(data))
	}
}

// TestGetEnvironment tests GET /api/v1/environments/:id
func TestGetEnvironment(t *testing.T) {
	app := setupTestApp()

	t.Run("existing environment", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/environments/550e8400-e29b-41d4-a716-446655440001", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if result["name"] != "alpha" {
			t.Errorf("Expected name 'alpha', got %v", result["name"])
		}
	})

	t.Run("non-existing environment", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/environments/non-existent-id", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

// TestListReleases tests GET /api/v1/releases
func TestListReleases(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/releases", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if _, ok := result["data"]; !ok {
		t.Error("Response should contain 'data'")
	}

	if _, ok := result["pagination"]; !ok {
		t.Error("Response should contain 'pagination'")
	}
}

// TestCreateRelease tests POST /api/v1/releases
func TestCreateRelease(t *testing.T) {
	app := setupTestApp()

	t.Run("valid release", func(t *testing.T) {
		payload := `{
			"version": "1.1.0",
			"git_sha": "def456abc789012345678901234567890abcdefg",
			"docker_image_backend": "ghcr.io/example/backend:1.1.0",
			"docker_image_frontend": "ghcr.io/example/frontend:1.1.0",
			"created_by": "ci-system"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/releases", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}
	})

	t.Run("missing required field", func(t *testing.T) {
		payload := `{
			"version": "1.1.0"
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/releases", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

// TestListDeployments tests GET /api/v1/deployments
func TestListDeployments(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/deployments", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// TestCreateDeployment tests POST /api/v1/deployments
func TestCreateDeployment(t *testing.T) {
	app := setupTestApp()

	payload := `{
		"environment_id": "550e8400-e29b-41d4-a716-446655440001",
		"release_id": "550e8400-e29b-41d4-a716-446655440010",
		"triggered_by": "test-user",
		"trigger_type": "manual"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/deployments", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if result["status"] != "pending" {
		t.Errorf("Expected status 'pending', got %v", result["status"])
	}
}

// TestGetDeployment tests GET /api/v1/deployments/:id
func TestGetDeployment(t *testing.T) {
	app := setupTestApp()

	t.Run("existing deployment with details", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/deployments/550e8400-e29b-41d4-a716-446655440030", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		// Should include related entities
		if _, ok := result["environment"]; !ok {
			t.Error("Response should include environment")
		}
		if _, ok := result["release"]; !ok {
			t.Error("Response should include release")
		}
		if _, ok := result["health_checks"]; !ok {
			t.Error("Response should include health_checks")
		}
	})

	t.Run("non-existing deployment", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/deployments/non-existent-id", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

// TestRollbackDeployment tests POST /api/v1/deployments/:id/rollback
func TestRollbackDeployment(t *testing.T) {
	app := setupTestApp()

	t.Run("successful rollback", func(t *testing.T) {
		payload := `{"triggered_by": "admin-user"}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/deployments/550e8400-e29b-41d4-a716-446655440030/rollback", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if result["trigger_type"] != "rollback" {
			t.Errorf("Expected trigger_type 'rollback', got %v", result["trigger_type"])
		}
		if result["rollback_of_id"] != "550e8400-e29b-41d4-a716-446655440030" {
			t.Error("rollback_of_id should reference original deployment")
		}
	})

	t.Run("non-existing deployment", func(t *testing.T) {
		payload := `{"triggered_by": "admin-user"}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/deployments/non-existent-id/rollback", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

// TestAPIResponseFormat verifies API response format consistency
func TestAPIResponseFormat(t *testing.T) {
	app := setupTestApp()

	t.Run("list response has data array", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/environments", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)

		if _, ok := result["data"].([]interface{}); !ok {
			t.Error("List response should have 'data' as array")
		}
	})

	t.Run("error response has error object", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/environments/non-existent", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)

		errorObj, ok := result["error"].(map[string]interface{})
		if !ok {
			t.Error("Error response should have 'error' object")
		}
		if _, ok := errorObj["code"]; !ok {
			t.Error("Error object should have 'code'")
		}
		if _, ok := errorObj["message"]; !ok {
			t.Error("Error object should have 'message'")
		}
	})
}
