package db_clients

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models/enums"
	"gorm.io/gorm"
)

type AgentDatabaseClientInterface interface {
	Migrate() error
	GetDriver() enums.DatabaseConnectionDriver
	GetDB() *gorm.DB
	Close() error
	HealthCheck() error
	BeginTransaction() (*gorm.DB, error)
	CommitTransaction(tx *gorm.DB) error
	RollbackTransaction(tx *gorm.DB) error
}
