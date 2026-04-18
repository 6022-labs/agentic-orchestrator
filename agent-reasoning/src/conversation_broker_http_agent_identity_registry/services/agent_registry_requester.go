package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_agent_identity_registry/clients"
	"github.com/samber/lo"
)

type AgentRegistryRequester struct {
	client clients.AgentIdentityRegistryClientInterface
}

func NewAgentRegistryRequester(client clients.AgentIdentityRegistryClientInterface) *AgentRegistryRequester {
	return &AgentRegistryRequester{client: client}
}

func (r *AgentRegistryRequester) GetAgent(collectionAddress string, collectionTokenId uint64) (*models.Participant, error) {
	external, err := r.client.GetAgent(collectionAddress, collectionTokenId)
	if err != nil {
		return nil, err
	}
	if external == nil {
		return nil, nil
	}

	var role enums.ParticipantRole
	if err := role.Set(external.Role); err != nil {
		return nil, fmt.Errorf("registry returned invalid role: %w", err)
	}

	return &models.Participant{
		Role:              &role,
		ID:                external.ID,
		CollectionTokenId: external.CollectionTokenId,
		EnsDomain:         lo.ToPtr(external.EnsDomain),
		CollectionAddress: lo.ToPtr(external.CollectionAddress),
	}, nil
}
