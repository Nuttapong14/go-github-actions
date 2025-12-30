package tests

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Nuttapong14/go-github-actions/backend/models"
)

// TestDeploymentStatus verifies deployment status constants
func TestDeploymentStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   models.DeploymentStatus
		expected string
	}{
		{"pending status", models.StatusPending, "pending"},
		{"running status", models.StatusRunning, "running"},
		{"success status", models.StatusSuccess, "success"},
		{"failed status", models.StatusFailed, "failed"},
		{"rolled_back status", models.StatusRolledBack, "rolled_back"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.status)
			}
		})
	}
}

// TestTriggerType verifies trigger type constants
func TestTriggerType(t *testing.T) {
	tests := []struct {
		name     string
		trigger  models.TriggerType
		expected string
	}{
		{"manual trigger", models.TriggerManual, "manual"},
		{"automatic trigger", models.TriggerAutomatic, "automatic"},
		{"hotfix trigger", models.TriggerHotfix, "hotfix"},
		{"rollback trigger", models.TriggerRollback, "rollback"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.trigger) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.trigger)
			}
		})
	}
}

// TestDeploymentIsFinal verifies IsFinal method
func TestDeploymentIsFinal(t *testing.T) {
	tests := []struct {
		name     string
		status   models.DeploymentStatus
		expected bool
	}{
		{"pending is not final", models.StatusPending, false},
		{"running is not final", models.StatusRunning, false},
		{"success is final", models.StatusSuccess, true},
		{"failed is final", models.StatusFailed, true},
		{"rolled_back is final", models.StatusRolledBack, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &models.Deployment{Status: tt.status}
			if d.IsFinal() != tt.expected {
				t.Errorf("IsFinal() = %v, expected %v", d.IsFinal(), tt.expected)
			}
		})
	}
}

// TestDeploymentIsRollback verifies IsRollback method
func TestDeploymentIsRollback(t *testing.T) {
	rollbackID := uuid.New()

	tests := []struct {
		name         string
		rollbackOfID *uuid.UUID
		triggerType  models.TriggerType
		expected     bool
	}{
		{"not a rollback - nil ID", nil, models.TriggerManual, false},
		{"not a rollback - wrong trigger", &rollbackID, models.TriggerManual, false},
		{"is a rollback", &rollbackID, models.TriggerRollback, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &models.Deployment{
				RollbackOfID: tt.rollbackOfID,
				TriggerType:  tt.triggerType,
			}
			if d.IsRollback() != tt.expected {
				t.Errorf("IsRollback() = %v, expected %v", d.IsRollback(), tt.expected)
			}
		})
	}
}

// TestDeploymentTableName verifies table name
func TestDeploymentTableName(t *testing.T) {
	d := models.Deployment{}
	if d.TableName() != "deployments" {
		t.Errorf("expected table name 'deployments', got %s", d.TableName())
	}
}

// TestEnvironmentTableName verifies environment table name
func TestEnvironmentTableName(t *testing.T) {
	e := models.Environment{}
	if e.TableName() != "environments" {
		t.Errorf("expected table name 'environments', got %s", e.TableName())
	}
}

// TestReleaseTableName verifies release table name
func TestReleaseTableName(t *testing.T) {
	r := models.Release{}
	if r.TableName() != "releases" {
		t.Errorf("expected table name 'releases', got %s", r.TableName())
	}
}

// TestDeploymentStructFields verifies Deployment struct fields
func TestDeploymentStructFields(t *testing.T) {
	now := time.Now()
	envID := uuid.New()
	releaseID := uuid.New()
	rollbackID := uuid.New()
	errMsg := "deployment failed"

	d := models.Deployment{
		ID:            uuid.New(),
		EnvironmentID: envID,
		ReleaseID:     releaseID,
		Status:        models.StatusRunning,
		StartedAt:     &now,
		CompletedAt:   nil,
		TriggeredBy:   "ci-system",
		TriggerType:   models.TriggerAutomatic,
		RollbackOfID:  &rollbackID,
		ErrorMessage:  &errMsg,
		CreatedAt:     now,
	}

	if d.EnvironmentID != envID {
		t.Error("EnvironmentID not set correctly")
	}
	if d.ReleaseID != releaseID {
		t.Error("ReleaseID not set correctly")
	}
	if d.Status != models.StatusRunning {
		t.Error("Status not set correctly")
	}
	if d.TriggeredBy != "ci-system" {
		t.Error("TriggeredBy not set correctly")
	}
	if d.TriggerType != models.TriggerAutomatic {
		t.Error("TriggerType not set correctly")
	}
	if *d.RollbackOfID != rollbackID {
		t.Error("RollbackOfID not set correctly")
	}
	if *d.ErrorMessage != errMsg {
		t.Error("ErrorMessage not set correctly")
	}
}

// TestEnvironmentStructFields verifies Environment struct fields
func TestEnvironmentStructFields(t *testing.T) {
	e := models.Environment{
		ID:               uuid.New(),
		Name:             "alpha",
		URL:              "https://alpha.example.com",
		ApprovalRequired: false,
		IsProduction:     false,
		SoakPeriodHours:  0,
	}

	if e.Name != "alpha" {
		t.Error("Name not set correctly")
	}
	if e.URL != "https://alpha.example.com" {
		t.Error("URL not set correctly")
	}
	if e.ApprovalRequired != false {
		t.Error("ApprovalRequired should be false for alpha")
	}
	if e.IsProduction != false {
		t.Error("IsProduction should be false for alpha")
	}
}

// TestReleaseStructFields verifies Release struct fields
func TestReleaseStructFields(t *testing.T) {
	changelog := "Initial release"
	r := models.Release{
		ID:                  uuid.New(),
		Version:             "1.0.0",
		GitSHA:              "abc123def456789012345678901234567890abcd",
		GitSHAShort:         "abc123d",
		Changelog:           &changelog,
		DockerImageBackend:  "ghcr.io/example/backend:1.0.0",
		DockerImageFrontend: "ghcr.io/example/frontend:1.0.0",
		CreatedBy:           "ci-system",
	}

	if r.Version != "1.0.0" {
		t.Error("Version not set correctly")
	}
	if len(r.GitSHA) != 40 {
		t.Error("GitSHA should be 40 characters")
	}
	if len(r.GitSHAShort) != 7 {
		t.Error("GitSHAShort should be 7 characters")
	}
	if *r.Changelog != changelog {
		t.Error("Changelog not set correctly")
	}
}

// TestHealthCheckTableName verifies health_check table name
func TestHealthCheckTableName(t *testing.T) {
	h := models.HealthCheck{}
	if h.TableName() != "health_checks" {
		t.Errorf("expected table name 'health_checks', got %s", h.TableName())
	}
}

// TestHealthCheckTypes verifies health check type constants
func TestHealthCheckTypes(t *testing.T) {
	tests := []struct {
		name      string
		checkType models.HealthCheckType
		expected  string
	}{
		{"liveness check", models.CheckTypeLiveness, "liveness"},
		{"readiness check", models.CheckTypeReadiness, "readiness"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.checkType) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.checkType)
			}
		})
	}
}

// TestHealthStatus verifies health status constants
func TestHealthStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   models.HealthStatus
		expected string
	}{
		{"passing status", models.HealthPassing, "passing"},
		{"failing status", models.HealthFailing, "failing"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.status)
			}
		})
	}
}
