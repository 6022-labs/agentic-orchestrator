package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingAllAgentsTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentRepository
}

func WhenGettingAllAgentsBeforeEach() *WhenGettingAllAgentsTestingSuite {
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

	sut := repositories.NewAgentRepository(dbClient)

	return &WhenGettingAllAgentsTestingSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenGettingAllAgentsAfterEach(suite *WhenGettingAllAgentsTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingAllAgents(t *testing.T) {
	t.Parallel()

	t.Run("Given entity exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAllAgentsBeforeEach()
			defer WhenGettingAllAgentsAfterEach(suite)

			entity := model_builders.NewAgentBuilder().Build()
			createdEntity, _ := suite.sut.Save(entity)

			fetchedEntity, err := suite.sut.GetById(createdEntity.ID)

			assert.Nil(t, err)
			assert.Equal(t, createdEntity.ID, fetchedEntity.ID)
			assert.Equal(t, createdEntity.Name, fetchedEntity.Name)
		})
	})

	t.Run("Given entities exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return entities", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAllAgentsBeforeEach()
			defer WhenGettingAllAgentsAfterEach(suite)

			createdEntities := []*models.Agent{}
			for i := 0; i < 3; i++ {
				entity := model_builders.NewAgentBuilder().Build()
				createdEntity, _ := suite.sut.Save(entity)
				createdEntities = append(createdEntities, createdEntity)
			}

			fetchedEntities, err := suite.sut.GetAll()

			assert.Nil(t, err)
			assert.Len(t, fetchedEntities, len(createdEntities))

			for i, createdEntity := range createdEntities {
				fetchedEntity := fetchedEntities[i]
				assert.Equal(t, createdEntity.ID, fetchedEntity.ID)
				assert.Equal(t, createdEntity.Name, fetchedEntity.Name)
			}
		})
	})

	t.Run("Given no entity in persistence", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil without error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAllAgentsBeforeEach()
			defer WhenGettingAllAgentsAfterEach(suite)

			nonExistentID := uuid.New()
			fetchedAgent, err := suite.sut.GetById(nonExistentID)

			assert.Nil(t, err)
			assert.Nil(t, fetchedAgent)
		})
	})
}
