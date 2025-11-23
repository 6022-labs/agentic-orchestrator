package event_handlers

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"go.uber.org/zap"
)

type ChannelIdChangedHandlerInterface interface {
	Handle(request *requests.ChannelIdChangedRequest) error
}

type ChannelIdChangedHandler struct {
	logger                 *zap.Logger
	slackChannelRepository repositories.SlackChannelRepositoryInterface
}

func NewChannelIdChangedHandler(
	logger *zap.Logger,
	slackChannelRepository repositories.SlackChannelRepositoryInterface,
) *ChannelIdChangedHandler {
	return &ChannelIdChangedHandler{
		logger:                 logger,
		slackChannelRepository: slackChannelRepository,
	}
}

func (c *ChannelIdChangedHandler) Handle(request *requests.ChannelIdChangedRequest) error {
	err := c.slackChannelRepository.UpdateId(request.Event.OldChannelId, request.Event.NewChannelId)
	if err != nil {
		return err
	}

	c.logger.Info("Channel id changed event handled successfully",
		zap.String("old_channel_id", request.Event.OldChannelId),
		zap.String("new_channel_id", request.Event.NewChannelId),
	)

	return nil
}
