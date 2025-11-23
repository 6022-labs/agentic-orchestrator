package event_handlers

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"go.uber.org/zap"
)

type ChannelArchiveHandlerInterface interface {
	Handle(request *requests.ChannelArchiveRequest) error
}

type ChannelArchiveHandler struct {
	logger                *zap.Logger
	slackThreadRepository repositories.SlackThreadRepositoryInterface
}

func NewChannelArchiveHandler(
	logger *zap.Logger,
	slackThreadRepository repositories.SlackThreadRepositoryInterface,
) *ChannelArchiveHandler {
	return &ChannelArchiveHandler{
		logger:                logger,
		slackThreadRepository: slackThreadRepository,
	}
}

func (c *ChannelArchiveHandler) Handle(request *requests.ChannelArchiveRequest) error {
	err := c.slackThreadRepository.ConcludeThreadsByChannelId(request.Event.Channel)
	if err != nil {
		return err
	}

	c.logger.Info("Channel archive event handled successfully",
		zap.String("channel_id", request.Event.Channel),
	)

	return nil
}
