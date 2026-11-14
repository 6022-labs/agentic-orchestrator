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

type WhenGettingSlackMessageBySlackIdsTestingSuite struct {
	sut *repositories.SlackMessageRepository

	dbName   string
	dbClient *slack_bridge_db.DatabaseClient
}

func WhenGettingSlackMessageBySlackIdsBeforeEach(t *testing.T) *WhenGettingSlackMessageBySlackIdsTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)

	return &WhenGettingSlackMessageBySlackIdsTestingSuite{
		sut:      repositories.NewSlackMessageRepository(dbClient),
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenGettingSlackMessageBySlackIdsAfterEach(s *WhenGettingSlackMessageBySlackIdsTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(s.dbName)
}

func TestWhenGettingSlackMessageBySlackIds(t *testing.T) {
	t.Parallel()

	t.Run("Given the slack message link exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the link", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackMessageBySlackIdsBeforeEach(t)
			defer WhenGettingSlackMessageBySlackIdsAfterEach(suite)

			brokerId := uuid.New()
			_, err := suite.sut.Save(context.Background(), &models.SlackMessage{
				SlackChannelId:  "C123",
				SlackMessageTs:  "1728391823.001200",
				BrokerMessageID: brokerId,
				CreatedAt:       time.Now(),
			})
			assert.NoError(t, err)

			got, err := suite.sut.GetBySlackIds(context.Background(), "C123", "1728391823.001200")

			assert.NoError(t, err)
			assert.Equal(t, brokerId, got.BrokerMessageID)
		})
	})

	t.Run("Given the slack message link does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackMessageBySlackIdsBeforeEach(t)
			defer WhenGettingSlackMessageBySlackIdsAfterEach(suite)

			entity, err := suite.sut.GetBySlackIds(context.Background(), "C-missing", "0.0")

			assert.NoError(t, err)
			assert.Nil(t, entity)
		})
	})
}
