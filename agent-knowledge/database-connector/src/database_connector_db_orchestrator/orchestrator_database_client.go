package orchestrator_agent_database_connector_db_orchestrator

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_orchestrator/settings"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

type OrchestratorDatabaseClient struct {
	*gorm.DB
	migrationsSettings *settings.MigrationsSettings
}

func NewOrchestratorDatabaseClient(
	connStr string,
	logger *zap.Logger,
	migrationsSettings *settings.MigrationsSettings,
) (*OrchestratorDatabaseClient, error) {
	dialector := postgres.Open(connStr)

	gormLogger := zapgorm2.New(logger)
	gormLogger.IgnoreRecordNotFoundError = true
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	return &OrchestratorDatabaseClient{
		DB:                 db,
		migrationsSettings: migrationsSettings,
	}, nil
}

func (c *OrchestratorDatabaseClient) Migrate() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Set Goose dialect
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// Apply all migrations in the folder
	if err := goose.Up(sqlDB, c.migrationsSettings.MigrationsFolderPath); err != nil {
		return fmt.Errorf("goose migration failed: %w", err)
	}

	return nil
}

func (c *OrchestratorDatabaseClient) GetDB() *gorm.DB {
	return c.DB
}

func (c *OrchestratorDatabaseClient) HealthCheck() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (c *OrchestratorDatabaseClient) BeginTransaction() *gorm.DB {
	return c.DB.Begin()
}

func (c *OrchestratorDatabaseClient) CommitTransaction(tx *gorm.DB) {
	tx.Commit()
}

func (c *OrchestratorDatabaseClient) RollbackTransaction(tx *gorm.DB) {
	tx.Rollback()
}
