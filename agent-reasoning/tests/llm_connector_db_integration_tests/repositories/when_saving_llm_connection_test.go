package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/db_clients"
	orchestrator_agent_llm_connector_db "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector_db"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type WhenSavingLLMConnectionTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.LLMConnectionRepository
}

func WhenSavingLLMConnectionBeforeEach() *WhenSavingLLMConnectionTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}
	dbClient, _ := orchestrator_agent_llm_connector_db.NewOrchestratorDbClient(
		connectionString,
		zap.NewNop(),
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
	)
	dbClient.Migrate()

	sut := repositories.NewLLMConnectionRepository(dbClient)

	return &WhenSavingLLMConnectionTestingSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenSavingLLMConnectionAfterEach(suite *WhenSavingLLMConnectionTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenSavingEntity(t *testing.T) {
	t.Parallel()

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return created entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingLLMConnectionBeforeEach()
			defer WhenSavingLLMConnectionAfterEach(suite)

			entityToCreate := model_builders.NewLLMConnectionBuilder().Build()
			createdEntity, err := suite.sut.Save(entityToCreate)

			assert.Nil(t, err)
			assert.Equal(t, entityToCreate, createdEntity)
		})

		t.Run("Should store entity in database", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingLLMConnectionBeforeEach()
			defer WhenSavingLLMConnectionAfterEach(suite)

			entityToCreate := model_builders.NewLLMConnectionBuilder().Build()
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

			suite := WhenSavingLLMConnectionBeforeEach()
			defer WhenSavingLLMConnectionAfterEach(suite)

			entityToCreate := model_builders.NewLLMConnectionBuilder().Build()
			createdEntity, _ := suite.sut.Save(entityToCreate)

			entityToUpdate := model_builders.NewLLMConnectionBuilder().WithID(createdEntity.ID).Build()
			updatedEntity, err := suite.sut.Save(entityToUpdate)

			assert.Nil(t, err)
			assert.Equal(t, entityToUpdate, updatedEntity)
		})

		t.Run("Should update entity in database", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingLLMConnectionBeforeEach()
			defer WhenSavingLLMConnectionAfterEach(suite)

			entityToCreate := model_builders.NewLLMConnectionBuilder().Build()
			createdEntity, _ := suite.sut.Save(entityToCreate)

			entityToUpdate := model_builders.NewLLMConnectionBuilder().WithID(createdEntity.ID).Build()
			suite.sut.Save(entityToUpdate)

			suite.dbClient.GetDB().First(&entityToCreate, createdEntity.ID)

			assert.Equal(t, entityToCreate.ID, entityToUpdate.ID)
			assert.Equal(t, entityToCreate.BaseUrl, entityToUpdate.BaseUrl)
			assert.Equal(t, entityToCreate.AgentID, entityToUpdate.AgentID)
			assert.Equal(t, entityToCreate.Platform, entityToUpdate.Platform)
			assert.Equal(t, entityToCreate.CreatedAt, entityToUpdate.CreatedAt)
		})
	})
}
