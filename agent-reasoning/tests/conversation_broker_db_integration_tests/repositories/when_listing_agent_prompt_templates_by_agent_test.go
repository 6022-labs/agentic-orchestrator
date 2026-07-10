package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
)

type WhenListingAgentPromptTemplatesByAgentTestingSuite struct {
	dbName   string
	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentPromptTemplateRepository
}

func WhenListingAgentPromptTemplatesByAgentBeforeEach() *WhenListingAgentPromptTemplatesByAgentTestingSuite {
	dbName, dbClient := newMigratedDatabase()
	return &WhenListingAgentPromptTemplatesByAgentTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,
		sut:      repositories.NewAgentPromptTemplateRepository(dbClient),
	}
}

func WhenListingAgentPromptTemplatesByAgentAfterEach(suite *WhenListingAgentPromptTemplatesByAgentTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenListingAgentPromptTemplatesByAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given no templates for the agent", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an empty list", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentPromptTemplatesByAgentBeforeEach()
			defer WhenListingAgentPromptTemplatesByAgentAfterEach(suite)

			templates, err := suite.sut.ListByAgent("0xabc", 1)

			assert.NoError(t, err)
			assert.Empty(t, templates)
		})
	})

	t.Run("Given templates for two different agents", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return only the requested agent's templates ordered by key", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentPromptTemplatesByAgentBeforeEach()
			defer WhenListingAgentPromptTemplatesByAgentAfterEach(suite)

			first := model_builders.NewAgentPromptTemplateBuilder().
				WithCollectionAddress("0xabc").
				WithCollectionTokenId(1).
				WithKey(enums.AgentPromptTemplateKeyCloneSystem).
				Build()
			second := model_builders.NewAgentPromptTemplateBuilder().
				WithCollectionAddress("0xabc").
				WithCollectionTokenId(1).
				WithKey(enums.AgentPromptTemplateKeyCloneDefault).
				Build()
			other := model_builders.NewAgentPromptTemplateBuilder().
				WithCollectionAddress("0xother").
				WithCollectionTokenId(2).
				WithKey(enums.AgentPromptTemplateKeyCloneSystem).
				Build()

			_, err := suite.sut.Save(first)
			assert.NoError(t, err)
			_, err = suite.sut.Save(second)
			assert.NoError(t, err)
			_, err = suite.sut.Save(other)
			assert.NoError(t, err)

			templates, err := suite.sut.ListByAgent("0xabc", 1)

			assert.NoError(t, err)
			assert.Len(t, templates, 2)
			assert.Equal(t, enums.AgentPromptTemplateKeyCloneDefault, templates[0].Key)
			assert.Equal(t, enums.AgentPromptTemplateKeyCloneSystem, templates[1].Key)
		})
	})
}
