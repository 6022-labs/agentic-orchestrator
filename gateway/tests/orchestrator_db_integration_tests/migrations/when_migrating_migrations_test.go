package migrations_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_db"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_db/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/container_builders"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenMigratingMigrationsTestingSuite struct {
	dbName            string
	migrationSettings *settings.MigrationsSettings

	sut *orchestrator_db.DatabaseClient
}

func WhenMigratingMigrationsBeforeEach() *WhenMigratingMigrationsTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	migrationsSettings := &settings.MigrationsSettings{
		MigrationsFolderPath: *migrationFolderPath,
	}

	dbClient, _ := orchestrator_db.NewDatabaseClient(
		connectionString,
		zap.NewNop(),
		migrationsSettings,
	)

	return &WhenMigratingMigrationsTestingSuite{
		dbName:            dbName,
		migrationSettings: migrationsSettings,

		sut: dbClient,
	}
}

func WhenMigratingMigrationsAfterEach(suite *WhenMigratingMigrationsTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
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
