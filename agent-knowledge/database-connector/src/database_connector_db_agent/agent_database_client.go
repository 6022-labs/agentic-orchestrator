package database_connector_db_agent

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_agent/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_agent/utils"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

type AgentDatabaseClient struct {
	*gorm.DB
	driver             enums.DatabaseConnectionDriver
	migrationsSettings *settings.MigrationsSettings
}

func NewAgentDbClient(
	connStr string,
	logger *zap.Logger,
	driver enums.DatabaseConnectionDriver,
	migrationsSettings *settings.MigrationsSettings,
) (*AgentDatabaseClient, error) {
	var dialector gorm.Dialector

	switch driver {
	case enums.DatabaseConnectionDriverMySQL, enums.DatabaseConnectionDriverMariaDB:
		dialector = mysql.Open(connStr)
	case enums.DatabaseConnectionDriverPostgres:
		dialector = postgres.Open(connStr)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	gormLogger := zapgorm2.New(logger)
	gormLogger.IgnoreRecordNotFoundError = true
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	return &AgentDatabaseClient{
		DB:                 db,
		driver:             driver,
		migrationsSettings: migrationsSettings,
	}, nil
}

func (c *AgentDatabaseClient) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sqlDB: %v", err)
	}
	return sqlDB.Close()
}

func (c *AgentDatabaseClient) Migrate() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := goose.SetDialect(c.driver.ToGooseDialect()); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(sqlDB, utils.GetMigrationFolderPath(c.migrationsSettings.MigrationsFolderPath, c.driver.String())); err != nil {
		return fmt.Errorf("goose migration failed: %w", err)
	}

	return nil
}

func (c *AgentDatabaseClient) GetDriver() enums.DatabaseConnectionDriver {
	return c.driver
}

func (c *AgentDatabaseClient) GetDB() *gorm.DB {
	return c.DB
}

func (c *AgentDatabaseClient) HealthCheck() error {
	sqlDB, err := c.GetDB().DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (c *AgentDatabaseClient) BeginTransaction() (*gorm.DB, error) {
	transaction := c.GetDB().Begin()
	if transaction.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", transaction.Error)
	}

	return transaction, nil
}

func (c *AgentDatabaseClient) CommitTransaction(tx *gorm.DB) error {
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (c *AgentDatabaseClient) RollbackTransaction(tx *gorm.DB) error {
	if err := tx.Rollback().Error; err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	return nil
}
