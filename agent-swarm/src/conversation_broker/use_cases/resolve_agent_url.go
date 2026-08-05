package use_cases

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/responses"
)

type ResolveAgentUrlInterface interface {
	Resolve(request *requests.ResolveAgentUrlRequest) (*responses.AgentUrlResponse, error)
}

type ResolveAgentUrl struct {
	agentUrlResolver interfaces.AgentUrlResolverInterface
}

func NewResolveAgentUrl(
	agentUrlResolver interfaces.AgentUrlResolverInterface,
) *ResolveAgentUrl {
	return &ResolveAgentUrl{
		agentUrlResolver: agentUrlResolver,
	}
}

func (uc *ResolveAgentUrl) Resolve(request *requests.ResolveAgentUrlRequest) (*responses.AgentUrlResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	url, err := uc.agentUrlResolver.ResolveUrl(strings.TrimSpace(request.EnsDomain))
	if err != nil {
		return nil, err
	}
	return &responses.AgentUrlResponse{Url: url}, nil
}
