package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/db_clients"
	orchestrator_agent_compute_db "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_db"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type WhenSavingMCPConnectionTestingSuite struct {
	dbName   string
	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.MCPConnectionRepository
}

func WhenSavingMCPConnectionBeforeEach() *WhenSavingMCPConnectionTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}
	dbClient, _ := orchestrator_agent_compute_db.NewOrchestratorDbClient(
		connectionString,
		zap.NewNop(),
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
	)
	dbClient.Migrate()

	sut := repositories.NewMCPConnectionRepository(dbClient)

	return &WhenSavingMCPConnectionTestingSuite{
		sut:      sut,
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenSavingMCPConnectionAfterEach(suite *WhenSavingMCPConnectionTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenSavingMCPConnectionEntity(t *testing.T) {
	t.Parallel()

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return created entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingMCPConnectionBeforeEach()
			defer WhenSavingMCPConnectionAfterEach(suite)

			entityToCreate := model_builders.NewMCPConnectionBuilder().Build()
			createdEntity, err := suite.sut.Save(entityToCreate)

			assert.Nil(t, err)
			assert.Equal(t, entityToCreate, createdEntity)
		})

		t.Run("Should store entity in database", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingMCPConnectionBeforeEach()
			defer WhenSavingMCPConnectionAfterEach(suite)

			entityToCreate := model_builders.NewMCPConnectionBuilder().Build()
			createdEntity, _ := suite.sut.Save(entityToCreate)

			suite.dbClient.GetDB().First(&entityToCreate, createdEntity.ID)

			assert.Equal(t, entityToCreate, createdEntity)
			assert.NotEqual(t, uuid.Nil, createdEntity.ID)
		})
	})

	t.Run("Given entity exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return updated entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingMCPConnectionBeforeEach()
			defer WhenSavingMCPConnectionAfterEach(suite)

			entityToCreate := model_builders.NewMCPConnectionBuilder().Build()
			createdEntity, err := suite.sut.Save(entityToCreate)
			if err != nil {
				t.Fatal(err)
			}

			entityToUpdate := model_builders.NewMCPConnectionBuilder().WithID(createdEntity.ID).Build()
			updatedEntity, err := suite.sut.Save(entityToUpdate)

			assert.Nil(t, err)
			assert.Equal(t, entityToUpdate, updatedEntity)
		})

		t.Run("Should update entity in database", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingMCPConnectionBeforeEach()
			defer WhenSavingMCPConnectionAfterEach(suite)

			entityToCreate := model_builders.NewMCPConnectionBuilder().Build()
			createdEntity, _ := suite.sut.Save(entityToCreate)

			entityToUpdate := model_builders.NewMCPConnectionBuilder().WithID(createdEntity.ID).Build()
			suite.sut.Save(entityToUpdate)

			suite.dbClient.GetDB().First(&entityToCreate, createdEntity.ID)

			assert.Equal(t, entityToCreate.ID, entityToUpdate.ID)
			assert.Equal(t, entityToCreate.Name, entityToUpdate.Name)
			assert.Equal(t, entityToCreate.AgentID, entityToUpdate.AgentID)
		})
	})
}
