package services

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_remote_agent/clients"
)

type RemoteAgentCardRequester struct {
	client clients.RemoteAgentClientInterface
}

func NewRemoteAgentCardRequester(client clients.RemoteAgentClientInterface) *RemoteAgentCardRequester {
	return &RemoteAgentCardRequester{client: client}
}

func (r *RemoteAgentCardRequester) GetCard(
	ctx context.Context,
	agentUrl string,
) (*models.RemoteAgentCard, []byte, error) {
	external, raw, err := r.client.GetWellKnownCard(ctx, agentUrl)
	if err != nil {
		return nil, nil, err
	}

	card, err := external.ToModel(raw)
	if err != nil {
		return nil, nil, err
	}
	return card, raw, nil
}
