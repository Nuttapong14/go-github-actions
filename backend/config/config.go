package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all application configuration loaded from environment variables
type Config struct {
	// Simplified fields for direct access (used by tests and simple lookups)
	DatabaseURL string
	Port        string
	LogLevel    string
	Environment string

	// Structured fields for detailed configuration
	Database DatabaseConfig
	Server   ServerConfig
	Logging  LoggingConfig
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port int
	Host string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Format string
}

// Load reads configuration from environment variables and returns a Config struct
// Default values are provided for unset variables
func Load() *Config {
	cfg := &Config{
		Environment: getEnv("ENV", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Port:        getEnv("PORT", "8080"),
	}

	// Database configuration
	cfg.Database = DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "cicd_training"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Server configuration
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		port = 8080 // fallback to default port
	}
	cfg.Server = ServerConfig{
		Host: getEnv("HOST", "0.0.0.0"),
		Port: port,
	}

	// Logging configuration
	cfg.Logging = LoggingConfig{
		Level:  cfg.LogLevel,
		Format: getEnv("LOG_FORMAT", "json"),
	}

	// Set DatabaseURL from DATABASE_URL env var or construct from parts
	cfg.DatabaseURL = getEnv("DATABASE_URL", cfg.buildDatabaseURL())

	return cfg
}

// buildDatabaseURL constructs a PostgreSQL connection string from components
func (c *Config) buildDatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// Validate checks that all required configuration values are set correctly
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if _, err := strconv.Atoi(c.Port); err != nil {
		return fmt.Errorf("PORT must be a valid number, got: %s", c.Port)
	}
	if c.Database.DBName == "" && c.DatabaseURL == "" {
		return fmt.Errorf("database name cannot be empty")
	}
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535, got: %d", c.Server.Port)
	}
	return nil
}

// GetDatabaseURL returns the PostgreSQL connection string
func (c *Config) GetDatabaseURL() string {
	return c.DatabaseURL
}

// ServerAddr returns the server address in host:port format
func (c *Config) ServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// GetListenAddr returns the listen address in :port format
func (c *Config) GetListenAddr() string {
	return ":" + c.Port
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	env := strings.ToLower(c.Environment)
	return env == "development" || env == "dev" || env == "local"
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	env := strings.ToLower(c.Environment)
	return env == "production" || env == "prod"
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
