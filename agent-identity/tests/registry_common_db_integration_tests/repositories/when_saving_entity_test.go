package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type WhenSavingEntityTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentRepository
}

func WhenSavingEntityBeforeEach(t *testing.T) *WhenSavingEntityTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)

	sut := repositories.NewAgentRepository(dbClient)

	return &WhenSavingEntityTestingSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenSavingEntityAfterEach(suite *WhenSavingEntityTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenSavingEntity(t *testing.T) {
	t.Parallel()

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return created entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingEntityBeforeEach(t)
			defer WhenSavingEntityAfterEach(suite)

			entityToCreate := model_builders.NewAgentBuilder().Build()
			createdEntity, err := suite.sut.Save(context.Background(), entityToCreate)

			assert.Nil(t, err)
			assert.Equal(t, entityToCreate, createdEntity)
		})

		t.Run("Should store entity in database", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingEntityBeforeEach(t)
			defer WhenSavingEntityAfterEach(suite)

			entityToCreate := model_builders.NewAgentBuilder().Build()
			createdEntity, _ := suite.sut.Save(context.Background(), entityToCreate)

			fetchedEntity, _ := suite.sut.GetById(context.Background(), createdEntity.ID)

			assert.Equal(t, createdEntity.ID, fetchedEntity.ID)
			assert.Equal(t, createdEntity.Name, fetchedEntity.Name)
			assert.NotEqual(t, uuid.Nil, createdEntity.ID)
		})
	})

	t.Run("Given entity exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return updated entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingEntityBeforeEach(t)
			defer WhenSavingEntityAfterEach(suite)

			entityToCreate := model_builders.NewAgentBuilder().Build()
			createdEntity, _ := suite.sut.Save(context.Background(), entityToCreate)

			entityToUpdate := model_builders.NewAgentBuilder().WithID(createdEntity.ID).Build()
			updatedEntity, err := suite.sut.Save(context.Background(), entityToUpdate)

			assert.Nil(t, err)
			assert.Equal(t, entityToUpdate, updatedEntity)
		})

		t.Run("Should update entity in database", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingEntityBeforeEach(t)
			defer WhenSavingEntityAfterEach(suite)

			entityToCreate := model_builders.NewAgentBuilder().Build()
			createdEntity, _ := suite.sut.Save(context.Background(), entityToCreate)

			entityToUpdate := model_builders.NewAgentBuilder().WithID(createdEntity.ID).Build()
			suite.sut.Save(context.Background(), entityToUpdate)

			entityToCreate, _ = suite.sut.GetById(context.Background(), createdEntity.ID)

			assert.Equal(t, entityToCreate.ID, entityToUpdate.ID)
			assert.Equal(t, entityToCreate.Name, entityToUpdate.Name)
			assert.Equal(t, entityToCreate.Owner, entityToUpdate.Owner)
			assert.Equal(t, entityToCreate.CollectionTokenId, entityToUpdate.CollectionTokenId)

			assert.NotNil(t, entityToCreate.CreatedAt)
			assert.NotNil(t, entityToCreate.UpdatedAt)
		})
	})
}
