package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	db_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"go.uber.org/zap"
)

type WhenGettingMessagesByThreadIdAndChannelIdTestingSuite struct {
	dbName string

	dbClient *slack_bridge_db.DatabaseClient

	slackThreadRepository  *repositories.SlackThreadRepository
	slackChannelRepository *repositories.SlackChannelRepository

	sut *repositories.SlackMessageRepository
}

func WhenGettingMessagesByThreadIdAndChannelIdBeforeEach() *WhenGettingMessagesByThreadIdAndChannelIdTestingSuite {
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

	sut := repositories.NewSlackMessageRepository(dbClient)
	slackThreadRepository := repositories.NewSlackThreadRepository(dbClient)
	slackChannelRepository := repositories.NewSlackChannelRepository(dbClient)

	return &WhenGettingMessagesByThreadIdAndChannelIdTestingSuite{
		sut: sut,

		dbName:   dbName,
		dbClient: dbClient,

		slackThreadRepository:  slackThreadRepository,
		slackChannelRepository: slackChannelRepository,
	}
}

func WhenGettingMessagesByThreadIdAndChannelIdAfterEach(suite *WhenGettingMessagesByThreadIdAndChannelIdTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingMessagesByThreadIdAndChannelId(t *testing.T) {
	t.Parallel()

	t.Run("Given slack message does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return empty conversation", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingMessagesByThreadIdAndChannelIdBeforeEach()
			defer WhenGettingMessagesByThreadIdAndChannelIdAfterEach(suite)

			threadId := "non-existing-slack-thread-id"
			channelId := "non-existing-slack-channel-id"

			messages, err := suite.sut.GetByThreadIdAndChannelId(threadId, channelId)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if len(messages) != 0 {
				t.Fatalf("Expected empty messages, got %v", messages)
			}
		})
	})

	t.Run("Given slack message exists for this thread id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return empty list", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingMessagesByThreadIdAndChannelIdBeforeEach()
			defer WhenGettingMessagesByThreadIdAndChannelIdAfterEach(suite)

			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			slackThread := model_builders.NewSlackThreadBuilder().
				WithChannelId(slackChannel.Id).
				Build()

			firstMessage := model_builders.NewSlackMessageBuilder().
				WithThreadId(slackThread.Id).
				WithChannelId("non-existing-slack-channel-id").
				Build()

			secondMessage := model_builders.NewSlackMessageBuilder().
				WithThreadId(slackThread.Id).
				WithChannelId("non-existing-slack-channel-id").
				Build()

			suite.slackChannelRepository.Save(slackChannel)
			suite.slackThreadRepository.Save(slackThread)
			suite.sut.SaveMany([]*models.SlackMessage{firstMessage, secondMessage})

			messages, err := suite.sut.GetByThreadIdAndChannelId(slackThread.Id, slackChannel.Id)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if len(messages) != 0 {
				t.Fatalf("Expected empty messages, got %v", messages)
			}
		})
	})

	t.Run("Given slack message exists for this channel id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return empty list", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingMessagesByThreadIdAndChannelIdBeforeEach()
			defer WhenGettingMessagesByThreadIdAndChannelIdAfterEach(suite)

			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			slackThread := model_builders.NewSlackThreadBuilder().
				WithChannelId(slackChannel.Id).
				Build()

			firstMessage := model_builders.NewSlackMessageBuilder().
				WithThreadId("existing-slack-thread-id").
				WithChannelId(slackChannel.Id).
				Build()

			secondMessage := model_builders.NewSlackMessageBuilder().
				WithThreadId("non-existing-slack-thread-id").
				WithChannelId(slackChannel.Id).
				Build()

			suite.slackChannelRepository.Save(slackChannel)
			suite.slackThreadRepository.Save(slackThread)
			suite.sut.SaveMany([]*models.SlackMessage{firstMessage, secondMessage})

			messages, err := suite.sut.GetByThreadIdAndChannelId(slackThread.Id, slackChannel.Id)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if len(messages) != 0 {
				t.Fatalf("Expected empty messages, got %v", messages)
			}
		})
	})

	t.Run("Given slack message exists for this thread id and channel id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return messages", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingMessagesByThreadIdAndChannelIdBeforeEach()
			defer WhenGettingMessagesByThreadIdAndChannelIdAfterEach(suite)

			slackChannel := model_builders.NewSlackChannelBuilder().Build()

			slackThread := model_builders.NewSlackThreadBuilder().
				WithChannelId(slackChannel.Id).
				Build()

			firstMessage := model_builders.NewSlackMessageBuilder().
				WithThreadId(slackThread.Id).
				WithChannelId(slackChannel.Id).
				Build()

			secondMessage := model_builders.NewSlackMessageBuilder().
				WithThreadId(slackThread.Id).
				WithChannelId(slackChannel.Id).
				Build()

			suite.slackChannelRepository.Save(slackChannel)
			suite.slackThreadRepository.Save(slackThread)
			suite.sut.SaveMany([]*models.SlackMessage{firstMessage, secondMessage})

			messages, err := suite.sut.GetByThreadIdAndChannelId(slackThread.Id, slackChannel.Id)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if len(messages) == 0 {
				t.Fatalf("Expected messages, got empty list")
			}
		})
	})
}
