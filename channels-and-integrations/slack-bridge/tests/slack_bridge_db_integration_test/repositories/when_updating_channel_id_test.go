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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenUpdatingChannelIdTestingSuite struct {
	dbName   string
	dbClient *slack_bridge_db.DatabaseClient

	slackThreadRepository  *repositories.SlackThreadRepository
	slackMessageRepository *repositories.SlackMessageRepository
	slackAgentRepository   *repositories.SlackAgentRepository

	sut *repositories.SlackChannelRepository
}

func WhenUpdatingChannelIdBeforeEach() *WhenUpdatingChannelIdTestingSuite {
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

	slackChannelRepository := repositories.NewSlackChannelRepository(dbClient)
	slackThreadRepository := repositories.NewSlackThreadRepository(processingSettings, dbClient)
	slackMessageRepository := repositories.NewSlackMessageRepository(dbClient)
	slackAgentRepository := repositories.NewSlackAgentRepository(processingSettings, dbClient)

	return &WhenUpdatingChannelIdTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,

		slackThreadRepository:  slackThreadRepository,
		slackMessageRepository: slackMessageRepository,
		slackAgentRepository:   slackAgentRepository,

		sut: slackChannelRepository,
	}
}

func WhenUpdatingChannelIdAfterEach(suite *WhenUpdatingChannelIdTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenUpdatingChannelId(t *testing.T) {
	t.Parallel()

	seedSlackChannelRelations := func(t *testing.T, suite *WhenUpdatingChannelIdTestingSuite) (string, string, string, string, uint64) {
		t.Helper()

		slackChannel := model_builders.NewSlackChannelBuilder().Build()
		slackChannel, err := suite.sut.Save(slackChannel)
		assert.NoError(t, err)

		slackThread := model_builders.NewSlackThreadBuilder().
			WithChannelId(slackChannel.Id).
			Build()
		slackThread, err = suite.slackThreadRepository.Save(slackThread)
		assert.NoError(t, err)

		slackMessage := model_builders.NewSlackMessageBuilder().
			WithThreadId(slackThread.Id).
			WithChannelId(slackChannel.Id).
			Build()
		slackMessage, err = suite.slackMessageRepository.Save(slackMessage)
		assert.NoError(t, err)

		slackAgent := model_builders.NewSlackAgentBuilder().Build()
		slackAgent, err = suite.slackAgentRepository.Save(slackAgent)
		assert.NoError(t, err)

		_, err = suite.sut.AddAgent(slackChannel, slackAgent)
		assert.NoError(t, err)

		return slackChannel.Id, slackThread.Id, slackMessage.Id, slackAgent.CollectionAddress, slackAgent.CollectionTokenId
	}

	t.Run("Given slack channel has stored threads", func(t *testing.T) {
		t.Parallel()

		t.Run("Should move related slack threads to the new channel id", func(t *testing.T) {
			t.Parallel()

			suite := WhenUpdatingChannelIdBeforeEach()
			defer WhenUpdatingChannelIdAfterEach(suite)

			channelId, threadId, _, _, _ := seedSlackChannelRelations(t, suite)
			newChannelId := uuid.NewString()

			err := suite.sut.UpdateId(channelId, newChannelId)
			assert.NoError(t, err)

			updatedChannel, err := suite.sut.GetById(newChannelId)
			assert.NoError(t, err)
			assert.Equal(t, newChannelId, updatedChannel.Id)

			updatedThread, err := suite.slackThreadRepository.GetByIdAndChannelId(threadId, newChannelId)
			assert.NoError(t, err)
			assert.Equal(t, threadId, updatedThread.Id)
			assert.Equal(t, newChannelId, updatedThread.ChannelId)

			_, err = suite.slackThreadRepository.GetByIdAndChannelId(threadId, channelId)
			assert.Error(t, err)
		})
	})

	t.Run("Given slack channel has stored messages", func(t *testing.T) {
		t.Parallel()

		t.Run("Should move related slack messages to the new channel id", func(t *testing.T) {
			t.Parallel()

			suite := WhenUpdatingChannelIdBeforeEach()
			defer WhenUpdatingChannelIdAfterEach(suite)

			channelId, threadId, messageId, _, _ := seedSlackChannelRelations(t, suite)
			newChannelId := uuid.NewString()

			err := suite.sut.UpdateId(channelId, newChannelId)
			assert.NoError(t, err)

			updatedMessages, err := suite.slackMessageRepository.GetByThreadIdAndChannelId(threadId, newChannelId)
			assert.NoError(t, err)
			assert.Len(t, updatedMessages, 1)
			assert.Equal(t, messageId, updatedMessages[0].Id)
			assert.Equal(t, newChannelId, updatedMessages[0].ChannelId)

			staleMessages, err := suite.slackMessageRepository.GetByThreadIdAndChannelId(threadId, channelId)
			assert.NoError(t, err)
			assert.Empty(t, staleMessages)
		})
	})

	t.Run("Given slack channel is paired with agents", func(t *testing.T) {
		t.Parallel()

		t.Run("Should update slack_channel_agents rows to the new channel id", func(t *testing.T) {
			t.Parallel()

			suite := WhenUpdatingChannelIdBeforeEach()
			defer WhenUpdatingChannelIdAfterEach(suite)

			channelId, _, _, agentAddress, agentTokenId := seedSlackChannelRelations(t, suite)
			newChannelId := uuid.NewString()

			err := suite.sut.UpdateId(channelId, newChannelId)
			assert.NoError(t, err)

			var channelAgents []struct {
				SlackChannelId              string
				SlackAgentCollectionAddress string
				SlackAgentCollectionTokenId uint64
			}

			err = suite.dbClient.GetDB().
				Table("slack_channel_agents").
				Where(
					"slack_agent_collection_address = ? AND slack_agent_collection_token_id = ?",
					agentAddress,
					agentTokenId,
				).
				Find(&channelAgents).Error
			assert.NoError(t, err)
			assert.Len(t, channelAgents, 1)
			assert.Equal(t, newChannelId, channelAgents[0].SlackChannelId)

			var legacyCount int64
			err = suite.dbClient.GetDB().
				Table("slack_channel_agents").
				Where("slack_channel_id = ?", channelId).
				Count(&legacyCount).Error
			assert.NoError(t, err)
			assert.Zero(t, legacyCount)
		})
	})
}
