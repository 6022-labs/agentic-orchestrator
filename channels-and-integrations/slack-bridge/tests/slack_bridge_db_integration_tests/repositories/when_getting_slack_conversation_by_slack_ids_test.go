package repositories_test

import (
	"context"

	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type WhenGettingSlackConversationBySlackIdsTestingSuite struct {
	sut *repositories.SlackConversationRepository

	dbName   string
	dbClient *slack_bridge_db.DatabaseClient
}

func WhenGettingSlackConversationBySlackIdsBeforeEach(
	t *testing.T,
) *WhenGettingSlackConversationBySlackIdsTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)

	return &WhenGettingSlackConversationBySlackIdsTestingSuite{
		sut:      repositories.NewSlackConversationRepository(dbClient),
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenGettingSlackConversationBySlackIdsAfterEach(s *WhenGettingSlackConversationBySlackIdsTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(s.dbName)
}

func TestWhenGettingSlackConversationBySlackIds(t *testing.T) {
	t.Parallel()

	t.Run("Given the slack conversation link exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the link", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackConversationBySlackIdsBeforeEach(t)
			defer WhenGettingSlackConversationBySlackIdsAfterEach(suite)

			brokerId := uuid.New()
			_, err := suite.sut.Save(context.Background(), &models.SlackConversation{
				SlackChannelId:       "C123",
				SlackThreadTs:        "1728391823.001200",
				BrokerConversationID: brokerId,
				CreatedAt:            time.Now(),
			})
			assert.NoError(t, err)

			got, err := suite.sut.GetBySlackIds(context.Background(), "C123", "1728391823.001200")

			assert.NoError(t, err)
			assert.Equal(t, brokerId, got.BrokerConversationID)
		})
	})

	t.Run("Given the slack conversation link does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackConversationBySlackIdsBeforeEach(t)
			defer WhenGettingSlackConversationBySlackIdsAfterEach(suite)

			entity, err := suite.sut.GetBySlackIds(context.Background(), "C-missing", "0.0")

			assert.NoError(t, err)
			assert.Nil(t, entity)
		})
	})
}
