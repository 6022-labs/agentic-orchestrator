package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
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

type WhenGettingAgentsByEnsDomainsTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentRepository
}

func WhenGettingAgentsByEnsDomainsBeforeEach(t *testing.T) *WhenGettingAgentsByEnsDomainsTestingSuite {
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

	return &WhenGettingAgentsByEnsDomainsTestingSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenGettingAgentsByEnsDomainsAfterEach(suite *WhenGettingAgentsByEnsDomainsTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingAgentsByEnsDomains(t *testing.T) {
	t.Parallel()

	t.Run("Given entities exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return entities", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentsByEnsDomainsBeforeEach(t)
			defer WhenGettingAgentsByEnsDomainsAfterEach(suite)

			existingEntity := model_builders.NewAgentBuilder().Build()
			createdEntity, err := suite.sut.Save(context.Background(), existingEntity)
			if err != nil {
				t.Errorf("Error creating entity: %v", err)
			}

			fetchedEntities, err := suite.sut.GetByEnsDomains(context.Background(), []string{createdEntity.EnsDomain})

			assert.Nil(t, err)
			assert.Equal(t, 1, len(fetchedEntities))
		})
	})

	t.Run("Given no entities exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return empty list", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentsByEnsDomainsBeforeEach(t)
			defer WhenGettingAgentsByEnsDomainsAfterEach(suite)

			fetchedEntities, err := suite.sut.GetByEnsDomains(context.Background(), []string{"nonexistent.ens"})

			assert.Nil(t, err)
			assert.Equal(t, 0, len(fetchedEntities))
		})
	})
}
