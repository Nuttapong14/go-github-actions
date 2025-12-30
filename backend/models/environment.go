package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Environment represents a deployment target in the promotion pipeline
type Environment struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name             string    `gorm:"uniqueIndex;not null;size:50"`
	URL              string    `gorm:"not null;size:255"`
	ApprovalRequired bool      `gorm:"default:false"`
	IsProduction     bool      `gorm:"default:false"`
	SoakPeriodHours  int       `gorm:"default:0"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Deployments      []Deployment `gorm:"foreignKey:EnvironmentID"`
}

// BeforeCreate hook to generate UUID if not set
func (e *Environment) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for Environment model
func (Environment) TableName() string {
	return "environments"
}
