package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
)

type WhenSavingAgentPromptTemplateTestingSuite struct {
	dbName   string
	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentPromptTemplateRepository
}

func WhenSavingAgentPromptTemplateBeforeEach(t *testing.T) *WhenSavingAgentPromptTemplateTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)
	return &WhenSavingAgentPromptTemplateTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,
		sut:      repositories.NewAgentPromptTemplateRepository(dbClient),
	}
}

func WhenSavingAgentPromptTemplateAfterEach(suite *WhenSavingAgentPromptTemplateTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenSavingAgentPromptTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Given a new template", func(t *testing.T) {
		t.Parallel()

		t.Run("Should be retrievable by agent and key", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			defer WhenSavingAgentPromptTemplateAfterEach(suite)

			entity := model_builders.NewAgentPromptTemplateBuilder().Build()
			_, err := suite.sut.Save(context.Background(), entity)
			assert.NoError(t, err)

			got, err := suite.sut.GetByAgentAndKey(
				context.Background(),
				entity.CollectionAddress,
				entity.CollectionTokenId,
				entity.Key,
			)
			assert.NoError(t, err)
			assert.Equal(t, entity.ID, got.ID)
			assert.Equal(t, entity.Content, got.Content)
		})
	})

	t.Run("Given an existing template", func(t *testing.T) {
		t.Parallel()

		t.Run("Should overwrite the content on save", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			defer WhenSavingAgentPromptTemplateAfterEach(suite)

			entity := model_builders.NewAgentPromptTemplateBuilder().WithContent("old").Build()
			_, err := suite.sut.Save(context.Background(), entity)
			assert.NoError(t, err)

			entity.Content = "new"
			_, err = suite.sut.Save(context.Background(), entity)
			assert.NoError(t, err)

			got, err := suite.sut.GetByAgentAndKey(
				context.Background(),
				entity.CollectionAddress,
				entity.CollectionTokenId,
				entity.Key,
			)
			assert.NoError(t, err)
			assert.Equal(t, "new", got.Content)
		})
	})
}
