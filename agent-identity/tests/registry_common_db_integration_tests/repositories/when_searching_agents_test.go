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
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenSearchingAgentsTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *db_repositories.AgentRepository
}

func WhenSearchingAgentsBeforeEach() *WhenSearchingAgentsTestingSuite {
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

	return &WhenSearchingAgentsTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,
		sut:      sut,
	}
}

func WhenSearchingAgentsAfterEach(suite *WhenSearchingAgentsTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenSearchingAgents(t *testing.T) {
	t.Parallel()

	t.Run("Given agents in the database", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return all agents", func(t *testing.T) {
			t.Parallel()
			suite := WhenSearchingAgentsBeforeEach()
			defer WhenSearchingAgentsAfterEach(suite)

			firstAgent := model_builders.NewAgentBuilder().Build()
			secondAgent := model_builders.NewAgentBuilder().Build()
			suite.sut.Save(firstAgent)
			suite.sut.Save(secondAgent)

			params := repositories.AgentSearchParamsWithPagination{}
			agents, err := suite.sut.Search(params)

			assert.Nil(t, err)
			assert.Len(t, agents, 2)
		})
	})

	t.Run("Given pagination parameters", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return correct number of agents per page", func(t *testing.T) {
			t.Parallel()
			suite := WhenSearchingAgentsBeforeEach()
			defer WhenSearchingAgentsAfterEach(suite)

			for i := 0; i < 5; i++ {
				agent := model_builders.NewAgentBuilder().Build()
				suite.sut.Save(agent)
			}

			params := repositories.AgentSearchParamsWithPagination{
				ElementsPerPage: lo.ToPtr(uint(2)),
				Page:            lo.ToPtr(uint(1)),
			}
			agents, err := suite.sut.Search(params)

			assert.Nil(t, err)
			assert.Len(t, agents, 2)
		})

		t.Run("Should return correct agents for second page", func(t *testing.T) {
			t.Parallel()
			suite := WhenSearchingAgentsBeforeEach()
			defer WhenSearchingAgentsAfterEach(suite)

			for i := 0; i < 5; i++ {
				agent := model_builders.NewAgentBuilder().Build()
				suite.sut.Save(agent)
			}

			params := repositories.AgentSearchParamsWithPagination{
				ElementsPerPage: lo.ToPtr(uint(2)),
				Page:            lo.ToPtr(uint(2)),
			}
			agents, err := suite.sut.Search(params)

			assert.Nil(t, err)
			assert.Len(t, agents, 2)
		})
	})

	t.Run("Given minted and proposal agents in the database", func(t *testing.T) {
		t.Parallel()

		saveMixedAgents := func(suite *WhenSearchingAgentsTestingSuite) {
			for i := 0; i < 2; i++ {
				suite.sut.Save(model_builders.NewAgentBuilder().Build())
			}
			suite.sut.Save(model_builders.NewAgentBuilder().WithCollectionTokenId(nil).Build())
		}

		t.Run("Given the filter selects minted agents", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return only minted agents", func(t *testing.T) {
				t.Parallel()
				suite := WhenSearchingAgentsBeforeEach()
				defer WhenSearchingAgentsAfterEach(suite)
				saveMixedAgents(suite)

				params := repositories.AgentSearchParamsWithPagination{
					AgentSearchParams: repositories.AgentSearchParams{IsMinted: lo.ToPtr(true)},
				}
				agents, err := suite.sut.Search(params)

				assert.Nil(t, err)
				assert.Len(t, agents, 2)
			})
		})

		t.Run("Given the filter selects proposal agents", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return only proposal agents", func(t *testing.T) {
				t.Parallel()
				suite := WhenSearchingAgentsBeforeEach()
				defer WhenSearchingAgentsAfterEach(suite)
				saveMixedAgents(suite)

				params := repositories.AgentSearchParamsWithPagination{
					AgentSearchParams: repositories.AgentSearchParams{IsMinted: lo.ToPtr(false)},
				}
				agents, err := suite.sut.Search(params)

				assert.Nil(t, err)
				assert.Len(t, agents, 1)
			})
		})
	})
}
