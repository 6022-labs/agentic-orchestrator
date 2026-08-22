package event_handlers

import (
	"context"

	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	slack_bridge_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"go.uber.org/zap"
)

type MessageChangedHandlerInterface interface {
	Handle(ctx context.Context, request *requests.MessageChangedRequest) error
}

type MessageChangedHandler struct {
	logger                      *zap.Logger
	slackConversationRepository repositories.SlackConversationRepositoryInterface
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface
}

func NewMessageChangedHandler(
	logger *zap.Logger,
	slackConversationRepository repositories.SlackConversationRepositoryInterface,
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface,
) *MessageChangedHandler {
	return &MessageChangedHandler{
		logger:                      logger,
		slackConversationRepository: slackConversationRepository,
		conversationBrokerRequester: conversationBrokerRequester,
	}
}

func (h *MessageChangedHandler) Handle(ctx context.Context, request *requests.MessageChangedRequest) error {
	messageTs := request.Event.PreviousMessage.TS

	h.logger.Debug("Handling message changed event", zap.String("message_ts", messageTs))

	if !request.Event.Message.Hidden {
		return nil
	}

	conversation, err := h.slackConversationRepository.GetBySlackIds(ctx, request.Event.Channel, messageTs)
	if err != nil {
		return fmt.Errorf("failed to lookup slack conversation: %w", err)
	}
	if conversation == nil {
		return nil
	}

	if err := h.conversationBrokerRequester.ForceConcludeConversation(
		ctx,
		conversation.BrokerConversationID,
	); err != nil {
		return fmt.Errorf("failed to force-conclude broker conversation %s: %w", conversation.BrokerConversationID, err)
	}

	return nil
}
