package event_handlers

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"go.uber.org/zap"
)

type ChannelLeftHandlerInterface interface {
	Handle(request *requests.ChannelLeftRequest) error
}

type ChannelLeftHandler struct {
	logger                *zap.Logger
	slackThreadRepository repositories.SlackThreadRepositoryInterface
}

func NewChannelLeftHandler(
	logger *zap.Logger,
	slackThreadRepository repositories.SlackThreadRepositoryInterface,
) *ChannelLeftHandler {
	return &ChannelLeftHandler{
		logger:                logger,
		slackThreadRepository: slackThreadRepository,
	}
}

func (c *ChannelLeftHandler) Handle(request *requests.ChannelLeftRequest) error {
	err := c.slackThreadRepository.ConcludeThreadsByChannelId(request.Event.Channel)
	if err != nil {
		return err
	}

	c.logger.Info("Channel left event handled successfully",
		zap.String("channel_id", request.Event.Channel),
	)

	return nil
}
