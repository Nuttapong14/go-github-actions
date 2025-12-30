package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/Nuttapong14/go-github-actions/backend/config"
	"github.com/Nuttapong14/go-github-actions/backend/database"
	"github.com/Nuttapong14/go-github-actions/backend/handlers"
	"github.com/Nuttapong14/go-github-actions/backend/middleware"
)

func main() {
	// Parse command line flags
	autoMigrate := flag.Bool("auto-migrate", false, "Run database auto-migration on startup")
	flag.Parse()

	// Initialize logger
	logger := middleware.NewLogger("info")
	slog.SetDefault(logger)

	slog.Info("Starting CI/CD Training Application Backend")

	// Load configuration
	cfg := config.Load()
	slog.Info("Configuration loaded",
		slog.String("environment", cfg.Environment),
		slog.String("port", cfg.Port),
		slog.String("log_level", cfg.LogLevel),
	)

	// Connect to database
	db, err := database.New(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("Failed to close database connection", slog.String("error", err.Error()))
		}
	}()

	slog.Info("Database connection established")

	// Auto-migrate if flag is set
	if *autoMigrate {
		slog.Info("Running auto-migration")
		// TODO: Add model auto-migration when models are ready
		// if err := db.AutoMigrate(&models.Environment{}, &models.Release{}, &models.Deployment{}, &models.HealthCheck{}); err != nil {
		// 	slog.Error("Auto-migration failed", slog.String("error", err.Error()))
		// 	os.Exit(1)
		// }
		slog.Info("Auto-migration completed")
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:               "cicd-training-backend",
		DisableStartupMessage: false,
		EnablePrintRoutes:     cfg.IsDevelopment(),
	})

	// Setup global middleware
	app.Use(recover.New())
	app.Use(middleware.CorrelationID())
	app.Use(middleware.Logger())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-Correlation-ID",
	}))

	// Setup Prometheus metrics
	prometheus := handlers.SetupMetrics(app, "cicd_training")
	app.Use(prometheus.Middleware)

	// Health check routes
	app.Get("/health/live", handlers.Liveness)
	app.Get("/health/ready", handlers.ReadinessWithChecker(db.Health))

	// API routes (v1)
	api := app.Group("/api/v1")
	setupAPIRoutes(api)

	// Start server in goroutine
	go func() {
		addr := cfg.GetListenAddr()
		slog.Info("Starting HTTP server", slog.String("address", addr))
		if err := app.Listen(addr); err != nil {
			slog.Error("Server error", slog.String("error", err.Error()))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		slog.Error("Server shutdown error", slog.String("error", err.Error()))
	}

	slog.Info("Server stopped")
}

// setupAPIRoutes configures the API v1 routes
func setupAPIRoutes(api fiber.Router) {
	// Placeholder routes for future implementation
	api.Get("/environments", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "List environments - coming soon",
		})
	})

	api.Get("/releases", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "List releases - coming soon",
		})
	})

	api.Get("/deployments", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "List deployments - coming soon",
		})
	})
}
