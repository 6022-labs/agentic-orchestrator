package services

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_agent_identity_registry/clients"
)

type AgentRegistryRequester struct {
	client clients.AgentIdentityRegistryClientInterface
}

func NewAgentRegistryRequester(client clients.AgentIdentityRegistryClientInterface) *AgentRegistryRequester {
	return &AgentRegistryRequester{client: client}
}

func (r *AgentRegistryRequester) GetAgent(collectionAddress string, collectionTokenId uint64) (*models.RegistryAgent, error) {
	external, err := r.client.GetAgent(collectionAddress, collectionTokenId)
	if err != nil {
		return nil, err
	}

	return external.ToRegistryAgent(), nil
}

func (r *AgentRegistryRequester) GetMintedAgentByEnsDomain(ensDomain string) (*models.RegistryAgent, error) {
	external, err := r.client.GetMintedAgentByEnsDomain(ensDomain)
	if err != nil {
		return nil, err
	}

	return external.ToRegistryAgent(), nil
}
