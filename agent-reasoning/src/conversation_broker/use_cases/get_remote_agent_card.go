package use_cases

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services/interfaces"
)

type GetRemoteAgentCardInterface interface {
	Get(ensDomain string) (*models.RemoteAgentCard, error)
}

type GetRemoteAgentCard struct {
	remoteAgentCardResolver interfaces.RemoteAgentCardResolverInterface
}

func NewGetRemoteAgentCard(remoteAgentCardResolver interfaces.RemoteAgentCardResolverInterface) *GetRemoteAgentCard {
	return &GetRemoteAgentCard{remoteAgentCardResolver: remoteAgentCardResolver}
}

func (uc *GetRemoteAgentCard) Get(ensDomain string) (*models.RemoteAgentCard, error) {
	ensDomain = strings.TrimSpace(ensDomain)
	if ensDomain == "" {
		return nil, errors.NewValidationError("ensDomain", "must not be empty")
	}

	return uc.remoteAgentCardResolver.ResolveCard(ensDomain)
}
