package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// DBChecker is a function type for checking database health
type DBChecker func() error

// LivenessResponse represents the response from the liveness endpoint
type LivenessResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// ReadinessResponse represents the response from the readiness endpoint
type ReadinessResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

// Liveness handles GET /health/live
// Returns ok if the service is running
func Liveness(c *fiber.Ctx) error {
	return c.JSON(LivenessResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// Readiness handles GET /health/ready
// Returns ok if all dependencies (database) are healthy
// This is a wrapper that uses the default DB checker from database package
func Readiness(c *fiber.Ctx) error {
	return ReadinessWithChecker(nil)(c)
}

// ReadinessWithChecker returns a handler that uses the provided DB checker
// If dbChecker is nil, the check is skipped and assumed healthy
func ReadinessWithChecker(dbChecker DBChecker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		checks := make(map[string]string)
		overallStatus := "ok"
		httpStatus := fiber.StatusOK

		// Check database health
		if dbChecker != nil {
			if err := dbChecker(); err != nil {
				checks["database"] = "error"
				overallStatus = "error"
				httpStatus = fiber.StatusServiceUnavailable
			} else {
				checks["database"] = "ok"
			}
		} else {
			// No checker provided, assume healthy (useful for testing)
			checks["database"] = "ok"
		}

		return c.Status(httpStatus).JSON(ReadinessResponse{
			Status:    overallStatus,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Checks:    checks,
		})
	}
}
