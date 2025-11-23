package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	db_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"go.uber.org/zap"
)

type WhenSavingMultipleSlackMessagesTestingSuite struct {
	dbName string

	dbClient *slack_bridge_db.DatabaseClient

	slackThreadRepository  *repositories.SlackThreadRepository
	slackChannelRepository *repositories.SlackChannelRepository

	sut *repositories.SlackMessageRepository
}

func WhenSavingMultipleSlackMessagesBeforeEach() *WhenSavingMultipleSlackMessagesTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbClient, _ := slack_bridge_db.NewDatabaseClient(
		connectionString,
		zap.NewNop(),
		&db_settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
	)
	err = dbClient.Migrate()
	if err != nil {
		panic(err)
	}

	processingSettings := &settings.ProcessingSettings{
		AgentProcessingInterval: 20,
		FacilitatorTickInterval: 30,
	}

	sut := repositories.NewSlackMessageRepository(dbClient)
	slackThreadRepository := repositories.NewSlackThreadRepository(processingSettings, dbClient)
	slackChannelRepository := repositories.NewSlackChannelRepository(dbClient)

	return &WhenSavingMultipleSlackMessagesTestingSuite{
		sut: sut,

		dbName:   dbName,
		dbClient: dbClient,

		slackThreadRepository:  slackThreadRepository,
		slackChannelRepository: slackChannelRepository,
	}
}

func WhenSavingMultipleSlackMessagesAfterEach(suite *WhenSavingMultipleSlackMessagesTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenSavingMultipleSlackMessages(t *testing.T) {
	t.Parallel()

	t.Run("Given slack message does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should save messages", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingMultipleSlackMessagesBeforeEach()
			defer WhenSavingMultipleSlackMessagesAfterEach(suite)

			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			_, err := suite.slackChannelRepository.Save(slackChannel)
			if err != nil {
				t.Fatalf("Failed to save slack channel: %v", err)
			}

			slackThread := model_builders.NewSlackThreadBuilder().
				WithChannelId(slackChannel.Id).
				Build()
			_, err = suite.slackThreadRepository.Save(slackThread)
			if err != nil {
				t.Fatalf("Failed to save slack thread: %v", err)
			}

			messages := []*models.SlackMessage{
				model_builders.NewSlackMessageBuilder().
					WithThreadId(slackThread.Id).
					WithChannelId(slackChannel.Id).
					Build(),
				model_builders.NewSlackMessageBuilder().
					WithThreadId(slackThread.Id).
					WithChannelId(slackChannel.Id).
					Build(),
			}

			savedMessages, err := suite.sut.SaveMany(messages)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if len(savedMessages) != 2 {
				t.Errorf("Expected 2 messages, got %d", len(savedMessages))
			}
		})
	})
}
