package command_handlers

import (
	"context"
	"fmt"
	"time"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	slack_bridge_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	traces_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/traces/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"github.com/samber/lo"
)

type AddAgentToChannelInterface interface {
	Add(ctx context.Context, request *requests.AddAgentToChannelRequest) error
}

type AddAgentToChannel struct {
	channelMetrics metrics_interfaces.ChannelMetricsInterface
	channelTracer  traces_interfaces.ChannelTracerInterface

	slackAgentSynchronizer      interfaces.SlackAgentSynchronizerInterface
	slackChannelSynchronizer    interfaces.SlackChannelSynchronizerInterface
	slackChannelRepository      repositories.SlackChannelRepositoryInterface
	slackConversationRepository repositories.SlackConversationRepositoryInterface
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface
}

func NewAddAgentToChannel(
	channelMetrics metrics_interfaces.ChannelMetricsInterface,
	channelTracer traces_interfaces.ChannelTracerInterface,
	slackAgentSynchronizer interfaces.SlackAgentSynchronizerInterface,
	slackChannelRepository repositories.SlackChannelRepositoryInterface,
	slackChannelSynchronizer interfaces.SlackChannelSynchronizerInterface,
	slackConversationRepository repositories.SlackConversationRepositoryInterface,
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface,
) *AddAgentToChannel {
	return &AddAgentToChannel{
		channelMetrics: channelMetrics,
		channelTracer:  channelTracer,

		slackAgentSynchronizer:      slackAgentSynchronizer,
		slackChannelRepository:      slackChannelRepository,
		slackChannelSynchronizer:    slackChannelSynchronizer,
		slackConversationRepository: slackConversationRepository,
		conversationBrokerRequester: conversationBrokerRequester,
	}
}

func (h *AddAgentToChannel) Add(ctx context.Context, request *requests.AddAgentToChannelRequest) (err error) {
	ctx, span := h.channelTracer.StartCommandHandling(ctx, "agent.add")
	start := time.Now()
	defer func() {
		outcome := metrics_interfaces.ChannelOutcomeHandled
		if err != nil {
			outcome = metrics_interfaces.ChannelOutcomeFailed
			span.Fail(err)
		}
		h.channelMetrics.RecordCommand(
			ctx,
			"agent.add",
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

		if channelContainAgent { // Do not add the agent if it's already in the channel
			return nil
		}
	}

	slackChannel, err = h.slackChannelSynchronizer.Sync(ctx, request.ChannelID)
	if err != nil {
		return err
	}

	slackAgent, err := h.slackAgentSynchronizer.Sync(ctx, request.CollectionAddress, request.CollectionTokenId)
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

	if _, err := h.slackChannelRepository.AddAgent(ctx, slackChannel, slackAgent); err != nil {
		return fmt.Errorf("failed to add agent to channel: %w", err)
	}

	conversations, err := h.slackConversationRepository.ListByChannelId(ctx, request.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to list slack conversations for channel %s: %w", request.ChannelID, err)
	}
	for _, conversation := range conversations {
		if err := h.conversationBrokerRequester.AttachParticipant(ctx,
			conversation.BrokerConversationID,
			*slackAgent.BrokerParticipantId,
		); err != nil {
			return fmt.Errorf(
				"failed to attach agent to broker conversation %s: %w",
				conversation.BrokerConversationID,
				err,
			)
		}
	}

	return nil
}
