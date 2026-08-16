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

type WhenGettingSlackConversationByBrokerConversationIdTestingSuite struct {
	sut *repositories.SlackConversationRepository

	dbName   string
	dbClient *slack_bridge_db.DatabaseClient
}

func WhenGettingSlackConversationByBrokerConversationIdBeforeEach() *WhenGettingSlackConversationByBrokerConversationIdTestingSuite {
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

	return &WhenGettingSlackConversationByBrokerConversationIdTestingSuite{
		sut:      repositories.NewSlackConversationRepository(dbClient),
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenGettingSlackConversationByBrokerConversationIdAfterEach(
	s *WhenGettingSlackConversationByBrokerConversationIdTestingSuite,
) {
	container_builders.DestroyPostgreSQLDatabase(s.dbName)
}

func TestWhenGettingSlackConversationByBrokerConversationId(t *testing.T) {
	t.Parallel()

	t.Run("Given the slack conversation link exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the link", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackConversationByBrokerConversationIdBeforeEach()
			defer WhenGettingSlackConversationByBrokerConversationIdAfterEach(suite)

			brokerId := uuid.New()
			_, err := suite.sut.Save(&models.SlackConversation{
				SlackChannelId:       "C123",
				SlackThreadTs:        "1728391823.001200",
				BrokerConversationID: brokerId,
				CreatedAt:            time.Now(),
			})
			assert.NoError(t, err)

			got, err := suite.sut.GetByBrokerConversationId(brokerId)

			assert.NoError(t, err)
			assert.Equal(t, "C123", got.SlackChannelId)
			assert.Equal(t, "1728391823.001200", got.SlackThreadTs)
		})
	})

	t.Run("Given the slack conversation link does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackConversationByBrokerConversationIdBeforeEach()
			defer WhenGettingSlackConversationByBrokerConversationIdAfterEach(suite)

			entity, err := suite.sut.GetByBrokerConversationId(uuid.New())

			assert.NoError(t, err)
			assert.Nil(t, entity)
		})
	})
}
