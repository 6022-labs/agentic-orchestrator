package command_handlers

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	slack_bridge_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"github.com/samber/lo"
)

type AddAgentToChannelInterface interface {
	Add(request *requests.AddAgentToChannelRequest) error
}

type AddAgentToChannel struct {
	slackAgentSynchronizer      interfaces.SlackAgentSynchronizerInterface
	slackChannelSynchronizer    interfaces.SlackChannelSynchronizerInterface
	slackChannelRepository      repositories.SlackChannelRepositoryInterface
	slackConversationRepository repositories.SlackConversationRepositoryInterface
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface
}

func NewAddAgentToChannel(
	slackAgentSynchronizer interfaces.SlackAgentSynchronizerInterface,
	slackChannelRepository repositories.SlackChannelRepositoryInterface,
	slackChannelSynchronizer interfaces.SlackChannelSynchronizerInterface,
	slackConversationRepository repositories.SlackConversationRepositoryInterface,
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface,
) *AddAgentToChannel {
	return &AddAgentToChannel{
		slackAgentSynchronizer:      slackAgentSynchronizer,
		slackChannelRepository:      slackChannelRepository,
		slackChannelSynchronizer:    slackChannelSynchronizer,
		slackConversationRepository: slackConversationRepository,
		conversationBrokerRequester: conversationBrokerRequester,
	}
}

func (h *AddAgentToChannel) Add(request *requests.AddAgentToChannelRequest) error {
	slackChannel, err := h.slackChannelRepository.GetById(request.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to get channel with ID %s: %w", request.ChannelID, err)
	}
	if slackChannel != nil {
		channelContainAgent := lo.ContainsBy(slackChannel.Agents, func(agent models.SlackAgent) bool {
			return agent.CollectionAddress == request.CollectionAddress && agent.CollectionTokenId == request.CollectionTokenId
		})

		if channelContainAgent { // Do not add the agent if it's already in the channel
			return nil
		}
	}

	slackChannel, err = h.slackChannelSynchronizer.Sync(request.ChannelID)
	if err != nil {
		return err
	}

	slackAgent, err := h.slackAgentSynchronizer.Sync(request.CollectionAddress, request.CollectionTokenId)
	if err != nil {
		return fmt.Errorf("failed to sync agent with collection address %s and token ID %d: %w",
			request.CollectionAddress,
			request.CollectionTokenId,
			err,
		)
	}
	if slackAgent.BrokerParticipantId == nil {
		return fmt.Errorf("slack agent %s/%d has no broker participant id; SlackAgentSynchronizer must populate it",
			slackAgent.CollectionAddress, slackAgent.CollectionTokenId)
	}

	if _, err := h.slackChannelRepository.AddAgent(slackChannel, slackAgent); err != nil {
		return fmt.Errorf("failed to add agent to channel: %w", err)
	}

	conversations, err := h.slackConversationRepository.ListByChannelId(request.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to list slack conversations for channel %s: %w", request.ChannelID, err)
	}
	for _, conversation := range conversations {
		if err := h.conversationBrokerRequester.AttachParticipant(conversation.BrokerConversationID, *slackAgent.BrokerParticipantId); err != nil {
			return fmt.Errorf("failed to attach agent to broker conversation %s: %w", conversation.BrokerConversationID, err)
		}
	}

	return nil
}
