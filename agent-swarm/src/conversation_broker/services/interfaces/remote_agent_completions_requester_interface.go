package interfaces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
)

// RemoteAgentCompletionsRequesterInterface posts one prepared input to the
// endpoint advertised by the agent's card (responses preferred, chatCompletions
// as fallback). Returns *errors.PaymentRequiredError on HTTP 402.
type RemoteAgentCompletionsRequesterInterface interface {
	RequestCompletions(
		ctx context.Context,
		card *models.RemoteAgentCard,
		input *models.RemoteAgentInput,
	) ([]string, error)
}
