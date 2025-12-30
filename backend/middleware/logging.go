package middleware

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Logger middleware provides structured JSON logging for all HTTP requests
// It logs request duration, status code, method, path, and correlation ID
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer
		start := time.Now()

		// Get correlation ID
		correlationID := GetCorrelationID(c)

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Create log attributes
		attrs := []slog.Attr{
			slog.String("correlation_id", correlationID),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.String("query", c.Context().QueryArgs().String()),
			slog.String("ip", c.IP()),
			slog.String("user_agent", c.Get("User-Agent")),
			slog.Int("status", c.Response().StatusCode()),
			slog.Duration("duration_ms", duration),
			slog.String("protocol", c.Protocol()),
		}

		// Add host information
		if c.Hostname() != "" {
			attrs = append(attrs, slog.String("host", c.Hostname()))
		}

		// Log based on status code
		status := c.Response().StatusCode()
		ctx := context.Background()
		switch {
		case status >= 500:
			slog.LogAttrs(ctx, slog.LevelError, "Server error", attrs...)
		case status >= 400:
			slog.LogAttrs(ctx, slog.LevelWarn, "Client error", attrs...)
		case status >= 300:
			slog.LogAttrs(ctx, slog.LevelInfo, "Redirect", attrs...)
		default:
			slog.LogAttrs(ctx, slog.LevelInfo, "Request completed", attrs...)
		}

		return err
	}
}

// NewLogger initializes a structured JSON logger
// It returns a configured slog.Logger that outputs JSON to stdout
func NewLogger(level string) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	return logger
}
