package services

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_channels_and_integrations_bridge/clients"
	http_callback_models "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_channels_and_integrations_bridge/models"
)

type BridgeCallbackRequester struct {
	client clients.BridgeCallbackClientInterface
}

func NewBridgeCallbackRequester(client clients.BridgeCallbackClientInterface) *BridgeCallbackRequester {
	return &BridgeCallbackRequester{client: client}
}

func (r *BridgeCallbackRequester) PostCompletion(
	bridge *models.Bridge,
	conversation *models.Conversation,
	agent *models.Participant,
	content string,
) error {
	body := &http_callback_models.ExternalCompletionCallbackRequest{
		ConversationId: conversation.ID,
		ParticipantId:  agent.ID,
		Content:        content,
	}
	return r.client.PostCompletion(bridge.CallbackUrl, body)
}
