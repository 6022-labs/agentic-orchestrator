package command_handlers_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenListingAgentsInChannelTestingSuite struct {
	sut *command_handlers.ListAgentInChannel

	slackChannelSynchronizer interfaces_mocks.MockSlackChannelSynchronizerInterface
}

func WhenListingAgentsInChannelBeforeEach(t *testing.T) *WhenListingAgentsInChannelTestingSuite {
	mockController := gomock.NewController(t)

	slackChannelSynchronizer := interfaces_mocks.NewMockSlackChannelSynchronizerInterface(mockController)

	sut := command_handlers.NewListAgentInChannel(
		slackChannelSynchronizer,
	)

	return &WhenListingAgentsInChannelTestingSuite{
		sut: sut,

		slackChannelSynchronizer: *slackChannelSynchronizer,
	}
}

func TestWhenListingAgentsInChannel(t *testing.T) {
	t.Parallel()

	t.Run("should return error when failing to sync slack channel", func(t *testing.T) {
		t.Parallel()

		suite := WhenListingAgentsInChannelBeforeEach(t)

		request := &requests.ListAgentInChannelRequest{
			ChannelID: "channel-id",
		}

		suite.slackChannelSynchronizer.
			EXPECT().
			Sync(gomock.Any(), request.ChannelID).
			Return(nil, assert.AnError)

		agents, err := suite.sut.List(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, agents)
	})

	t.Run("should return agents in channel", func(t *testing.T) {
		t.Parallel()

		suite := WhenListingAgentsInChannelBeforeEach(t)

		request := &requests.ListAgentInChannelRequest{
			ChannelID: "channel-id",
		}

		slackChannel := model_builders.NewSlackChannelBuilder().
			WithId(request.ChannelID).
			WithAgents([]models.SlackAgent{
				(*model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleExpert).
					WithEnsDomain("expert.ens").
					WithCollectionAddress("0xexpert").
					WithCollectionTokenId(uint64(1)).
					Build()),
				(*model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleHuman).
					WithEnsDomain("human.ens").
					WithCollectionAddress("0xhuman").
					WithCollectionTokenId(uint64(2)).
					Build()),
			}).
			Build()

		suite.slackChannelSynchronizer.
			EXPECT().
			Sync(gomock.Any(), request.ChannelID).
			Return(slackChannel, nil)

		agents, err := suite.sut.List(context.Background(), request)

		assert.NoError(t, err)
		assert.Len(t, agents, 2)

		assert.Equal(t, enums.AgentRoleExpert, agents[0].Role)
		assert.Equal(t, "expert.ens", agents[0].EnsDomain)
		assert.Equal(t, "0xexpert", agents[0].CollectionAddress)
		assert.Equal(t, uint64(1), agents[0].CollectionTokenId)

		assert.Equal(t, enums.AgentRoleHuman, agents[1].Role)
		assert.Equal(t, "human.ens", agents[1].EnsDomain)
		assert.Equal(t, "0xhuman", agents[1].CollectionAddress)
		assert.Equal(t, uint64(2), agents[1].CollectionTokenId)
	})
}
