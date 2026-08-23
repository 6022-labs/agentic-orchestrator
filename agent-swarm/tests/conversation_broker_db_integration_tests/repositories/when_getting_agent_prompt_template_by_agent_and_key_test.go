package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
)

type WhenGettingAgentPromptTemplateByAgentAndKeyTestingSuite struct {
	dbName   string
	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentPromptTemplateRepository
}

func WhenGettingAgentPromptTemplateByAgentAndKeyBeforeEach(
	t *testing.T,
) *WhenGettingAgentPromptTemplateByAgentAndKeyTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)
	return &WhenGettingAgentPromptTemplateByAgentAndKeyTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,
		sut:      repositories.NewAgentPromptTemplateRepository(dbClient),
	}
}

func WhenGettingAgentPromptTemplateByAgentAndKeyAfterEach(
	suite *WhenGettingAgentPromptTemplateByAgentAndKeyTestingSuite,
) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingAgentPromptTemplateByAgentAndKey(t *testing.T) {
	t.Parallel()

	t.Run("Given no template for the agent and key", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil without error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateByAgentAndKeyBeforeEach(t)
			defer WhenGettingAgentPromptTemplateByAgentAndKeyAfterEach(suite)

			template, err := suite.sut.GetByAgentAndKey(
				context.Background(),
				"0xabc",
				1,
				enums.AgentPromptTemplateKeyCloneSystem,
			)

			assert.Nil(t, err)
			assert.Nil(t, template)
		})
	})

	t.Run("Given templates for several keys of the same agent", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the one matching the key", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateByAgentAndKeyBeforeEach(t)
			defer WhenGettingAgentPromptTemplateByAgentAndKeyAfterEach(suite)

			system := model_builders.NewAgentPromptTemplateBuilder().
				WithCollectionAddress("0xabc").
				WithCollectionTokenId(1).
				WithKey(enums.AgentPromptTemplateKeyCloneSystem).
				Build()
			defaultPrompt := model_builders.NewAgentPromptTemplateBuilder().
				WithCollectionAddress("0xabc").
				WithCollectionTokenId(1).
				WithKey(enums.AgentPromptTemplateKeyCloneDefault).
				Build()

			_, err := suite.sut.Save(context.Background(), system)
			assert.NoError(t, err)
			_, err = suite.sut.Save(context.Background(), defaultPrompt)
			assert.NoError(t, err)

			got, err := suite.sut.GetByAgentAndKey(
				context.Background(),
				"0xabc",
				1,
				enums.AgentPromptTemplateKeyCloneDefault,
			)

			assert.NoError(t, err)
			assert.Equal(t, defaultPrompt.ID, got.ID)
			assert.Equal(t, enums.AgentPromptTemplateKeyCloneDefault, got.Key)
		})
	})
}
