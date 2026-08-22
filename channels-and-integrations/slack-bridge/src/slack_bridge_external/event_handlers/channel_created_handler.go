package event_handlers

import (
	"context"

	"strings"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"go.uber.org/zap"
)

type ChannelCreatedHandlerInterface interface {
	Handle(ctx context.Context, request *requests.ChannelCreatedRequest) error
}

type ChannelCreatedHandler struct {
	logger                   *zap.Logger
	slackChannel             *settings.SlackSettings
	slackChannelSynchronizer interfaces.SlackChannelSynchronizerInterface
}

func NewChannelCreatedHandler(
	logger *zap.Logger,
	slackChannel *settings.SlackSettings,
	slackChannelSynchronizer interfaces.SlackChannelSynchronizerInterface,
) *ChannelCreatedHandler {
	return &ChannelCreatedHandler{
		logger:                   logger,
		slackChannel:             slackChannel,
		slackChannelSynchronizer: slackChannelSynchronizer,
	}
}

func (c *ChannelCreatedHandler) Handle(ctx context.Context, request *requests.ChannelCreatedRequest) error {
	if !strings.HasPrefix(request.Event.Channel.Name, c.slackChannel.AcceptedChannelPrefix) {
		c.logger.Info("Channel name does not start with accepted prefix, skipping join",
			zap.String("channel_id", request.Event.Channel.Id),
			zap.String("channel_name", request.Event.Channel.Name),
			zap.String("accepted_prefix", c.slackChannel.AcceptedChannelPrefix),
		)

		return nil
	}

	_, err := c.slackChannelSynchronizer.Sync(ctx, request.Event.Channel.Id)
	if err != nil {
		return err
	}

	c.logger.Info("Channel created event handled successfully",
		zap.String("channel_id", request.Event.Channel.Id),
		zap.String("channel_name", request.Event.Channel.Name),
	)

	return nil
}
