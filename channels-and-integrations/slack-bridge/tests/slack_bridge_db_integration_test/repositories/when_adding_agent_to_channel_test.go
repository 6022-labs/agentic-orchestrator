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

type WhenAddingAgentToChannelTestingSuite struct {
	sut *repositories.SlackChannelRepository

	dbName   string
	dbClient *slack_bridge_db.DatabaseClient

	slackAgentRepository *repositories.SlackAgentRepository
}

func WhenAddingAgentToChannelBeforeEach() *WhenAddingAgentToChannelTestingSuite {
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

	sut := repositories.NewSlackChannelRepository(dbClient)

	slackAgentRepository := repositories.NewSlackAgentRepository(processingSettings, dbClient)

	return &WhenAddingAgentToChannelTestingSuite{
		sut: sut,

		dbName:   dbName,
		dbClient: dbClient,

		slackAgentRepository: slackAgentRepository,
	}
}

func WhenAddingAgentToChannelAfterEach(suite *WhenAddingAgentToChannelTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenAddingAgentToChannel(t *testing.T) {
	t.Parallel()

	t.Run("Given agent is not in channel", func(t *testing.T) {
		t.Parallel()

		t.Run("Should add agent to channel", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAgentToChannelBeforeEach()
			defer WhenAddingAgentToChannelAfterEach(suite)

			channel := model_builders.NewSlackChannelBuilder().Build()
			agent := model_builders.NewSlackAgentBuilder().Build()

			_, err := suite.sut.Save(channel)
			if err != nil {
				t.Fatalf("Failed to save channel: %v", err)
			}

			_, err = suite.slackAgentRepository.Save(agent)
			if err != nil {
				t.Fatalf("Failed to save agent: %v", err)
			}

			_, err = suite.sut.AddAgent(channel, agent)
			assert.NoError(t, err)

			channel, err = suite.sut.GetById(channel.Id)

			assert.NoError(t, err)

			assert.Equal(t, 1, len(channel.Agents))
		})
	})

	t.Run("Given agent is already in channel", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not add agent to channel again", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAgentToChannelBeforeEach()
			defer WhenAddingAgentToChannelAfterEach(suite)

			channel := model_builders.NewSlackChannelBuilder().Build()
			agent := model_builders.NewSlackAgentBuilder().Build()

			_, err := suite.sut.Save(channel)
			if err != nil {
				t.Fatalf("Failed to save channel: %v", err)
			}

			_, err = suite.slackAgentRepository.Save(agent)
			if err != nil {
				t.Fatalf("Failed to save agent: %v", err)
			}

			_, err = suite.sut.AddAgent(channel, agent)
			assert.NoError(t, err)

			_, err = suite.sut.AddAgent(channel, agent)
			assert.NoError(t, err)

			channel, err = suite.sut.GetById(channel.Id)
			if err != nil {
				t.Fatalf("Failed to get channel: %v", err)
			}

			assert.Equal(t, 1, len(channel.Agents), "Expected only one agent in channel")
		})
	})
}
