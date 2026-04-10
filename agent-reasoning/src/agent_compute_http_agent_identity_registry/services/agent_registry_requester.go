package services

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_agent_identity_registry/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_agent_identity_registry/models"
)

type AgentRegistryRequester struct {
	client clients.AgentIdentityRegistryClientInterface
}

func NewAgentRegistryRequester(client clients.AgentIdentityRegistryClientInterface) *AgentRegistryRequester {
	return &AgentRegistryRequester{client: client}
}

func (r *AgentRegistryRequester) GetAgent(agentCollectionAddress string, agentCollectionTokenId uint64) (*responses.AgentResponse, error) {
	extResponse, err := r.client.GetAgent(agentCollectionAddress, agentCollectionTokenId)
	if err != nil {
		return nil, err
	}
	if extResponse == nil {
		return nil, nil
	}

	return mapExternalAgentToResponse(extResponse), nil
}

func (r *AgentRegistryRequester) GetMintedAgentByEnsDomain(ensDomain string) (*responses.AgentResponse, error) {
	extResponse, err := r.client.GetMintedAgentByEnsDomain(ensDomain)
	if err != nil {
		return nil, err
	}
	if extResponse == nil {
		return nil, nil
	}

	return mapExternalAgentToResponse(extResponse), nil
}

func (r *AgentRegistryRequester) GetAgentsByEnsDomains(ensDomains []string) ([]*responses.AgentResponse, error) {
	result := make([]*responses.AgentResponse, 0, len(ensDomains))
	for _, ensDomain := range ensDomains {
		agent, err := r.GetMintedAgentByEnsDomain(ensDomain)
		if err != nil {
			return nil, err
		}
		if agent != nil {
			result = append(result, agent)
		}
	}
	return result, nil
}

func mapExternalAgentToResponse(ext *models.ExternalAgentResponse) *responses.AgentResponse {
	return &responses.AgentResponse{
		ID:                ext.ID,
		Name:              ext.Name,
		Role:              ext.Role,
		Owner:             ext.Owner,
		Creator:           ext.Creator,
		Images:            ext.Images,
		CloneOf:           ext.CloneOf,
		Attributes:        ext.Attributes,
		EnsDomain:         ext.EnsDomain,
		IsInternal:        ext.IsInternal,
		WalletAddress:     ext.WalletAddress,
		ProposalId:        ext.ProposalId,
		CollectionTokenId: ext.CollectionTokenId,
		CollectionAddress: ext.CollectionAddress,
	}
}
