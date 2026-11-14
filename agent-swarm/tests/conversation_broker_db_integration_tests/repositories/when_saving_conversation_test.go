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

type WhenSavingConversationTestingSuite struct {
	dbName   string
	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.ConversationRepository
	bridges  *repositories.BridgeRepository
}

func WhenSavingConversationBeforeEach(t *testing.T) *WhenSavingConversationTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)
	return &WhenSavingConversationTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,
		sut:      repositories.NewConversationRepository(dbClient),
		bridges:  repositories.NewBridgeRepository(dbClient),
	}
}

func WhenSavingConversationAfterEach(suite *WhenSavingConversationTestingSuite) {
	container_builders.GetPostgresContainerBuilder().DropDatabase(suite.dbName)
}

func TestWhenSavingConversation(t *testing.T) {
	t.Parallel()

	t.Run("Given a conversation attached to an existing bridge", func(t *testing.T) {
		t.Parallel()

		t.Run("Should be retrievable by id", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingConversationBeforeEach(t)
			defer WhenSavingConversationAfterEach(suite)

			bridge := model_builders.NewBridgeBuilder().Build()
			_, err := suite.bridges.Save(context.Background(), bridge)
			assert.NoError(t, err)

			entity := model_builders.NewConversationBuilder().WithBridgeID(bridge.ID).Build()
			_, err = suite.sut.Save(context.Background(), entity)
			assert.NoError(t, err)

			got, err := suite.sut.GetById(context.Background(), entity.ID)
			assert.NoError(t, err)
			assert.Equal(t, entity.ID, got.ID)
			assert.Equal(t, bridge.ID, got.BridgeID)
		})
	})
}
