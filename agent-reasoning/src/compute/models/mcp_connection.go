package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	MCPAuthHeaderSecretVaultPathTemplate = "agent/%s/mcp_auth_header_secret/%s"
)

// MCPConnection holds the schema definition for the MCPConnection entity.
type MCPConnection struct {
	ID uuid.UUID `gorm:"type:uuid;not null" json:"id"`

	Name           string                      `gorm:"type:string;not null" json:"name"`
	Url            string                      `gorm:"type:string;not null" json:"url"`
	AgentID        uuid.UUID                   `gorm:"type:uuid;not null" json:"agentId"`
	HeaderAuthKeys MCPConnectionHeaderAuthKeys `gorm:"type:jsonb" json:"headerAuthKeys"`

	UpdatedAt time.Time `gorm:"type:timestamp;not null" json:"updatedAt"`
	CreatedAt time.Time `gorm:"type:timestamp;not null" json:"createdAt"`

	// Non persistent
	HeaderAuth map[string]string `gorm:"-" json:"headerAuth"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (entity *MCPConnection) BeforeCreate(tx *gorm.DB) (err error) {
	if entity.ID == uuid.Nil {
		entity.ID = uuid.New()
	}

	entity.CreatedAt = time.Now()
	return
}

// BeforeSave will update the UpdatedAt field.
func (entity *MCPConnection) BeforeSave(tx *gorm.DB) (err error) {
	entity.UpdatedAt = time.Now()
	return
}
