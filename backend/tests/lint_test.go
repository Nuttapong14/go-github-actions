package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGolangciLintConfigExists verifies the .golangci.yml configuration file exists
func TestGolangciLintConfigExists(t *testing.T) {
	// Get the project root (one level up from backend)
	projectRoot := filepath.Join("..", "..")
	configPath := filepath.Join(projectRoot, ".golangci.yml")

	_, err := os.Stat(configPath)
	require.NoError(t, err, ".golangci.yml should exist at project root")
}

// TestGolangciLintConfigHasRequiredLinters verifies essential linters are enabled
func TestGolangciLintConfigHasRequiredLinters(t *testing.T) {
	projectRoot := filepath.Join("..", "..")
	configPath := filepath.Join(projectRoot, ".golangci.yml")

	content, err := os.ReadFile(configPath)
	require.NoError(t, err, "should be able to read .golangci.yml")

	configStr := string(content)

	// Required linters that must be enabled for CI consistency
	requiredLinters := []string{
		"errcheck",    // Check for unchecked errors
		"gosec",       // Security problems
		"govet",       // Reports suspicious constructs
		"staticcheck", // Go static analysis
		"gofmt",       // Format code
		"goimports",   // Manage import lines
		"unused",      // Find unused code
	}

	for _, linter := range requiredLinters {
		assert.True(t, strings.Contains(configStr, linter),
			"Required linter %q should be configured in .golangci.yml", linter)
	}
}

// TestGolangciLintConfigSecurityRules verifies security-related rules are enabled
func TestGolangciLintConfigSecurityRules(t *testing.T) {
	projectRoot := filepath.Join("..", "..")
	configPath := filepath.Join(projectRoot, ".golangci.yml")

	content, err := os.ReadFile(configPath)
	require.NoError(t, err, "should be able to read .golangci.yml")

	configStr := string(content)

	// Important security rules from gosec
	securityRules := []string{
		"G101", // Hard coded credentials
		"G201", // SQL query construction
		"G304", // File path provided as taint input
		"G401", // Weak cryptographic primitive
	}

	for _, rule := range securityRules {
		assert.True(t, strings.Contains(configStr, rule),
			"Security rule %q should be configured in gosec settings", rule)
	}
}

// TestGolangciLintInstalled verifies golangci-lint is available
func TestGolangciLintInstalled(t *testing.T) {
	if os.Getenv("CI") == "" {
		t.Skip("Skipping golangci-lint installation check in local environment")
	}

	cmd := exec.Command("golangci-lint", "--version")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "golangci-lint should be installed: %s", string(output))
	assert.Contains(t, string(output), "golangci-lint", "output should contain version info")
}

// TestMakefileLintTargetExists verifies the Makefile has lint targets
func TestMakefileLintTargetExists(t *testing.T) {
	projectRoot := filepath.Join("..", "..")
	makefilePath := filepath.Join(projectRoot, "Makefile")

	content, err := os.ReadFile(makefilePath)
	require.NoError(t, err, "should be able to read Makefile")

	makefileStr := string(content)

	// Required make targets for linting
	requiredTargets := []string{
		"lint:",          // Main lint target
		"backend.lint:",  // Backend-specific lint
		"frontend.lint:", // Frontend-specific lint
	}

	for _, target := range requiredTargets {
		assert.True(t, strings.Contains(makefileStr, target),
			"Makefile should contain target %q", target)
	}
}

// TestMakefileTestTargetExists verifies the Makefile has test targets with coverage
func TestMakefileTestTargetExists(t *testing.T) {
	projectRoot := filepath.Join("..", "..")
	makefilePath := filepath.Join(projectRoot, "Makefile")

	content, err := os.ReadFile(makefilePath)
	require.NoError(t, err, "should be able to read Makefile")

	makefileStr := string(content)

	// Required make targets for testing
	requiredTargets := []string{
		"test:",          // Main test target
		"backend.test:",  // Backend-specific test
		"frontend.test:", // Frontend-specific test
	}

	for _, target := range requiredTargets {
		assert.True(t, strings.Contains(makefileStr, target),
			"Makefile should contain target %q", target)
	}

	// Verify coverage is configured
	assert.True(t, strings.Contains(makefileStr, "coverprofile"),
		"Makefile should use -coverprofile for coverage reporting")
}

// TestMakefileBuildTargetExists verifies the Makefile has build targets
func TestMakefileBuildTargetExists(t *testing.T) {
	projectRoot := filepath.Join("..", "..")
	makefilePath := filepath.Join(projectRoot, "Makefile")

	content, err := os.ReadFile(makefilePath)
	require.NoError(t, err, "should be able to read Makefile")

	makefileStr := string(content)

	// Required make targets for building
	requiredTargets := []string{
		"build:",          // Main build target
		"backend.build:",  // Backend-specific build
		"frontend.build:", // Frontend-specific build
	}

	for _, target := range requiredTargets {
		assert.True(t, strings.Contains(makefileStr, target),
			"Makefile should contain target %q", target)
	}
}

// TestComplexityThresholds verifies complexity thresholds are set appropriately
func TestComplexityThresholds(t *testing.T) {
	projectRoot := filepath.Join("..", "..")
	configPath := filepath.Join(projectRoot, ".golangci.yml")

	content, err := os.ReadFile(configPath)
	require.NoError(t, err, "should be able to read .golangci.yml")

	configStr := string(content)

	// Verify complexity settings exist
	assert.True(t, strings.Contains(configStr, "gocyclo"),
		"gocyclo complexity checker should be configured")
	assert.True(t, strings.Contains(configStr, "gocognit"),
		"gocognit cognitive complexity checker should be configured")
	assert.True(t, strings.Contains(configStr, "funlen"),
		"funlen function length checker should be configured")
}
