package model_builders

import (
	"strings"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models/enums"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type DatabaseConnectionBuilder struct {
	id uuid.UUID

	host     string
	port     uint16
	driver   enums.DatabaseConnectionDriver
	agentId  uuid.UUID
	username string

	setupStatus enums.DatabaseConnectionSetupStatus

	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

func NewDatabaseConnectionBuilder() *DatabaseConnectionBuilder {
	id := uuid.New()
	host := faker.URL()
	port, _ := faker.RandomInt(0, 65535, 1)
	driver := enums.DatabaseConnectionDriverPostgres

	agentId := uuid.New()
	username := strings.ReplaceAll(agentId.String(), "-", "")

	setupStatus := enums.DatabaseConnectionSetupStatusEstablished

	createdAtStr := faker.Date()
	updatedAtStr := faker.Date()

	createdAt, _ := time.Parse(time.RFC3339, createdAtStr)
	updatedAt, _ := time.Parse(time.RFC3339, updatedAtStr)

	return &DatabaseConnectionBuilder{
		id:          id,
		host:        host,
		driver:      driver,
		agentId:     agentId,
		setupStatus: setupStatus,
		port:        uint16(port[0]),
		username:    username,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		deletedAt:   nil,
	}
}

func (b *DatabaseConnectionBuilder) WithID(id uuid.UUID) *DatabaseConnectionBuilder {
	b.id = id
	return b
}

func (b *DatabaseConnectionBuilder) WithHost(host string) *DatabaseConnectionBuilder {
	b.host = host
	return b
}

func (b *DatabaseConnectionBuilder) WithPort(port uint16) *DatabaseConnectionBuilder {
	b.port = port
	return b
}

func (b *DatabaseConnectionBuilder) WithDriver(driver enums.DatabaseConnectionDriver) *DatabaseConnectionBuilder {
	b.driver = driver
	return b
}

func (b *DatabaseConnectionBuilder) WithUsername(username string) *DatabaseConnectionBuilder {
	b.username = username
	return b
}

func (b *DatabaseConnectionBuilder) WithAgentId(agentId uuid.UUID) *DatabaseConnectionBuilder {
	b.agentId = agentId
	return b
}

func (b *DatabaseConnectionBuilder) WithSetupStatus(setupStatus enums.DatabaseConnectionSetupStatus) *DatabaseConnectionBuilder {
	b.setupStatus = setupStatus
	return b
}

func (b *DatabaseConnectionBuilder) WithCreatedAt(createdAt time.Time) *DatabaseConnectionBuilder {
	b.createdAt = createdAt
	return b
}

func (b *DatabaseConnectionBuilder) WithUpdatedAt(updatedAt time.Time) *DatabaseConnectionBuilder {
	b.updatedAt = updatedAt
	return b
}

func (b *DatabaseConnectionBuilder) WithDeletedAt(deletedAt *time.Time) *DatabaseConnectionBuilder {
	b.deletedAt = deletedAt
	return b
}

func (b *DatabaseConnectionBuilder) Build() *models.DatabaseConnection {
	return &models.DatabaseConnection{
		ID:          b.id,
		Host:        b.host,
		Port:        b.port,
		Driver:      b.driver,
		AgentID:     b.agentId,
		Username:    b.username,
		CreatedAt:   b.createdAt,
		UpdatedAt:   b.updatedAt,
		DeletedAt:   b.deletedAt,
		SetupStatus: b.setupStatus,
	}
}
