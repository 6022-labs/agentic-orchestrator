package use_cases

import (
	"context"

	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/responses"
)

type ResolveAgentUrlInterface interface {
	Resolve(ctx context.Context, request *requests.ResolveAgentUrlRequest) (*responses.AgentUrlResponse, error)
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

func (uc *ResolveAgentUrl) Resolve(
	ctx context.Context,
	request *requests.ResolveAgentUrlRequest,
) (*responses.AgentUrlResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	url, err := uc.agentUrlResolver.ResolveUrl(ctx, strings.TrimSpace(request.EnsDomain))
	if err != nil {
		return nil, errors.NewUnavailableError("agent_url_unresolved", "failed to resolve agent url", err)
	}
	return &responses.AgentUrlResponse{Url: url}, nil
}
