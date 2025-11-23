package event_handlers

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"go.uber.org/zap"
)

type MessageDeletedHandlerInterface interface {
	Handle(request *requests.MessageDeletedRequest) error
}

type MessageDeletedHandler struct {
	logger              *zap.Logger
	slackThreadConclude interfaces.SlackThreadConcludeInterface
}

func NewMessageDeletedHandler(
	logger *zap.Logger,
	slackThreadConclude interfaces.SlackThreadConcludeInterface,
) *MessageDeletedHandler {
	return &MessageDeletedHandler{
		logger:              logger,
		slackThreadConclude: slackThreadConclude,
	}
}

func (h *MessageDeletedHandler) Handle(request *requests.MessageDeletedRequest) error {
	messageTs := request.Event.PreviousMessage.TS

	h.logger.Debug(
		"Handling message deleted event",
		zap.String("message_ts", messageTs),
	)

	return h.slackThreadConclude.ConcludeIfExists(messageTs, request.Event.Channel)
}
