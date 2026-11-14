package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db"
	db_repositories "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	db_metrics_mocks "github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_db_mocks/metrics_mocks/interfaces_mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenCountingSearchAgentsTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *db_repositories.AgentRepository
}

func WhenCountingSearchAgentsBeforeEach(t *testing.T) *WhenCountingSearchAgentsTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbQueryMetrics := db_metrics_mocks.NewMockDbQueryMetricsInterface(gomock.NewController(t))
	dbQueryMetrics.EXPECT().
		RecordQuery(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes()

	dbClient, _ := registry_common_db.NewDatabaseClient(
		&settings.DatabaseSettings{Url: connectionString},
		zap.NewNop(),
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
		dbQueryMetrics,
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
			suite := WhenCountingSearchAgentsBeforeEach(t)
			defer WhenCountingSearchAgentsAfterEach(suite)

			params := repositories.AgentSearchParams{}
			count, err := suite.sut.SearchCount(context.Background(), params)

			assert.Nil(t, err)
			assert.Equal(t, uint(0), count)
		})
	})

	t.Run("Given agents in the database", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return correct count", func(t *testing.T) {
			t.Parallel()

			suite := WhenSearchingAgentsBeforeEach(t)
			defer WhenSearchingAgentsAfterEach(suite)

			for i := 0; i < 3; i++ {
				agent := model_builders.NewAgentBuilder().Build()
				suite.sut.Save(context.Background(), agent)
			}

			params := repositories.AgentSearchParams{}
			count, err := suite.sut.SearchCount(context.Background(), params)

			assert.Nil(t, err)
			assert.Equal(t, uint(3), count)
		})
	})

	t.Run("Given minted and proposal agents in the database", func(t *testing.T) {
		t.Parallel()

		saveMixedAgents := func(suite *WhenCountingSearchAgentsTestingSuite) {
			for i := 0; i < 2; i++ {
				suite.sut.Save(context.Background(), model_builders.NewAgentBuilder().Build())
			}
			suite.sut.Save(context.Background(), model_builders.NewAgentBuilder().WithCollectionTokenId(nil).Build())
		}

		t.Run("Given the filter selects minted agents", func(t *testing.T) {
			t.Parallel()

			t.Run("Should count only minted agents", func(t *testing.T) {
				t.Parallel()
				suite := WhenCountingSearchAgentsBeforeEach(t)
				defer WhenCountingSearchAgentsAfterEach(suite)
				saveMixedAgents(suite)

				count, err := suite.sut.SearchCount(
					context.Background(),
					repositories.AgentSearchParams{IsMinted: lo.ToPtr(true)},
				)

				assert.Nil(t, err)
				assert.Equal(t, uint(2), count)
			})
		})

		t.Run("Given the filter selects proposal agents", func(t *testing.T) {
			t.Parallel()

			t.Run("Should count only proposal agents", func(t *testing.T) {
				t.Parallel()
				suite := WhenCountingSearchAgentsBeforeEach(t)
				defer WhenCountingSearchAgentsAfterEach(suite)
				saveMixedAgents(suite)

				count, err := suite.sut.SearchCount(
					context.Background(),
					repositories.AgentSearchParams{IsMinted: lo.ToPtr(false)},
				)

				assert.Nil(t, err)
				assert.Equal(t, uint(1), count)
			})
		})
	})
}
