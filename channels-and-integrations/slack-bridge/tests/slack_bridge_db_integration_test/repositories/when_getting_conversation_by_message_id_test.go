package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	db_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingConversationByMessageIdTestingSuite struct {
	dbName   string
	dbClient *slack_bridge_db.DatabaseClient

	slackThreadRepository  *repositories.SlackThreadRepository
	slackChannelRepository *repositories.SlackChannelRepository

	sut *repositories.SlackMessageRepository
}

func WhenGettingConversationByMessageIdBeforeEach() *WhenGettingConversationByMessageIdTestingSuite {
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

	return &WhenGettingConversationByMessageIdTestingSuite{
		sut: sut,

		dbName:   dbName,
		dbClient: dbClient,

		slackThreadRepository:  slackThreadRepository,
		slackChannelRepository: slackChannelRepository,
	}
}

func WhenGettingConversationByMessageIdAfterEach(suite *WhenGettingConversationByMessageIdTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingConversationBySlackMessageId(t *testing.T) {
	t.Parallel()

	t.Run("Given slack message does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return empty conversation", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationByMessageIdBeforeEach()
			defer WhenGettingConversationByMessageIdAfterEach(suite)

			slackMessageId := "non-existing-slack-message-id"

			conversation, err := suite.sut.GetConversationByMessageId(slackMessageId)

			assert.Nil(t, err)
			assert.Empty(t, conversation)
		})
	})

	t.Run("Given slack message exists and it's a parent message", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return entire conversation ordered by timestamp", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationByMessageIdBeforeEach()
			defer WhenGettingConversationByMessageIdAfterEach(suite)

			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			slackThread := model_builders.NewSlackThreadBuilder().
				WithChannelId(slackChannel.Id).
				Build()
			firstMessage := model_builders.NewSlackMessageBuilder().
				WithId("1").
				WithThreadId(slackThread.Id).
				WithChannelId(slackChannel.Id).
				Build()

			secondMessage := model_builders.NewSlackMessageBuilder().
				WithId("2").
				WithThreadId(slackThread.Id).
				WithChannelId(slackChannel.Id).
				Build()

			suite.slackChannelRepository.Save(slackChannel)
			suite.slackThreadRepository.Save(slackThread)
			suite.sut.Save(firstMessage)
			suite.sut.Save(secondMessage)

			// Assuming the parent message and its replies are already saved in the database

			conversation, err := suite.sut.GetConversationByMessageId(firstMessage.Id)

			assert.Nil(t, err)
			assert.NotEmpty(t, conversation)
			assert.Len(t, conversation, 2)
			assert.Equal(t, firstMessage.Id, conversation[0].Id)
			assert.Equal(t, secondMessage.Id, conversation[1].Id)
		})
	})

	t.Run("Given slack message exists and it's a child message", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return entire conversation ordered by timestamp", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationByMessageIdBeforeEach()
			defer WhenGettingConversationByMessageIdAfterEach(suite)

			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			slackThread := model_builders.NewSlackThreadBuilder().
				WithChannelId(slackChannel.Id).
				Build()
			firstMessage := model_builders.NewSlackMessageBuilder().
				WithId("1").
				WithThreadId(slackThread.Id).
				WithChannelId(slackChannel.Id).
				Build()

			secondMessage := model_builders.NewSlackMessageBuilder().
				WithId("2").
				WithThreadId(slackThread.Id).
				WithChannelId(slackChannel.Id).
				Build()

			suite.slackChannelRepository.Save(slackChannel)
			suite.slackThreadRepository.Save(slackThread)
			suite.sut.Save(firstMessage)
			suite.sut.Save(secondMessage)

			conversation, err := suite.sut.GetConversationByMessageId(secondMessage.Id)

			assert.Nil(t, err)
			assert.NotEmpty(t, conversation)
			assert.Len(t, conversation, 2)
			assert.Equal(t, firstMessage.Id, conversation[0].Id)
			assert.Equal(t, secondMessage.Id, conversation[1].Id)
		})
	})
}
