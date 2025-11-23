package model_builders

import (
	"time"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type DatabaseConnectionMigrationBuilder struct {
	id uuid.UUID

	name string

	createdAt time.Time
	updatedAt time.Time
}

func NewDatabaseConnectionMigrationBuilder() *DatabaseConnectionMigrationBuilder {
	id := uuid.New()
	name := faker.Word()

	createdAtStr := faker.Date()
	updatedAtStr := faker.Date()

	createdAt, _ := time.Parse(time.RFC3339, createdAtStr)
	updatedAt, _ := time.Parse(time.RFC3339, updatedAtStr)

	return &DatabaseConnectionMigrationBuilder{
		id:        id,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (b *DatabaseConnectionMigrationBuilder) WithID(id uuid.UUID) *DatabaseConnectionMigrationBuilder {
	b.id = id
	return b
}

func (b *DatabaseConnectionMigrationBuilder) WithName(name string) *DatabaseConnectionMigrationBuilder {
	b.name = name
	return b
}

func (b *DatabaseConnectionMigrationBuilder) Build() *models.DatabaseConnectionMigration {
	return &models.DatabaseConnectionMigration{
		ID:        b.id,
		Name:      b.name,
		CreatedAt: b.createdAt,
		UpdatedAt: b.updatedAt,
	}
}
