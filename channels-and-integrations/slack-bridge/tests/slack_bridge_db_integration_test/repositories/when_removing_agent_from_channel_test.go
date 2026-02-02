package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	db_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenRemovingAgentFromChannelTestingSuite struct {
	sut *repositories.SlackChannelRepository

	dbName   string
	dbClient *slack_bridge_db.DatabaseClient

	slackAgentRepository *repositories.SlackAgentRepository
}

func WhenRemovingAgentFromChannelBeforeEach() *WhenRemovingAgentFromChannelTestingSuite {
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

			suite := WhenRemovingAgentFromChannelBeforeEach()
			defer WhenRemovingAgentFromChannelAfterEach(suite)

			channel := model_builders.NewSlackChannelBuilder().Build()
			agentInChannel := model_builders.NewSlackAgentBuilder().Build()
			agentNotInChannel := model_builders.NewSlackAgentBuilder().Build()

			_, err := suite.sut.Save(channel)
			if err != nil {
				t.Fatalf("Unexpected error saving channel: %v", err)
			}

			_, err = suite.slackAgentRepository.Save(agentInChannel)
			if err != nil {
				t.Fatalf("Unexpected error saving agent: %v", err)
			}

			_, err = suite.slackAgentRepository.Save(agentNotInChannel)
			if err != nil {
				t.Fatalf("Unexpected error saving agent: %v", err)
			}

			_, err = suite.sut.AddAgent(channel, agentInChannel)
			if err != nil {
				t.Fatalf("Unexpected error adding agent to channel: %v", err)
			}

			_, err = suite.sut.RemoveAgent(channel, agentNotInChannel)

			assert.Nil(t, err)
			channel, err = suite.sut.GetById(channel.Id)
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

			suite := WhenRemovingAgentFromChannelBeforeEach()
			defer WhenRemovingAgentFromChannelAfterEach(suite)

			channel := model_builders.NewSlackChannelBuilder().Build()
			agent1InChannel := model_builders.NewSlackAgentBuilder().Build()
			agent2InChannel := model_builders.NewSlackAgentBuilder().Build()

			_, err := suite.sut.Save(channel)
			if err != nil {
				t.Fatalf("Unexpected error saving channel: %v", err)
			}

			_, err = suite.slackAgentRepository.Save(agent1InChannel)
			if err != nil {
				t.Fatalf("Unexpected error saving agent: %v", err)
			}

			_, err = suite.slackAgentRepository.Save(agent2InChannel)
			if err != nil {
				t.Fatalf("Unexpected error saving agent: %v", err)
			}

			_, err = suite.sut.AddAgent(channel, agent1InChannel)
			if err != nil {
				t.Fatalf("Unexpected error adding agent to channel: %v", err)
			}

			_, err = suite.sut.AddAgent(channel, agent2InChannel)
			if err != nil {
				t.Fatalf("Unexpected error adding agent to channel: %v", err)
			}

			_, err = suite.sut.RemoveAgent(channel, agent1InChannel)

			assert.Nil(t, err, "Expected no error when removing agent from channel")

			channel, err = suite.sut.GetById(channel.Id)
			if err != nil {
				t.Fatalf("Unexpected error getting channel: %v", err)
			}

			assert.Len(t, channel.Agents, 1, "Expected one agent in channel after removal")
		})
	})
}
