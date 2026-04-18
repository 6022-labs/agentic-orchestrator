package command_handlers

import (
	"errors"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	slack_bridge_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RemoveAgentFromChannelInterface interface {
	Remove(request *requests.RemoveAgentFromChannelRequest) error
}

type RemoveAgentFromChannel struct {
	logger                      *zap.Logger
	slackAgentRepository        repositories.SlackAgentRepositoryInterface
	slackChannelSynchronizer    interfaces.SlackChannelSynchronizerInterface
	slackChannelRepository      repositories.SlackChannelRepositoryInterface
	slackConversationRepository repositories.SlackConversationRepositoryInterface
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface
}

func NewRemoveAgentFromChannel(
	logger *zap.Logger,
	slackAgentRepository repositories.SlackAgentRepositoryInterface,
	slackChannelRepository repositories.SlackChannelRepositoryInterface,
	slackChannelSynchronizer interfaces.SlackChannelSynchronizerInterface,
	slackConversationRepository repositories.SlackConversationRepositoryInterface,
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface,
) *RemoveAgentFromChannel {
	return &RemoveAgentFromChannel{
		logger:                      logger,
		slackAgentRepository:        slackAgentRepository,
		slackChannelRepository:      slackChannelRepository,
		slackChannelSynchronizer:    slackChannelSynchronizer,
		slackConversationRepository: slackConversationRepository,
		conversationBrokerRequester: conversationBrokerRequester,
	}
}

func (h *RemoveAgentFromChannel) Remove(request *requests.RemoveAgentFromChannelRequest) error {
	slackChannel, err := h.slackChannelRepository.GetById(request.ChannelID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to get channel with ID %s: %w", request.ChannelID, err)
	}
	if slackChannel != nil {
		channelContainAgent := lo.ContainsBy(slackChannel.Agents, func(agent models.SlackAgent) bool {
			return agent.CollectionAddress == request.CollectionAddress && agent.CollectionTokenId == request.CollectionTokenId
		})

		if !channelContainAgent { // Do not attempt to remove the agent if it's not in the channel
			return nil
		}
	}

	slackChannel, err = h.slackChannelSynchronizer.Sync(request.ChannelID)
	if err != nil {
		return err
	}

	slackAgent, err := h.slackAgentRepository.GetByCollectionAddressAndCollectionTokenId(
		request.CollectionAddress,
		request.CollectionTokenId,
	)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to get agent with collection address %s and token ID %d: %w",
			request.CollectionAddress,
			request.CollectionTokenId,
			err,
		)
	}
	if slackAgent == nil {
		return nil
	}

	h.logger.Debug("removing agent from channel",
		zap.Any("agent", slackAgent),
		zap.String("channel_id", request.ChannelID),
	)

	if _, err = h.slackChannelRepository.RemoveAgent(slackChannel, slackAgent); err != nil {
		return fmt.Errorf("failed to remove agent from channel: %w", err)
	}

	if slackAgent.BrokerParticipantId == nil {
		return fmt.Errorf("slack agent %s/%d has no broker participant id; SlackAgentSynchronizer must populate it",
			slackAgent.CollectionAddress, slackAgent.CollectionTokenId)
	}

	conversations, err := h.slackConversationRepository.ListByChannelId(request.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to list slack conversations for channel %s: %w", request.ChannelID, err)
	}
	for _, conversation := range conversations {
		if err := h.conversationBrokerRequester.DetachParticipant(conversation.BrokerConversationID, *slackAgent.BrokerParticipantId); err != nil {
			return fmt.Errorf("failed to detach agent from broker conversation %s: %w", conversation.BrokerConversationID, err)
		}
	}

	return nil
}
