package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	db_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingSlackThreadByIdAndChannelIdTestingSuite struct {
	dbName   string
	dbClient *slack_bridge_db.DatabaseClient

	slackChannelRepository *repositories.SlackChannelRepository

	sut *repositories.SlackThreadRepository
}

func WhenGettingSlackThreadByIdAndChannelIdBeforeEach() *WhenGettingSlackThreadByIdAndChannelIdTestingSuite {
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

	sut := repositories.NewSlackThreadRepository(dbClient)
	slackChannelRepository := repositories.NewSlackChannelRepository(dbClient)

	return &WhenGettingSlackThreadByIdAndChannelIdTestingSuite{
		sut: sut,

		dbName:   dbName,
		dbClient: dbClient,

		slackChannelRepository: slackChannelRepository,
	}
}

func WhenGettingSlackThreadByIdAndChannelIdAfterEach(suite *WhenGettingSlackThreadByIdAndChannelIdTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingSlackThreadByIdAndChannelId(t *testing.T) {
	t.Parallel()

	t.Run("Given slack thread does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackThreadByIdAndChannelIdBeforeEach()
			defer WhenGettingSlackThreadByIdAndChannelIdAfterEach(suite)

			id := "non-existing-slack-thread-id"
			channelId := "non-existing-slack-channel-id"

			thread, err := suite.sut.GetByIdAndChannelId(id, channelId)

			assert.NotNil(t, err)
			assert.Nil(t, thread)
		})
	})

	t.Run("Given slack thread exists for this id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackThreadByIdAndChannelIdBeforeEach()
			defer WhenGettingSlackThreadByIdAndChannelIdAfterEach(suite)

			slackThread := model_builders.NewSlackThreadBuilder().Build()

			suite.sut.Save(slackThread)

			thread, err := suite.sut.GetByIdAndChannelId(slackThread.Id, "non-existing-slack-channel-id")

			assert.Error(t, err)
			assert.Nil(t, thread)
		})
	})

	t.Run("Given slack thread exists for this channel id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackThreadByIdAndChannelIdBeforeEach()
			defer WhenGettingSlackThreadByIdAndChannelIdAfterEach(suite)

			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			slackThread := model_builders.NewSlackThreadBuilder().
				WithChannelId(slackChannel.Id).
				Build()

			suite.sut.Save(slackThread)

			thread, err := suite.sut.GetByIdAndChannelId("non-existing-slack-thread-id", slackThread.ChannelId)

			assert.Error(t, err)
			assert.Nil(t, thread)
		})
	})

	t.Run("Given slack thread exists for this id and channel id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return slack thread", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackThreadByIdAndChannelIdBeforeEach()
			defer WhenGettingSlackThreadByIdAndChannelIdAfterEach(suite)

			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			slackThread := model_builders.NewSlackThreadBuilder().
				WithChannelId(slackChannel.Id).
				Build()

			suite.slackChannelRepository.Save(slackChannel)
			suite.sut.Save(slackThread)

			thread, err := suite.sut.GetByIdAndChannelId(slackThread.Id, slackThread.ChannelId)

			assert.Nil(t, err)
			assert.Equal(t, slackThread, thread)
		})
	})
}
