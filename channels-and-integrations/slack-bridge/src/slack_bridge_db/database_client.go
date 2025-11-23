package slack_bridge_db

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

type DatabaseClient struct {
	*gorm.DB
	migrationsSettings *settings.MigrationsSettings
}

func NewDatabaseClient(
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
