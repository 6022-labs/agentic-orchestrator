package migrations_test

import (
	"testing"

	orchestrator_agent_database_connector_db_orchestrator "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_orchestrator"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_orchestrator/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/common_tests/container_builders"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenMigratingMigrationsTestingSuite struct {
	dbName                   string
	migrationSettings        *settings.MigrationsSettings
	postgresContainerBuilder *container_builders.PostgresContainerBuilder

	sut *orchestrator_agent_database_connector_db_orchestrator.OrchestratorDatabaseClient
}

func WhenMigratingMigrationsBeforeEach() *WhenMigratingMigrationsTestingSuite {
	postgresContainerBuilder := container_builders.GetPostgresContainerBuilder()
	dbName, connectionString := postgresContainerBuilder.Create()

	migrationFolderPath, err := common_tests.FindOrchestratorMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	migrationsSettings := &settings.MigrationsSettings{
		MigrationsFolderPath: *migrationFolderPath,
	}

	dbClient, _ := orchestrator_agent_database_connector_db_orchestrator.NewOrchestratorDatabaseClient(
		connectionString,
		zap.NewNop(),
		migrationsSettings,
	)

	return &WhenMigratingMigrationsTestingSuite{
		dbName:                   dbName,
		migrationSettings:        migrationsSettings,
		postgresContainerBuilder: postgresContainerBuilder,

		sut: dbClient,
	}
}

func WhenMigratingMigrationsAfterEach(suite *WhenMigratingMigrationsTestingSuite) {
	suite.postgresContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenMigratingMigrations(t *testing.T) {
	t.Parallel()

	t.Run("Given migrations are already applied", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not apply migrations again and not return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenMigratingMigrationsBeforeEach()
			defer WhenMigratingMigrationsAfterEach(suite)

			err := suite.sut.Migrate()
			assert.NoError(t, err)

			err = suite.sut.Migrate()
			assert.NoError(t, err)
		})
	})

	t.Run("Given migrations are partially applied", func(t *testing.T) {
		t.Parallel()

		t.Run("Should apply the remaining migrations and not return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenMigratingMigrationsBeforeEach()
			defer WhenMigratingMigrationsAfterEach(suite)

			err := suite.sut.Migrate()
			assert.NoError(t, err)

			// Rollback the last migration to simulate partial migration using Goose
			sqlDB, _ := suite.sut.DB.DB()
			if err := goose.SetDialect("postgres"); err != nil {
				t.Fatalf("failed to set goose dialect: %v", err)
			}
			// Rollback one migration
			err = goose.Down(sqlDB, suite.migrationSettings.MigrationsFolderPath)
			assert.NoError(t, err)

			err = suite.sut.Migrate()
			assert.NoError(t, err)
		})
	})

	t.Run("Given no migrations are applied", func(t *testing.T) {
		t.Parallel()

		t.Run("Should apply all migrations and not return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenMigratingMigrationsBeforeEach()
			defer WhenMigratingMigrationsAfterEach(suite)

			err := suite.sut.Migrate()
			assert.NoError(t, err)
		})
	})
}
