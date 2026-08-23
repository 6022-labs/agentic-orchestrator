package repositories_test

import (
	"context"

	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	db_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	db_metrics_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_db_mocks/metrics_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenSavingSlackConversationTestingSuite struct {
	sut *repositories.SlackConversationRepository

	dbName   string
	dbClient *slack_bridge_db.DatabaseClient
}

func WhenSavingSlackConversationBeforeEach(t *testing.T) *WhenSavingSlackConversationTestingSuite {
	dbName, connStr := container_builders.CreatePostgreSQLDatabase()

	folder, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}
	dbQueryMetrics := db_metrics_mocks.NewMockDbQueryMetricsInterface(gomock.NewController(t))
	dbQueryMetrics.EXPECT().
		RecordQuery(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes()

	dbClient, _ := slack_bridge_db.NewDatabaseClient(
		connStr,
		zap.NewNop(),
		&db_settings.MigrationsSettings{MigrationsFolderPath: *folder},
		dbQueryMetrics,
	)
	if err := dbClient.Migrate(); err != nil {
		panic(err)
	}

	return &WhenSavingSlackConversationTestingSuite{
		sut:      repositories.NewSlackConversationRepository(dbClient),
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenSavingSlackConversationAfterEach(s *WhenSavingSlackConversationTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(s.dbName)
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
