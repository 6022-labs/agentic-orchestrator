package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	db_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	db_metrics_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_db_mocks/metrics_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenRemovingAgentFromChannelTestingSuite struct {
	sut *repositories.SlackChannelRepository

	dbName   string
	dbClient *slack_bridge_db.DatabaseClient

	slackAgentRepository *repositories.SlackAgentRepository
}

func WhenRemovingAgentFromChannelBeforeEach(t *testing.T) *WhenRemovingAgentFromChannelTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbQueryMetrics := db_metrics_mocks.NewMockDbQueryMetricsInterface(gomock.NewController(t))
	dbQueryMetrics.EXPECT().
		RecordQuery(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes()

	dbClient, _ := slack_bridge_db.NewDatabaseClient(
		&db_settings.DatabaseSettings{Url: connectionString},
		zap.NewNop(),
		&db_settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
		dbQueryMetrics,
	)
	err = dbClient.Migrate()
	if err != nil {
		panic(err)
	}

	sut := repositories.NewSlackChannelRepository(dbClient)

	slackAgentRepository := repositories.NewSlackAgentRepository(dbClient)

	return &WhenRemovingAgentFromChannelTestingSuite{
		sut: sut,

		dbName:   dbName,
		dbClient: dbClient,

		slackAgentRepository: slackAgentRepository,
	}
}

func WhenRemovingAgentFromChannelAfterEach(suite *WhenRemovingAgentFromChannelTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenRemovingAgentFromChannel(t *testing.T) {
	t.Parallel()

	t.Run("Given agent is not in channel", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not remove agents from channel", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromChannelBeforeEach(t)
			defer WhenRemovingAgentFromChannelAfterEach(suite)

			channel := model_builders.NewSlackChannelBuilder().Build()
			agentInChannel := model_builders.NewSlackAgentBuilder().Build()
			agentNotInChannel := model_builders.NewSlackAgentBuilder().Build()

			_, err := suite.sut.Save(context.Background(), channel)
			if err != nil {
				t.Fatalf("Unexpected error saving channel: %v", err)
			}

			_, err = suite.slackAgentRepository.Save(context.Background(), agentInChannel)
			if err != nil {
				t.Fatalf("Unexpected error saving agent: %v", err)
			}

			_, err = suite.slackAgentRepository.Save(context.Background(), agentNotInChannel)
			if err != nil {
				t.Fatalf("Unexpected error saving agent: %v", err)
			}

			_, err = suite.sut.AddAgent(context.Background(), channel, agentInChannel)
			if err != nil {
				t.Fatalf("Unexpected error adding agent to channel: %v", err)
			}

			_, err = suite.sut.RemoveAgent(context.Background(), channel, agentNotInChannel)

			assert.Nil(t, err)
			channel, err = suite.sut.GetById(context.Background(), channel.Id)
			if err != nil {
				t.Fatalf("Unexpected error getting channel: %v", err)
			}

			assert.Len(t, channel.Agents, 1, "Expected one agent in channel")
		})
	})

	t.Run("Given agent is in channel", func(t *testing.T) {
		t.Parallel()

		t.Run("Should remove agent from channel", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromChannelBeforeEach(t)
			defer WhenRemovingAgentFromChannelAfterEach(suite)

			channel := model_builders.NewSlackChannelBuilder().Build()
			agent1InChannel := model_builders.NewSlackAgentBuilder().Build()
			agent2InChannel := model_builders.NewSlackAgentBuilder().Build()

			_, err := suite.sut.Save(context.Background(), channel)
			if err != nil {
				t.Fatalf("Unexpected error saving channel: %v", err)
			}

			_, err = suite.slackAgentRepository.Save(context.Background(), agent1InChannel)
			if err != nil {
				t.Fatalf("Unexpected error saving agent: %v", err)
			}

			_, err = suite.slackAgentRepository.Save(context.Background(), agent2InChannel)
			if err != nil {
				t.Fatalf("Unexpected error saving agent: %v", err)
			}

			_, err = suite.sut.AddAgent(context.Background(), channel, agent1InChannel)
			if err != nil {
				t.Fatalf("Unexpected error adding agent to channel: %v", err)
			}

			_, err = suite.sut.AddAgent(context.Background(), channel, agent2InChannel)
			if err != nil {
				t.Fatalf("Unexpected error adding agent to channel: %v", err)
			}

			_, err = suite.sut.RemoveAgent(context.Background(), channel, agent1InChannel)

			assert.Nil(t, err, "Expected no error when removing agent from channel")

			channel, err = suite.sut.GetById(context.Background(), channel.Id)
			if err != nil {
				t.Fatalf("Unexpected error getting channel: %v", err)
			}

			assert.Len(t, channel.Agents, 1, "Expected one agent in channel after removal")
		})
	})
}
