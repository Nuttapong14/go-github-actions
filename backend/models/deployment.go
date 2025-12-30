package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeploymentStatus represents the status of a deployment
type DeploymentStatus string

// Deployment status constants
const (
	StatusPending    DeploymentStatus = "pending"
	StatusRunning    DeploymentStatus = "running"
	StatusSuccess    DeploymentStatus = "success"
	StatusFailed     DeploymentStatus = "failed"
	StatusRolledBack DeploymentStatus = "rolled_back"
)

// TriggerType represents how a deployment was triggered
type TriggerType string

// Trigger type constants
const (
	TriggerManual    TriggerType = "manual"
	TriggerAutomatic TriggerType = "automatic"
	TriggerHotfix    TriggerType = "hotfix"
	TriggerRollback  TriggerType = "rollback"
)

// Deployment represents a single deployment event to an environment
type Deployment struct {
	ID            uuid.UUID        `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	EnvironmentID uuid.UUID        `gorm:"type:uuid;not null;index"`
	ReleaseID     uuid.UUID        `gorm:"type:uuid;not null;index"`
	Status        DeploymentStatus `gorm:"type:deployment_status;not null;default:'pending';index"`
	StartedAt     *time.Time
	CompletedAt   *time.Time
	TriggeredBy   string      `gorm:"not null;size:100"`
	TriggerType   TriggerType `gorm:"type:trigger_type;not null"`
	RollbackOfID  *uuid.UUID  `gorm:"type:uuid"`
	ErrorMessage  *string     `gorm:"type:text"`
	CreatedAt     time.Time

	Environment  Environment   `gorm:"foreignKey:EnvironmentID"`
	Release      Release       `gorm:"foreignKey:ReleaseID"`
	RollbackOf   *Deployment   `gorm:"foreignKey:RollbackOfID"`
	HealthChecks []HealthCheck `gorm:"foreignKey:DeploymentID"`
}

// BeforeCreate hook to generate UUID if not set
func (d *Deployment) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for Deployment model
func (Deployment) TableName() string {
	return "deployments"
}

// IsFinal returns true if the deployment is in a terminal state
func (d *Deployment) IsFinal() bool {
	return d.Status == StatusSuccess || d.Status == StatusFailed || d.Status == StatusRolledBack
}

// IsRollback returns true if this deployment is a rollback
func (d *Deployment) IsRollback() bool {
	return d.RollbackOfID != nil && d.TriggerType == TriggerRollback
}
