package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Release represents a versioned software release with artifacts
type Release struct {
	ID                  uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Version             string    `gorm:"uniqueIndex;not null;size:50"`
	GitSHA              string    `gorm:"column:git_sha;uniqueIndex;not null;size:40"`
	GitSHAShort         string    `gorm:"column:git_sha_short;not null;size:7"`
	Changelog           *string   `gorm:"type:text"`
	DockerImageBackend  string    `gorm:"not null;size:255"`
	DockerImageFrontend string    `gorm:"not null;size:255"`
	CreatedAt           time.Time
	CreatedBy           string       `gorm:"not null;size:100"`
	Deployments         []Deployment `gorm:"foreignKey:ReleaseID"`
}

// BeforeCreate hook to generate UUID if not set
func (r *Release) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for Release model
func (Release) TableName() string {
	return "releases"
}
