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

type WhenSavingSlackConversationTestingSuite struct {
	sut *repositories.SlackConversationRepository

	dbName   string
	dbClient *slack_bridge_db.DatabaseClient
}

func WhenSavingSlackConversationBeforeEach(t *testing.T) *WhenSavingSlackConversationTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)

	return &WhenSavingSlackConversationTestingSuite{
		sut:      repositories.NewSlackConversationRepository(dbClient),
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenSavingSlackConversationAfterEach(s *WhenSavingSlackConversationTestingSuite) {
	container_builders.GetPostgresContainerBuilder().DropDatabase(s.dbName)
}

func TestWhenSavingSlackConversation(t *testing.T) {
	t.Parallel()

	t.Run("Given the slack conversation link does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should persist the link", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingSlackConversationBeforeEach(t)
			defer WhenSavingSlackConversationAfterEach(suite)

			entity := &models.SlackConversation{
				SlackChannelId:       "C123",
				SlackThreadTs:        "1728391823.001200",
				BrokerConversationID: uuid.New(),
				CreatedAt:            time.Now(),
			}

			_, err := suite.sut.Save(context.Background(), entity)

			assert.NoError(t, err)
		})
	})
}
