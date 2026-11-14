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

type WhenSavingSlackMessageTestingSuite struct {
	sut *repositories.SlackMessageRepository

	dbName   string
	dbClient *slack_bridge_db.DatabaseClient
}

func WhenSavingSlackMessageBeforeEach(t *testing.T) *WhenSavingSlackMessageTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)

	return &WhenSavingSlackMessageTestingSuite{
		sut:      repositories.NewSlackMessageRepository(dbClient),
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenSavingSlackMessageAfterEach(s *WhenSavingSlackMessageTestingSuite) {
	container_builders.GetPostgresContainerBuilder().DropDatabase(s.dbName)
}

func TestWhenSavingSlackMessage(t *testing.T) {
	t.Parallel()

	t.Run("Given the slack message link does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should persist the link", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingSlackMessageBeforeEach(t)
			defer WhenSavingSlackMessageAfterEach(suite)

			entity := &models.SlackMessage{
				SlackChannelId:  "C123",
				SlackMessageTs:  "1728391823.001200",
				BrokerMessageID: uuid.New(),
				CreatedAt:       time.Now(),
			}

			_, err := suite.sut.Save(context.Background(), entity)

			assert.NoError(t, err)
		})
	})
}
