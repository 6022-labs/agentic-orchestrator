package event_handlers

import (
	"context"
	"time"

	"fmt"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	slack_bridge_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	traces_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/traces/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"go.uber.org/zap"
)

type ChannelDeletedHandlerInterface interface {
	Handle(ctx context.Context, request *requests.ChannelDeletedRequest) error
}

type ChannelDeletedHandler struct {
	channelMetrics metrics_interfaces.ChannelMetricsInterface
	channelTracer  traces_interfaces.ChannelTracerInterface

	logger                      *zap.Logger
	slackConversationRepository repositories.SlackConversationRepositoryInterface
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface
}

func NewChannelDeletedHandler(
	channelMetrics metrics_interfaces.ChannelMetricsInterface,
	channelTracer traces_interfaces.ChannelTracerInterface,
	logger *zap.Logger,
	slackConversationRepository repositories.SlackConversationRepositoryInterface,
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface,
) *ChannelDeletedHandler {
	return &ChannelDeletedHandler{
		channelMetrics: channelMetrics,
		channelTracer:  channelTracer,

		logger:                      logger,
		slackConversationRepository: slackConversationRepository,
		conversationBrokerRequester: conversationBrokerRequester,
	}
}

func (c *ChannelDeletedHandler) Handle(ctx context.Context, request *requests.ChannelDeletedRequest) (err error) {
	ctx, span := c.channelTracer.StartEventHandling(ctx, "channel.deleted")
	start := time.Now()
	defer func() {
		outcome := metrics_interfaces.ChannelOutcomeHandled
		if err != nil {
			outcome = metrics_interfaces.ChannelOutcomeFailed
			span.Fail(err)
		}
		c.channelMetrics.RecordEvent(ctx, "channel.deleted", outcome, time.Since(start))
		span.End()
	}()

	conversations, err := c.slackConversationRepository.ListByChannelId(ctx, request.Event.Channel)
	if err != nil {
		return fmt.Errorf("failed to list slack conversations for channel %s: %w", request.Event.Channel, err)
	}

	for _, conversation := range conversations {
		if err := c.conversationBrokerRequester.ForceConcludeConversation(ctx,
			conversation.BrokerConversationID,
		); err != nil {
			return fmt.Errorf(
				"failed to force-conclude broker conversation %s: %w",
				conversation.BrokerConversationID,
				err,
			)
		}
	}

	c.logger.Info("Channel deleted event handled successfully",
		zap.String("channel_id", request.Event.Channel),
		zap.Int("conversationsConcluded", len(conversations)),
	)

	return nil
}
