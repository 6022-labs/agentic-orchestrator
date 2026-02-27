package models

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	DatabaseCredentialsVaultPathTemplate = "agent/%s/database_credentials"
)

// DatabaseConnection holds the schema definition for the DatabaseConnection entity.
type DatabaseConnection struct {
	ID uuid.UUID `gorm:"type:uuid;not null primaryKey" json:"id"`

	Host     string                         `gorm:"type:string;not null" json:"host"`
	Port     uint16                         `gorm:"type:integer;not null" json:"port"`
	AgentID  uuid.UUID                      `gorm:"type:uuid;not null" json:"agentId"`
	Driver   enums.DatabaseConnectionDriver `gorm:"type:string;not null" json:"driver"`
	Username string                         `gorm:"type:string;not null" json:"username"`

	SetupStatus enums.DatabaseConnectionSetupStatus `gorm:"type:string;not null" json:"setupStatus"`

	MigratedMigrations []*DatabaseConnectionMigration `gorm:"many2many:database_connection_migration_connections;"`

	UpdatedAt time.Time  `gorm:"type:timestamp;not null" json:"updatedAt"`
	CreatedAt time.Time  `gorm:"type:timestamp;not null" json:"createdAt"`
	DeletedAt *time.Time `gorm:"type:timestamp" json:"deletedAt"`
}

func (entity *DatabaseConnection) BeforeCreate(tx *gorm.DB) error {
	if entity.ID == uuid.Nil {
		entity.ID = uuid.New()
	}

	if entity.CreatedAt.IsZero() {
		entity.CreatedAt = time.Now()
	}

	return nil
}

func (entity *DatabaseConnection) BeforeSave(tx *gorm.DB) error {
	entity.UpdatedAt = time.Now()

	return nil
}
