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
	"gorm.io/gorm"
)

type WhenGettingSlackMessageLinkBySlackIdsTestingSuite struct {
	sut *repositories.SlackMessageLinkRepository

	dbName   string
	dbClient *slack_bridge_db.DatabaseClient
}

func WhenGettingSlackMessageLinkBySlackIdsBeforeEach() *WhenGettingSlackMessageLinkBySlackIdsTestingSuite {
	dbName, connStr := container_builders.CreatePostgreSQLDatabase()

	folder, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}
	dbClient, _ := slack_bridge_db.NewDatabaseClient(connStr, zap.NewNop(), &db_settings.MigrationsSettings{MigrationsFolderPath: *folder})
	if err := dbClient.Migrate(); err != nil {
		panic(err)
	}

	return &WhenGettingSlackMessageLinkBySlackIdsTestingSuite{
		sut:      repositories.NewSlackMessageLinkRepository(dbClient),
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenGettingSlackMessageLinkBySlackIdsAfterEach(s *WhenGettingSlackMessageLinkBySlackIdsTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(s.dbName)
}

func TestWhenGettingSlackMessageLinkBySlackIds(t *testing.T) {
	t.Parallel()

	t.Run("Given the slack message link exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the link", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackMessageLinkBySlackIdsBeforeEach()
			defer WhenGettingSlackMessageLinkBySlackIdsAfterEach(suite)

			brokerId := uuid.New()
			_, err := suite.sut.Save(&models.SlackMessageLink{
				SlackChannelId:  "C123",
				SlackMessageTs:  "1728391823.001200",
				BrokerMessageID: brokerId,
				CreatedAt:       time.Now(),
			})
			assert.NoError(t, err)

			got, err := suite.sut.GetBySlackIds("C123", "1728391823.001200")

			assert.NoError(t, err)
			assert.Equal(t, brokerId, got.BrokerMessageID)
		})
	})

	t.Run("Given the slack message link does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackMessageLinkBySlackIdsBeforeEach()
			defer WhenGettingSlackMessageLinkBySlackIdsAfterEach(suite)

			_, err := suite.sut.GetBySlackIds("C-missing", "0.0")

			assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		})
	})
}
