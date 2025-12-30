package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nuttapong14/go-github-actions/backend/config"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Clear any existing env vars
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("PORT")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("ENV")

	cfg := config.Load()

	assert.Equal(t, "postgres://postgres:postgres@localhost:5432/cicd_training?sslmode=disable", cfg.DatabaseURL)
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "development", cfg.Environment)
}

func TestLoadConfig_FromEnv(t *testing.T) {
	// Set custom env vars
	os.Setenv("DATABASE_URL", "postgres://custom:pass@host:5432/db")
	os.Setenv("PORT", "9090")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("ENV", "production")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("PORT")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("ENV")
	}()

	cfg := config.Load()

	assert.Equal(t, "postgres://custom:pass@host:5432/db", cfg.DatabaseURL)
	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, "production", cfg.Environment)
}

func TestConfig_IsDevelopment(t *testing.T) {
	tests := []struct {
		env      string
		expected bool
	}{
		{"development", true},
		{"dev", true},
		{"local", true},
		{"production", false},
		{"prod", false},
		{"staging", false},
		{"alpha", false},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			cfg := &config.Config{Environment: tt.env}
			assert.Equal(t, tt.expected, cfg.IsDevelopment())
		})
	}
}

func TestConfig_IsProduction(t *testing.T) {
	tests := []struct {
		env      string
		expected bool
	}{
		{"production", true},
		{"prod", true},
		{"development", false},
		{"staging", false},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			cfg := &config.Config{Environment: tt.env}
			assert.Equal(t, tt.expected, cfg.IsProduction())
		})
	}
}

func TestConfig_GetListenAddr(t *testing.T) {
	cfg := &config.Config{Port: "3000"}
	assert.Equal(t, ":3000", cfg.GetListenAddr())
}

func TestConfig_Validate(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := &config.Config{
			DatabaseURL: "postgres://user:pass@host:5432/db",
			Port:        "8080",
			LogLevel:    "info",
			Environment: "development",
			Server: config.ServerConfig{
				Port: 8080,
				Host: "0.0.0.0",
			},
		}
		err := cfg.Validate()
		require.NoError(t, err)
	})

	t.Run("missing database URL", func(t *testing.T) {
		cfg := &config.Config{
			DatabaseURL: "",
			Port:        "8080",
			Server: config.ServerConfig{
				Port: 8080,
			},
		}
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "DATABASE_URL")
	})

	t.Run("invalid port", func(t *testing.T) {
		cfg := &config.Config{
			DatabaseURL: "postgres://user:pass@host:5432/db",
			Port:        "invalid",
			Server: config.ServerConfig{
				Port: 8080,
			},
		}
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "PORT")
	})
}
