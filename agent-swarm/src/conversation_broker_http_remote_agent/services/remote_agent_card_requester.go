package services

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_remote_agent/clients"
)

type RemoteAgentCardRequester struct {
	client clients.RemoteAgentClientInterface
}

func NewRemoteAgentCardRequester(client clients.RemoteAgentClientInterface) *RemoteAgentCardRequester {
	return &RemoteAgentCardRequester{client: client}
}

func (r *RemoteAgentCardRequester) GetCard(agentUrl string) (*models.RemoteAgentCard, error) {
	external, err := r.client.GetWellKnownCard(agentUrl)
	if err != nil {
		return nil, err
	}

	return external.ToModel(), nil
}
