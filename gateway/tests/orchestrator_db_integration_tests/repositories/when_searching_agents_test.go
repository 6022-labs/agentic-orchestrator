package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/db_clients"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_db"
	db_repositories "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_db/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_db/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/model_builders"
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

	dbClient, _ := orchestrator_db.NewDatabaseClient(
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

	t.Run("Given IsInternal parameters", func(t *testing.T) {
		t.Parallel()

		t.Run("And searching with IsInternal true", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return only internal agents", func(t *testing.T) {
				t.Parallel()
				suite := WhenSearchingAgentsBeforeEach()
				defer WhenSearchingAgentsAfterEach(suite)

				internalAgent := model_builders.NewAgentBuilder().WithIsInternal(true).Build()
				nonInternalAgent := model_builders.NewAgentBuilder().WithIsInternal(false).Build()
				suite.sut.Save(internalAgent)
				suite.sut.Save(nonInternalAgent)

				params := repositories.AgentSearchParamsWithPagination{
					AgentSearchParams: repositories.AgentSearchParams{IsInternal: ptrBool(true)},
				}
				agents, err := suite.sut.Search(params)

				assert.Nil(t, err)
				assert.Len(t, agents, 1)
				assert.True(t, agents[0].IsInternal)
			})
		})

		t.Run("And searching with IsInternal false", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return only external agents", func(t *testing.T) {
				t.Parallel()
				suite := WhenSearchingAgentsBeforeEach()
				defer WhenSearchingAgentsAfterEach(suite)

				internalAgent := model_builders.NewAgentBuilder().WithIsInternal(true).Build()
				nonInternalAgent := model_builders.NewAgentBuilder().WithIsInternal(false).Build()
				suite.sut.Save(internalAgent)
				suite.sut.Save(nonInternalAgent)

				params := repositories.AgentSearchParamsWithPagination{
					AgentSearchParams: repositories.AgentSearchParams{IsInternal: ptrBool(false)},
				}
				agents, err := suite.sut.Search(params)

				assert.Nil(t, err)
				assert.Len(t, agents, 1)
				assert.False(t, agents[0].IsInternal)
			})
		})

		t.Run("And not searching with IsInternal", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return all agents", func(t *testing.T) {
				t.Parallel()
				suite := WhenSearchingAgentsBeforeEach()
				defer WhenSearchingAgentsAfterEach(suite)

				internalAgent := model_builders.NewAgentBuilder().WithIsInternal(true).Build()
				nonInternalAgent := model_builders.NewAgentBuilder().WithIsInternal(false).Build()
				suite.sut.Save(internalAgent)
				suite.sut.Save(nonInternalAgent)

				params := repositories.AgentSearchParamsWithPagination{}
				agents, err := suite.sut.Search(params)

				assert.Nil(t, err)
				assert.Len(t, agents, 2)
			})
		})
	})

	t.Run("Given pagination parameters", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return correct number of agents per page", func(t *testing.T) {
			t.Parallel()
			suite := WhenSearchingAgentsBeforeEach()
			defer WhenSearchingAgentsAfterEach(suite)

			for i := 0; i < 5; i++ {
				agent := model_builders.NewAgentBuilder().WithIsInternal(true).Build()
				suite.sut.Save(agent)
			}

			params := repositories.AgentSearchParamsWithPagination{
				AgentSearchParams: repositories.AgentSearchParams{IsInternal: ptrBool(true)},
				ElementsPerPage:   ptrUint(2),
				Page:              ptrUint(1),
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
				agent := model_builders.NewAgentBuilder().WithIsInternal(true).Build()
				suite.sut.Save(agent)
			}

			params := repositories.AgentSearchParamsWithPagination{
				AgentSearchParams: repositories.AgentSearchParams{IsInternal: ptrBool(true)},
				ElementsPerPage:   ptrUint(2),
				Page:              ptrUint(2),
			}
			agents, err := suite.sut.Search(params)

			assert.Nil(t, err)
			assert.Len(t, agents, 2)
		})
	})
}

func ptrBool(b bool) *bool { return &b }
func ptrUint(i uint) *uint { return &i }
