package command_handlers

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/responses"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
)

type ListAgentInChannelInterface interface {
	List(request *requests.ListAgentInChannelRequest) ([]*responses.ListAgentsInChannelResponse, error)
}

type ListAgentInChannel struct {
	slackChannelSynchronizer interfaces.SlackChannelSynchronizerInterface
}

func NewListAgentInChannel(
	slackChannelSynchronizer interfaces.SlackChannelSynchronizerInterface,
) *ListAgentInChannel {
	return &ListAgentInChannel{
		slackChannelSynchronizer: slackChannelSynchronizer,
	}
}

func (handler *ListAgentInChannel) List(request *requests.ListAgentInChannelRequest) ([]*responses.ListAgentsInChannelResponse, error) {
	slackChannel, err := handler.slackChannelSynchronizer.Sync(request.ChannelID)
	if err != nil {
		return nil, err
	}

	result := make([]*responses.ListAgentsInChannelResponse, len(slackChannel.Agents))
	for i, agent := range slackChannel.Agents {
		result[i] = &responses.ListAgentsInChannelResponse{
			Role:              agent.Role,
			EnsDomain:         agent.EnsDomain,
			CollectionAddress: agent.CollectionAddress,
			CollectionTokenId: agent.CollectionTokenId,
		}
	}

	return result, nil
}
