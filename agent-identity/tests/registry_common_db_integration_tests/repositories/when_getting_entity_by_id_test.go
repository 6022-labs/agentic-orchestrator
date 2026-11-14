package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

type WhenGettingEntityByIdTestingSuite struct {
	dbName string

	sut *repositories.AgentRepository
}

func WhenGettingEntityByIdBeforeEach(t *testing.T) *WhenGettingEntityByIdTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)

	sut := repositories.NewAgentRepository(dbClient)

	return &WhenGettingEntityByIdTestingSuite{
		sut:    sut,
		dbName: dbName,
	}
}

func WhenGettingEntityByIdBeforeEachAfterEach(suite *WhenGettingEntityByIdTestingSuite) {
	container_builders.GetPostgresContainerBuilder().DropDatabase(suite.dbName)
}

func TestWhenGettingEntityById(t *testing.T) {
	t.Parallel()

	t.Run("Given entity exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingEntityByIdBeforeEach(t)
			defer WhenGettingEntityByIdBeforeEachAfterEach(suite)

			entity := model_builders.NewAgentBuilder().Build()
			createdEntity, _ := suite.sut.Save(context.Background(), entity)

			fetchedEntity, err := suite.sut.GetById(context.Background(), createdEntity.ID)

			assert.Nil(t, err)
			assert.Equal(t, createdEntity.ID, fetchedEntity.ID)
			assert.Equal(t, createdEntity.Name, fetchedEntity.Name)
		})
	})

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil without error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingEntityByIdBeforeEach(t)
			defer WhenGettingEntityByIdBeforeEachAfterEach(suite)

			nonExistentID := uuid.New()
			fetchedAgent, err := suite.sut.GetById(context.Background(), nonExistentID)

			assert.Nil(t, err)
			assert.Nil(t, fetchedAgent)
		})
	})
}
