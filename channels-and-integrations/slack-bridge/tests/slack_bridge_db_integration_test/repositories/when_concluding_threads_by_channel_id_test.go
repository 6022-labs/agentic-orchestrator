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
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenConcludingThreadsByChannelIdTestingSuite struct {
	sut *repositories.SlackThreadRepository

	dbName   string
	dbClient *slack_bridge_db.DatabaseClient

	slackAgentRepository   *repositories.SlackAgentRepository
	slackChannelRepository *repositories.SlackChannelRepository
}

func WhenConcludingThreadsByChannelIdBeforeEach() *WhenConcludingThreadsByChannelIdTestingSuite {
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

	sut := repositories.NewSlackThreadRepository(processingSettings, dbClient)

	slackAgentRepository := repositories.NewSlackAgentRepository(processingSettings, dbClient)
	slackChannelRepository := repositories.NewSlackChannelRepository(dbClient)

	return &WhenConcludingThreadsByChannelIdTestingSuite{
		sut: sut,

		dbName:   dbName,
		dbClient: dbClient,

		slackAgentRepository:   slackAgentRepository,
		slackChannelRepository: slackChannelRepository,
	}
}

func WhenConcludingThreadsByChannelIdAfterEach(suite *WhenConcludingThreadsByChannelIdTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenConcludingThreadsByChannelId(t *testing.T) {
	t.Parallel()

	t.Run("Given multiple channels have multiple slack threads", func(t *testing.T) {
		t.Parallel()

		t.Run("Should mark only target channel threads as concluded", func(t *testing.T) {
			t.Parallel()

			suite := WhenConcludingThreadsByChannelIdBeforeEach()
			defer WhenConcludingThreadsByChannelIdAfterEach(suite)

			targetChannel := model_builders.NewSlackChannelBuilder().Build()
			otherChannel := model_builders.NewSlackChannelBuilder().Build()

			_, err := suite.slackChannelRepository.Save(targetChannel)
			assert.NoError(t, err)
			_, err = suite.slackChannelRepository.Save(otherChannel)
			assert.NoError(t, err)

			targetThreads := []*models.SlackThread{
				model_builders.NewSlackThreadBuilder().WithChannelId(targetChannel.Id).WithIsConcluded(false).Build(),
				model_builders.NewSlackThreadBuilder().WithChannelId(targetChannel.Id).WithIsConcluded(true).Build(),
			}
			otherThread := model_builders.NewSlackThreadBuilder().WithChannelId(otherChannel.Id).Build()

			for _, thread := range append(targetThreads, otherThread) {
				_, err = suite.sut.Save(thread)
				assert.NoError(t, err)
			}

			err = suite.sut.ConcludeThreadsByChannelId(targetChannel.Id)
			assert.NoError(t, err)

			for _, tc := range targetThreads {
				storedThread, err := suite.sut.GetByIdAndChannelId(tc.Id, tc.ChannelId)
				assert.NoError(t, err)
				assert.True(t, storedThread.IsConcluded)
			}

			storedOtherThread, err := suite.sut.GetByIdAndChannelId(otherThread.Id, otherThread.ChannelId)
			assert.NoError(t, err)
			assert.False(t, storedOtherThread.IsConcluded)
		})
	})

	t.Run("Given single channel has multiple slack threads", func(t *testing.T) {
		t.Parallel()

		t.Run("Should mark all threads as concluded", func(t *testing.T) {
			t.Parallel()

			suite := WhenConcludingThreadsByChannelIdBeforeEach()
			defer WhenConcludingThreadsByChannelIdAfterEach(suite)

			channel := model_builders.NewSlackChannelBuilder().Build()
			_, err := suite.slackChannelRepository.Save(channel)
			assert.NoError(t, err)

			threads := []*models.SlackThread{
				model_builders.NewSlackThreadBuilder().WithChannelId(channel.Id).WithIsConcluded(false).Build(),
				model_builders.NewSlackThreadBuilder().WithChannelId(channel.Id).WithIsConcluded(false).Build(),
			}

			for _, thread := range threads {
				_, err = suite.sut.Save(thread)
				assert.NoError(t, err)
			}

			err = suite.sut.ConcludeThreadsByChannelId(channel.Id)
			assert.NoError(t, err)

			for _, thread := range threads {
				storedThread, err := suite.sut.GetByIdAndChannelId(thread.Id, thread.ChannelId)
				assert.NoError(t, err)
				assert.True(t, storedThread.IsConcluded)
			}
		})
	})

	t.Run("Given channel has no slack threads", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil and leave other threads untouched", func(t *testing.T) {
			t.Parallel()

			suite := WhenConcludingThreadsByChannelIdBeforeEach()
			defer WhenConcludingThreadsByChannelIdAfterEach(suite)

			emptyChannel := model_builders.NewSlackChannelBuilder().Build()
			otherChannel := model_builders.NewSlackChannelBuilder().Build()
			_, err := suite.slackChannelRepository.Save(emptyChannel)
			assert.NoError(t, err)
			_, err = suite.slackChannelRepository.Save(otherChannel)
			assert.NoError(t, err)

			otherThread := model_builders.NewSlackThreadBuilder().WithChannelId(otherChannel.Id).Build()
			_, err = suite.sut.Save(otherThread)
			assert.NoError(t, err)

			err = suite.sut.ConcludeThreadsByChannelId(emptyChannel.Id)
			assert.NoError(t, err)

			storedOtherThread, err := suite.sut.GetByIdAndChannelId(otherThread.Id, otherThread.ChannelId)
			assert.NoError(t, err)
			assert.False(t, storedOtherThread.IsConcluded)
		})
	})

}
