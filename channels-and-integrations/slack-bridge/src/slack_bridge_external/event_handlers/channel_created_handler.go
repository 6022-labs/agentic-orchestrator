package event_handlers

import (
	"context"
	"time"

	"strings"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	traces_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/traces/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"go.uber.org/zap"
)

type ChannelCreatedHandlerInterface interface {
	Handle(ctx context.Context, request *requests.ChannelCreatedRequest) error
}

type ChannelCreatedHandler struct {
	channelMetrics metrics_interfaces.ChannelMetricsInterface
	channelTracer  traces_interfaces.ChannelTracerInterface

	logger                   *zap.Logger
	slackChannel             *settings.SlackSettings
	slackChannelSynchronizer interfaces.SlackChannelSynchronizerInterface
}

func NewChannelCreatedHandler(
	channelMetrics metrics_interfaces.ChannelMetricsInterface,
	channelTracer traces_interfaces.ChannelTracerInterface,
	logger *zap.Logger,
	slackChannel *settings.SlackSettings,
	slackChannelSynchronizer interfaces.SlackChannelSynchronizerInterface,
) *ChannelCreatedHandler {
	return &ChannelCreatedHandler{
		channelMetrics: channelMetrics,
		channelTracer:  channelTracer,

		logger:                   logger,
		slackChannel:             slackChannel,
		slackChannelSynchronizer: slackChannelSynchronizer,
	}
}

func (c *ChannelCreatedHandler) Handle(ctx context.Context, request *requests.ChannelCreatedRequest) (err error) {
	ctx, span := c.channelTracer.StartEventHandling(ctx, "channel.created")
	start := time.Now()
	defer func() {
		outcome := metrics_interfaces.ChannelOutcomeHandled
		if err != nil {
			outcome = metrics_interfaces.ChannelOutcomeFailed
			span.Fail(err)
		}
		c.channelMetrics.RecordEvent(ctx, "channel.created", outcome, time.Since(start))
		span.End()
	}()

	if !strings.HasPrefix(request.Event.Channel.Name, c.slackChannel.AcceptedChannelPrefix) {
		c.logger.Info("Channel name does not start with accepted prefix, skipping join",
			zap.String("channel_id", request.Event.Channel.Id),
			zap.String("channel_name", request.Event.Channel.Name),
			zap.String("accepted_prefix", c.slackChannel.AcceptedChannelPrefix),
		)

		return nil
	}

	if _, err = c.slackChannelSynchronizer.Sync(ctx, request.Event.Channel.Id); err != nil {
		return err
	}

	c.logger.Info("Channel created event handled successfully",
		zap.String("channel_id", request.Event.Channel.Id),
		zap.String("channel_name", request.Event.Channel.Name),
	)

	return nil
}
