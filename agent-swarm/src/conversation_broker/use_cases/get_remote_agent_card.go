package use_cases

import (
	"context"

	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/responses"
)

type GetRemoteAgentCardInterface interface {
	Get(ctx context.Context, ensDomain string) (*responses.RemoteAgentCardResponse, error)
}

type GetRemoteAgentCard struct {
	remoteAgentCardResolver interfaces.RemoteAgentCardResolverInterface
}

func NewGetRemoteAgentCard(remoteAgentCardResolver interfaces.RemoteAgentCardResolverInterface) *GetRemoteAgentCard {
	return &GetRemoteAgentCard{remoteAgentCardResolver: remoteAgentCardResolver}
}

func (uc *GetRemoteAgentCard) Get(ctx context.Context, ensDomain string) (*responses.RemoteAgentCardResponse, error) {
	ensDomain = strings.TrimSpace(ensDomain)
	if ensDomain == "" {
		return nil, errors.NewValidationError("ensDomain", "must not be empty")
	}

	card, err := uc.remoteAgentCardResolver.ResolveCard(ctx, ensDomain)
	if err != nil {
		return nil, err
	}
	return responses.RemoteAgentCardResponseFromModel(card), nil
}
