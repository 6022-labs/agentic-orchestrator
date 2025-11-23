package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/db_clients"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_db"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_db/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_db/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingAgentsByEnsDomainsTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentRepository
}

func WhenGettingAgentsByEnsDomainsBeforeEach() *WhenGettingAgentsByEnsDomainsTestingSuite {
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

			suite := WhenGettingAgentsByEnsDomainsBeforeEach()
			defer WhenGettingAgentsByEnsDomainsAfterEach(suite)

			existingEntity := model_builders.NewAgentBuilder().Build()
			createdEntity, err := suite.sut.Save(existingEntity)
			if err != nil {
				t.Errorf("Error creating entity: %v", err)
			}

			fetchedEntities, err := suite.sut.GetByEnsDomains([]string{createdEntity.EnsDomain})

			assert.Nil(t, err)
			assert.Equal(t, 1, len(fetchedEntities))
		})
	})

	t.Run("Given no entities exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return empty list", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentsByEnsDomainsBeforeEach()
			defer WhenGettingAgentsByEnsDomainsAfterEach(suite)

			fetchedEntities, err := suite.sut.GetByEnsDomains([]string{"nonexistent.ens"})

			assert.Nil(t, err)
			assert.Equal(t, 0, len(fetchedEntities))
		})
	})
}
