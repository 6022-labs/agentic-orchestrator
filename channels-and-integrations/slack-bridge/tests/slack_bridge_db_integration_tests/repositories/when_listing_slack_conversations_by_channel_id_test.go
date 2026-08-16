package repositories_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	db_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenListingSlackConversationsByChannelIdTestingSuite struct {
	sut *repositories.SlackConversationRepository

	dbName   string
	dbClient *slack_bridge_db.DatabaseClient
}

func WhenListingSlackConversationsByChannelIdBeforeEach() *WhenListingSlackConversationsByChannelIdTestingSuite {
	dbName, connStr := container_builders.CreatePostgreSQLDatabase()

	folder, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}
	dbClient, _ := slack_bridge_db.NewDatabaseClient(
		connStr,
		zap.NewNop(),
		&db_settings.MigrationsSettings{MigrationsFolderPath: *folder},
	)
	if err := dbClient.Migrate(); err != nil {
		panic(err)
	}

	return &WhenListingSlackConversationsByChannelIdTestingSuite{
		sut:      repositories.NewSlackConversationRepository(dbClient),
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenListingSlackConversationsByChannelIdAfterEach(s *WhenListingSlackConversationsByChannelIdTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(s.dbName)
}

func TestWhenListingSlackConversationsByChannelId(t *testing.T) {
	t.Parallel()

	t.Run("Given multiple links exist for the channel and one for another channel", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return only the links for that channel", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingSlackConversationsByChannelIdBeforeEach()
			defer WhenListingSlackConversationsByChannelIdAfterEach(suite)

			channelId := "C-target"
			otherChannelId := "C-other"

			_, err := suite.sut.Save(&models.SlackConversation{
				SlackChannelId:       channelId,
				SlackThreadTs:        "1.0",
				BrokerConversationID: uuid.New(),
				CreatedAt:            time.Now(),
			})
			assert.NoError(t, err)
			_, err = suite.sut.Save(&models.SlackConversation{
				SlackChannelId:       channelId,
				SlackThreadTs:        "2.0",
				BrokerConversationID: uuid.New(),
				CreatedAt:            time.Now(),
			})
			assert.NoError(t, err)
			_, err = suite.sut.Save(&models.SlackConversation{
				SlackChannelId:       otherChannelId,
				SlackThreadTs:        "3.0",
				BrokerConversationID: uuid.New(),
				CreatedAt:            time.Now(),
			})
			assert.NoError(t, err)

			got, err := suite.sut.ListByChannelId(channelId)

			assert.NoError(t, err)
			assert.Len(t, got, 2)
			for _, link := range got {
				assert.Equal(t, channelId, link.SlackChannelId)
			}
		})
	})

	t.Run("Given no link exists for the channel", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an empty list", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingSlackConversationsByChannelIdBeforeEach()
			defer WhenListingSlackConversationsByChannelIdAfterEach(suite)

			got, err := suite.sut.ListByChannelId("C-empty")

			assert.NoError(t, err)
			assert.Empty(t, got)
		})
	})
}
