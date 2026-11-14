package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
)

type WhenAddingAgentToChannelTestingSuite struct {
	sut *repositories.SlackChannelRepository

	dbName   string
	dbClient *slack_bridge_db.DatabaseClient

	slackAgentRepository *repositories.SlackAgentRepository
}

func WhenAddingAgentToChannelBeforeEach(t *testing.T) *WhenAddingAgentToChannelTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)

	sut := repositories.NewSlackChannelRepository(dbClient)

	slackAgentRepository := repositories.NewSlackAgentRepository(dbClient)

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

			suite := WhenAddingAgentToChannelBeforeEach(t)
			defer WhenAddingAgentToChannelAfterEach(suite)

			channel := model_builders.NewSlackChannelBuilder().Build()
			agent := model_builders.NewSlackAgentBuilder().Build()

			_, err := suite.sut.Save(context.Background(), channel)
			if err != nil {
				t.Fatalf("Failed to save channel: %v", err)
			}

			_, err = suite.slackAgentRepository.Save(context.Background(), agent)
			if err != nil {
				t.Fatalf("Failed to save agent: %v", err)
			}

			_, err = suite.sut.AddAgent(context.Background(), channel, agent)
			assert.NoError(t, err)

			channel, err = suite.sut.GetById(context.Background(), channel.Id)

			assert.NoError(t, err)

			assert.Equal(t, 1, len(channel.Agents))
		})
	})

	t.Run("Given agent is already in channel", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not add agent to channel again", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAgentToChannelBeforeEach(t)
			defer WhenAddingAgentToChannelAfterEach(suite)

			channel := model_builders.NewSlackChannelBuilder().Build()
			agent := model_builders.NewSlackAgentBuilder().Build()

			_, err := suite.sut.Save(context.Background(), channel)
			if err != nil {
				t.Fatalf("Failed to save channel: %v", err)
			}

			_, err = suite.slackAgentRepository.Save(context.Background(), agent)
			if err != nil {
				t.Fatalf("Failed to save agent: %v", err)
			}

			_, err = suite.sut.AddAgent(context.Background(), channel, agent)
			assert.NoError(t, err)

			_, err = suite.sut.AddAgent(context.Background(), channel, agent)
			assert.NoError(t, err)

			channel, err = suite.sut.GetById(context.Background(), channel.Id)
			if err != nil {
				t.Fatalf("Failed to get channel: %v", err)
			}

			assert.Equal(t, 1, len(channel.Agents), "Expected only one agent in channel")
		})
	})
}
