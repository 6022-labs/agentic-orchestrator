package command_handlers

import (
	"context"
	"time"

	"fmt"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	slack_bridge_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	traces_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/traces/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type RemoveAgentFromChannelInterface interface {
	Remove(ctx context.Context, request *requests.RemoveAgentFromChannelRequest) error
}

type RemoveAgentFromChannel struct {
	channelMetrics metrics_interfaces.ChannelMetricsInterface
	channelTracer  traces_interfaces.ChannelTracerInterface

	logger                      *zap.Logger
	slackAgentRepository        repositories.SlackAgentRepositoryInterface
	slackChannelSynchronizer    interfaces.SlackChannelSynchronizerInterface
	slackChannelRepository      repositories.SlackChannelRepositoryInterface
	slackConversationRepository repositories.SlackConversationRepositoryInterface
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface
}

func NewRemoveAgentFromChannel(
	channelMetrics metrics_interfaces.ChannelMetricsInterface,
	channelTracer traces_interfaces.ChannelTracerInterface,
	logger *zap.Logger,
	slackAgentRepository repositories.SlackAgentRepositoryInterface,
	slackChannelRepository repositories.SlackChannelRepositoryInterface,
	slackChannelSynchronizer interfaces.SlackChannelSynchronizerInterface,
	slackConversationRepository repositories.SlackConversationRepositoryInterface,
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface,
) *RemoveAgentFromChannel {
	return &RemoveAgentFromChannel{
		channelMetrics: channelMetrics,
		channelTracer:  channelTracer,

		logger:                      logger,
		slackAgentRepository:        slackAgentRepository,
		slackChannelRepository:      slackChannelRepository,
		slackChannelSynchronizer:    slackChannelSynchronizer,
		slackConversationRepository: slackConversationRepository,
		conversationBrokerRequester: conversationBrokerRequester,
	}
}

func (h *RemoveAgentFromChannel) Remove(
	ctx context.Context,
	request *requests.RemoveAgentFromChannelRequest,
) (err error) {
	ctx, span := h.channelTracer.StartCommandHandling(ctx, "agent.remove")
	start := time.Now()
	defer func() {
		outcome := metrics_interfaces.ChannelOutcomeHandled
		if err != nil {
			outcome = metrics_interfaces.ChannelOutcomeFailed
			span.Fail(err)
		}
		h.channelMetrics.RecordCommand(
			ctx,
			"agent.remove",
			metrics_interfaces.CommandTriggerFrom(ctx),
			outcome,
			time.Since(start),
		)
		span.End()
	}()

	slackChannel, err := h.slackChannelRepository.GetById(ctx, request.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to get channel with ID %s: %w", request.ChannelID, err)
	}
	if slackChannel != nil {
		channelContainAgent := lo.ContainsBy(slackChannel.Agents, func(agent models.SlackAgent) bool {
			return agent.CollectionAddress == request.CollectionAddress &&
				agent.CollectionTokenId == request.CollectionTokenId
		})

		if !channelContainAgent { // Do not attempt to remove the agent if it's not in the channel
			return nil
		}
	}

	slackChannel, err = h.slackChannelSynchronizer.Sync(ctx, request.ChannelID)
	if err != nil {
		return err
	}

	slackAgent, err := h.slackAgentRepository.GetByCollectionAddressAndCollectionTokenId(ctx,
		request.CollectionAddress,
		request.CollectionTokenId,
	)
	if err != nil {
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

	if _, err = h.slackChannelRepository.RemoveAgent(ctx, slackChannel, slackAgent); err != nil {
		return fmt.Errorf("failed to remove agent from channel: %w", err)
	}

	if slackAgent.BrokerParticipantId == nil {
		return fmt.Errorf("slack agent %s/%d has no broker participant id; SlackAgentSynchronizer must populate it",
			slackAgent.CollectionAddress, slackAgent.CollectionTokenId)
	}

	conversations, err := h.slackConversationRepository.ListByChannelId(ctx, request.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to list slack conversations for channel %s: %w", request.ChannelID, err)
	}
	for _, conversation := range conversations {
		if err := h.conversationBrokerRequester.DetachParticipant(ctx,
			conversation.BrokerConversationID,
			*slackAgent.BrokerParticipantId,
		); err != nil {
			return fmt.Errorf(
				"failed to detach agent from broker conversation %s: %w",
				conversation.BrokerConversationID,
				err,
			)
		}
	}

	return nil
}
