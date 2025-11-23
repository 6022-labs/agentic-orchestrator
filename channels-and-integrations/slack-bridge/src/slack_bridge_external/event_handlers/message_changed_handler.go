package event_handlers

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"go.uber.org/zap"
)

type MessageChangedHandlerInterface interface {
	Handle(request *requests.MessageChangedRequest) error
}

type MessageChangedHandler struct {
	logger              *zap.Logger
	slackThreadConclude interfaces.SlackThreadConcludeInterface
}

func NewMessageChangedHandler(
	logger *zap.Logger,
	slackThreadConclude interfaces.SlackThreadConcludeInterface,
) *MessageChangedHandler {
	return &MessageChangedHandler{
		logger:              logger,
		slackThreadConclude: slackThreadConclude,
	}
}

func (h *MessageChangedHandler) Handle(request *requests.MessageChangedRequest) error {
	messageTs := request.Event.PreviousMessage.TS

	h.logger.Debug(
		"Handling message changed event",
		zap.String("message_ts", messageTs),
	)

	if request.Event.Message.Hidden {
		return h.slackThreadConclude.ConcludeIfExists(messageTs, request.Event.Channel)
	}

	return nil
}
