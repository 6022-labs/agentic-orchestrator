package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DatabaseConnectionMigration holds the schema definition for the DatabaseConnectionMigration entity.
type DatabaseConnectionMigration struct {
	ID uuid.UUID `gorm:"type:uuid;not null primaryKey" json:"id"`

	Name string `gorm:"type:string;not null" json:"name"`

	DatabaseConnections []*DatabaseConnection `gorm:"many2many:database_connection_migration_connections;"`

	UpdatedAt time.Time `gorm:"type:timestamp;not null" json:"updatedAt"`
	CreatedAt time.Time `gorm:"type:timestamp;not null" json:"createdAt"`
}

func (entity *DatabaseConnectionMigration) BeforeCreate(tx *gorm.DB) error {
	if entity.ID == uuid.Nil {
		entity.ID = uuid.New()
	}

	if entity.CreatedAt.IsZero() {
		entity.CreatedAt = time.Now()
	}

	return nil
}

func (entity *DatabaseConnectionMigration) BeforeSave(tx *gorm.DB) error {
	entity.UpdatedAt = time.Now()

	return nil
}
