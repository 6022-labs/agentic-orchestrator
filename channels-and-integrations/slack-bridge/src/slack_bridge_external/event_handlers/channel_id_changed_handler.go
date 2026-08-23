package event_handlers

import (
	"context"
	"time"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	traces_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/traces/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"go.uber.org/zap"
)

type ChannelIdChangedHandlerInterface interface {
	Handle(ctx context.Context, request *requests.ChannelIdChangedRequest) error
}

type ChannelIdChangedHandler struct {
	channelMetrics metrics_interfaces.ChannelMetricsInterface
	channelTracer  traces_interfaces.ChannelTracerInterface

	logger                 *zap.Logger
	slackChannelRepository repositories.SlackChannelRepositoryInterface
}

func NewChannelIdChangedHandler(
	channelMetrics metrics_interfaces.ChannelMetricsInterface,
	channelTracer traces_interfaces.ChannelTracerInterface,
	logger *zap.Logger,
	slackChannelRepository repositories.SlackChannelRepositoryInterface,
) *ChannelIdChangedHandler {
	return &ChannelIdChangedHandler{
		channelMetrics: channelMetrics,
		channelTracer:  channelTracer,

		logger:                 logger,
		slackChannelRepository: slackChannelRepository,
	}
}

func (c *ChannelIdChangedHandler) Handle(ctx context.Context, request *requests.ChannelIdChangedRequest) (err error) {
	ctx, span := c.channelTracer.StartEventHandling(ctx, "channel.id_changed")
	start := time.Now()
	defer func() {
		outcome := metrics_interfaces.ChannelOutcomeHandled
		if err != nil {
			outcome = metrics_interfaces.ChannelOutcomeFailed
			span.Fail(err)
		}
		c.channelMetrics.RecordEvent(ctx, "channel.id_changed", outcome, time.Since(start))
		span.End()
	}()

	if err = c.slackChannelRepository.UpdateId(
		ctx,
		request.Event.OldChannelId,
		request.Event.NewChannelId,
	); err != nil {
		return err
	}

	c.logger.Info("Channel id changed event handled successfully",
		zap.String("old_channel_id", request.Event.OldChannelId),
		zap.String("new_channel_id", request.Event.NewChannelId),
	)

	return nil
}
