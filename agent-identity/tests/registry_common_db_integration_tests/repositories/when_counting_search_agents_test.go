package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db"
	db_repositories "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenCountingSearchAgentsTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *db_repositories.AgentRepository
}

func WhenCountingSearchAgentsBeforeEach() *WhenCountingSearchAgentsTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbClient, _ := registry_common_db.NewDatabaseClient(
		connectionString,
		zap.NewNop(),
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
	)
	if err = dbClient.Migrate(); err != nil {
		panic(err)
	}

	sut := db_repositories.NewAgentRepository(dbClient)

	return &WhenCountingSearchAgentsTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,
		sut:      sut,
	}
}

func WhenCountingSearchAgentsAfterEach(suite *WhenCountingSearchAgentsTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenCountingSearchAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given no agents in the database", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return zero count", func(t *testing.T) {
			t.Parallel()
			suite := WhenCountingSearchAgentsBeforeEach()
			defer WhenCountingSearchAgentsAfterEach(suite)

			params := repositories.AgentSearchParams{}
			count, err := suite.sut.SearchCount(params)

			assert.Nil(t, err)
			assert.Equal(t, uint(0), count)
		})
	})

	t.Run("Given agents in the database", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return correct count", func(t *testing.T) {
			t.Parallel()

			suite := WhenSearchingAgentsBeforeEach()
			defer WhenSearchingAgentsAfterEach(suite)

			for i := 0; i < 3; i++ {
				agent := model_builders.NewAgentBuilder().Build()
				suite.sut.Save(agent)
			}

			params := repositories.AgentSearchParams{}
			count, err := suite.sut.SearchCount(params)

			assert.Nil(t, err)
			assert.Equal(t, uint(3), count)
		})
	})
}
