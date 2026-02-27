package event_handlers

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"go.uber.org/zap"
)

type ChannelDeletedHandlerInterface interface {
	Handle(request *requests.ChannelDeletedRequest) error
}

type ChannelDeletedHandler struct {
	logger                *zap.Logger
	slackThreadRepository repositories.SlackThreadRepositoryInterface
}

func NewChannelDeletedHandler(
	logger *zap.Logger,
	slackThreadRepository repositories.SlackThreadRepositoryInterface,
) *ChannelDeletedHandler {
	return &ChannelDeletedHandler{
		logger:                logger,
		slackThreadRepository: slackThreadRepository,
	}
}

func (c *ChannelDeletedHandler) Handle(request *requests.ChannelDeletedRequest) error {
	err := c.slackThreadRepository.ConcludeThreadsByChannelId(request.Event.Channel)
	if err != nil {
		return err
	}

	c.logger.Info("Channel deleted event handled successfully",
		zap.String("channel_id", request.Event.Channel),
	)

	return nil
}
