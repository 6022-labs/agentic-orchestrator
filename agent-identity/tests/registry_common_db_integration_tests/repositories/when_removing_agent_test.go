package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	db_metrics_mocks "github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_db_mocks/metrics_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenRemovingAgentTestingSuite struct {
	dbName string

	sut *repositories.AgentRepository
}

func WhenRemovingAgentBeforeEach(t *testing.T) *WhenRemovingAgentTestingSuite {
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

	sut := repositories.NewAgentRepository(dbClient)

	return &WhenRemovingAgentTestingSuite{
		sut:    sut,
		dbName: dbName,
	}
}

func WhenRemovingAgentAfterEach(suite *WhenRemovingAgentTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenRemovingAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given agent does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentBeforeEach(t)
			defer WhenRemovingAgentAfterEach(suite)

			nonExistentAgent := model_builders.NewAgentBuilder().Build()

			err := suite.sut.Remove(context.Background(), nonExistentAgent)

			assert.Nil(t, err)
		})
	})

	t.Run("Given agent exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentBeforeEach(t)
			defer WhenRemovingAgentAfterEach(suite)

			agent := model_builders.NewAgentBuilder().Build()
			_, err := suite.sut.Save(context.Background(), agent)
			if err != nil {
				t.Fatalf("Error saving agent: %v", err)
			}

			err = suite.sut.Remove(context.Background(), agent)

			assert.Nil(t, err)

			fetchedAgent, err := suite.sut.GetById(context.Background(), agent.ID)

			assert.Nil(t, err)
			assert.Nil(t, fetchedAgent)
		})
	})
}
