package event_handlers

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	slack_bridge_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"go.uber.org/zap"
)

type ChannelLeftHandlerInterface interface {
	Handle(request *requests.ChannelLeftRequest) error
}

type ChannelLeftHandler struct {
	logger                      *zap.Logger
	slackConversationRepository repositories.SlackConversationRepositoryInterface
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface
}

func NewChannelLeftHandler(
	logger *zap.Logger,
	slackConversationRepository repositories.SlackConversationRepositoryInterface,
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface,
) *ChannelLeftHandler {
	return &ChannelLeftHandler{
		logger:                      logger,
		slackConversationRepository: slackConversationRepository,
		conversationBrokerRequester: conversationBrokerRequester,
	}
}

func (c *ChannelLeftHandler) Handle(request *requests.ChannelLeftRequest) error {
	conversations, err := c.slackConversationRepository.ListByChannelId(request.Event.Channel)
	if err != nil {
		return fmt.Errorf("failed to list slack conversations for channel %s: %w", request.Event.Channel, err)
	}

	for _, conversation := range conversations {
		if err := c.conversationBrokerRequester.ForceConcludeConversation(conversation.BrokerConversationID); err != nil {
			return fmt.Errorf("failed to force-conclude broker conversation %s: %w", conversation.BrokerConversationID, err)
		}
	}

	c.logger.Info("Channel left event handled successfully",
		zap.String("channel_id", request.Event.Channel),
		zap.Int("conversationsConcluded", len(conversations)),
	)

	return nil
}
