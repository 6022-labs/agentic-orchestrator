package interfaces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
)

// RemoteAgentCompletionsRequesterInterface posts one prepared input to the endpoint the agent's card advertises.
// Returns *errors.PaymentRequiredError on HTTP 402, *errors.RemoteAgentUnavailableError on 429, 502, 503 and 504.
type RemoteAgentCompletionsRequesterInterface interface {
	RequestCompletions(
		ctx context.Context,
		card *models.RemoteAgentCard,
		input *models.RemoteAgentInput,
	) ([]string, error)
}
