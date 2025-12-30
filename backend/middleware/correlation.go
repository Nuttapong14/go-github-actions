package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	// CorrelationIDHeader is the header name for correlation ID
	CorrelationIDHeader = "X-Correlation-ID"
	// CorrelationIDContextKey is the context key for storing correlation ID
	CorrelationIDContextKey = "correlation_id"
)

// CorrelationID middleware adds a unique correlation ID to each request
// If a correlation ID is provided in the header, it uses that value
// Otherwise, it generates a new UUID v4
func CorrelationID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get correlation ID from header or generate new one
		correlationID := c.Get(CorrelationIDHeader)
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		// Store in context for later use
		c.Locals(CorrelationIDContextKey, correlationID)

		// Set correlation ID in response header
		c.Set(CorrelationIDHeader, correlationID)

		return c.Next()
	}
}

// GetCorrelationID retrieves the correlation ID from the fiber context
func GetCorrelationID(c *fiber.Ctx) string {
	if correlationID, ok := c.Locals(CorrelationIDContextKey).(string); ok {
		return correlationID
	}
	return ""
}
