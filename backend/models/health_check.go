package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HealthCheckType represents the type of health check
type HealthCheckType string

// Health check type constants
const (
	CheckTypeLiveness  HealthCheckType = "liveness"
	CheckTypeReadiness HealthCheckType = "readiness"
)

// HealthStatus represents the status of a health check
type HealthStatus string

// Health status constants
const (
	HealthPassing HealthStatus = "passing"
	HealthFailing HealthStatus = "failing"
)

// HealthCheck represents health check results after deployment
type HealthCheck struct {
	ID             uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	DeploymentID   uuid.UUID       `gorm:"type:uuid;not null;index"`
	CheckType      HealthCheckType `gorm:"type:health_check_type;not null"`
	Status         HealthStatus    `gorm:"type:health_status;not null"`
	ResponseTimeMs *int
	StatusCode     *int
	ErrorMessage   *string    `gorm:"type:text"`
	CheckedAt      time.Time  `gorm:"default:now()"`
	Deployment     Deployment `gorm:"foreignKey:DeploymentID"`
}

// BeforeCreate hook to generate UUID if not set
func (h *HealthCheck) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for HealthCheck model
func (HealthCheck) TableName() string {
	return "health_checks"
}

// IsPassing returns true if the health check passed
func (h *HealthCheck) IsPassing() bool {
	return h.Status == HealthPassing
}
