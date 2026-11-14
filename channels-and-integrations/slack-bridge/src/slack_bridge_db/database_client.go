package slack_bridge_db

import (
	"fmt"
	"sync"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/metrics"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

// goose does not support concurrent migrations, so a global mutex ensures only one runs at a time.
var gooseMigrationMu sync.Mutex

type DatabaseClient struct {
	*gorm.DB
	migrationsSettings *settings.MigrationsSettings
}

func NewDatabaseClient(
	databaseSettings *settings.DatabaseSettings,
	logger *zap.Logger,
	migrationsSettings *settings.MigrationsSettings,
	dbQueryMetrics metrics_interfaces.DbQueryMetricsInterface,
) (*DatabaseClient, error) {
	dialector := postgres.Open(databaseSettings.Url)
	if databaseSettings.TlsConfig != nil {
		pgxConfig, err := pgx.ParseConfig(databaseSettings.Url)
		if err != nil {
			return nil, fmt.Errorf("failed to parse database url: %w", err)
		}
		tlsConfig := databaseSettings.TlsConfig.Clone()
		// Keep what pgx derived from the url: ALPN for sslnegotiation=direct, sslcert/sslkey client certificates.
		if pgxConfig.TLSConfig != nil {
			tlsConfig.NextProtos = pgxConfig.TLSConfig.NextProtos
			tlsConfig.Certificates = pgxConfig.TLSConfig.Certificates
			tlsConfig.GetClientCertificate = pgxConfig.TLSConfig.GetClientCertificate
		} else if pgxConfig.SSLNegotiation == "direct" {
			tlsConfig.NextProtos = []string{"postgresql"}
		}
		pgxConfig.TLSConfig = tlsConfig
		// sslmode prefer/allow fallbacks would retry in plaintext.
		pgxConfig.Fallbacks = nil
		dialector = postgres.New(postgres.Config{Conn: stdlib.OpenDB(*pgxConfig)})
	}

	gormLogger := zapgorm2.New(logger)
	gormLogger.IgnoreRecordNotFoundError = true
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Use(metrics.NewGormMetricsPlugin(dbQueryMetrics)); err != nil {
		return nil, fmt.Errorf("failed to register db metrics plugin: %w", err)
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
