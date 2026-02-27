package migrations_test

import (
	"fmt"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models/enums"
	database_connector_db_agent "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_agent"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_agent/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/common_tests/container_builders"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenMigratingMigrationsTestCase struct {
	name            string
	dialect         string
	clientBuilderFn func() container_builders.DatabaseContainerBuilderInterface
	driver          enums.DatabaseConnectionDriver
}

var (
	WhenMigratingMigrationsTestCases = []*WhenMigratingMigrationsTestCase{
		{
			name:    "PostgreSQL",
			dialect: "postgres",
			driver:  enums.DatabaseConnectionDriverPostgres,
			clientBuilderFn: func() container_builders.DatabaseContainerBuilderInterface {
				return container_builders.GetPostgresContainerBuilder()
			},
		},
		{
			name:    "MariaDB",
			dialect: "mysql",
			driver:  enums.DatabaseConnectionDriverMariaDB,
			clientBuilderFn: func() container_builders.DatabaseContainerBuilderInterface {
				return container_builders.GetMariaDBContainerBuilder()
			},
		},
		{
			name:    "MySQL",
			dialect: "mysql",
			driver:  enums.DatabaseConnectionDriverMySQL,
			clientBuilderFn: func() container_builders.DatabaseContainerBuilderInterface {
				return container_builders.GetMySQLContainerBuilder()
			},
		},
	}
)

type WhenMigratingMigrationsTestingSuite struct {
	dbName                   string
	migrationSettings        *settings.MigrationsSettings
	databaseContainerBuilder container_builders.DatabaseContainerBuilderInterface

	sut *database_connector_db_agent.AgentDatabaseClient
}

func WhenMigratingMigrationsBeforeEach(testCase *WhenMigratingMigrationsTestCase) *WhenMigratingMigrationsTestingSuite {
	databaseContainerBuilder := testCase.clientBuilderFn()
	dbName, connectionString := databaseContainerBuilder.Create()

	migrationFolderPath, err := common_tests.FindAgentMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	migrationsSettings := &settings.MigrationsSettings{
		MigrationsFolderPath: *migrationFolderPath,
	}

	dbClient, err := database_connector_db_agent.NewAgentDbClient(
		connectionString,
		zap.NewNop(),
		testCase.driver,
		migrationsSettings,
	)
	if err != nil {
		panic(err)
	}

	return &WhenMigratingMigrationsTestingSuite{
		dbName:                   dbName,
		migrationSettings:        migrationsSettings,
		databaseContainerBuilder: databaseContainerBuilder,

		sut: dbClient,
	}
}

func WhenMigratingMigrationsAfterEach(suite *WhenMigratingMigrationsTestingSuite) {
	suite.databaseContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenMigratingMigrations(t *testing.T) {
	t.Parallel()

	for _, testCase := range WhenMigratingMigrationsTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			t.Run("Given migrations are already applied", func(t *testing.T) {
				t.Parallel()

				t.Run("Should not apply migrations again and not return an error", func(t *testing.T) {
					t.Parallel()

					suite := WhenMigratingMigrationsBeforeEach(testCase)
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

					suite := WhenMigratingMigrationsBeforeEach(testCase)
					defer WhenMigratingMigrationsAfterEach(suite)

					err := suite.sut.Migrate()
					assert.NoError(t, err)

					// Rollback the last migration to simulate partial migration using Goose
					sqlDB, _ := suite.sut.DB.DB()
					if err := goose.SetDialect(testCase.dialect); err != nil {
						t.Fatalf("failed to set goose dialect: %v", err)
					}
					// Rollback one migration
					err = goose.Down(sqlDB, fmt.Sprintf("%s/%s", suite.migrationSettings.MigrationsFolderPath, testCase.driver.String()))
					assert.NoError(t, err)

					err = suite.sut.Migrate()
					assert.NoError(t, err)
				})
			})

			t.Run("Given no migrations are applied", func(t *testing.T) {
				t.Parallel()

				t.Run("Should apply all migrations and not return an error", func(t *testing.T) {
					t.Parallel()

					suite := WhenMigratingMigrationsBeforeEach(testCase)
					defer WhenMigratingMigrationsAfterEach(suite)

					err := suite.sut.Migrate()
					assert.NoError(t, err)
				})
			})
		})
	}
}
