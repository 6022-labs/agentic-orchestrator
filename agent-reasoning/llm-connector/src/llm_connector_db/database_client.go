package llm_connector_db

import (
	"fmt"
	"sync"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector_db/settings"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

// goose does not support concurrent migrations, so we use a global mutex to ensure that only one migration runs at a time.
var gooseMigrationMu sync.Mutex

type DatabaseClient struct {
	*gorm.DB
	migrationsSettings *settings.MigrationsSettings
}

func NewOrchestratorDbClient(
	connStr string,
	logger *zap.Logger,
	migrationsSettings *settings.MigrationsSettings,
) (*DatabaseClient, error) {
	dialector := postgres.Open(connStr)

	gormLogger := zapgorm2.New(logger)
	gormLogger.IgnoreRecordNotFoundError = true
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	return &DatabaseClient{
		DB:                 db,
		migrationsSettings: migrationsSettings,
	}, nil
}

func (c *DatabaseClient) Migrate() error {
	gooseMigrationMu.Lock()
	defer gooseMigrationMu.Unlock()

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

func (c *DatabaseClient) GetDB() *gorm.DB {
	return c.DB
}

func (c *DatabaseClient) HealthCheck() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (c *DatabaseClient) BeginTransaction() *gorm.DB {
	return c.DB.Begin()
}

func (c *DatabaseClient) CommitTransaction(tx *gorm.DB) {
	tx.Commit()
}

func (c *DatabaseClient) RollbackTransaction(tx *gorm.DB) {
	tx.Rollback()
}
