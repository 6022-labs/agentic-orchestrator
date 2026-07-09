package services

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/models"
)

type AgentRegistryRequester struct {
	client clients.AgentIdentityRegistryClientInterface
}

func NewAgentRegistryRequester(client clients.AgentIdentityRegistryClientInterface) *AgentRegistryRequester {
	return &AgentRegistryRequester{client: client}
}

func (r *AgentRegistryRequester) ListAgents(request *interfaces.ListRegistryAgentsRequest) (*responses.ListAgentsResponse, error) {
	extResponse, err := r.client.ListAgents(&models.ExternalListAgentsRequest{
		Page:            request.Page,
		ElementsPerPage: request.ElementsPerPage,
		SortBy:          request.SortBy,
		SortDirection:   request.SortDirection,
	})
	if err != nil {
		return nil, err
	}

	agents := make([]*responses.AgentResponse, len(extResponse.Agents))
	for i, a := range extResponse.Agents {
		agents[i] = &responses.AgentResponse{
			ID:                a.ID,
			Name:              a.Name,
			Role:              a.Role,
			Owner:             a.Owner,
			Creator:           a.Creator,
			Images:            a.Images,
			CloneOf:           a.CloneOf,
			Attributes:        a.Attributes,
			EnsDomain:         a.EnsDomain,
			WalletAddress:     a.WalletAddress,
			ProposalId:        a.ProposalId,
			CollectionTokenId: a.CollectionTokenId,
			CollectionAddress: a.CollectionAddress,
		}
	}

	return &responses.ListAgentsResponse{
		PaginateResponse: responses.PaginateResponse{
			Page:            extResponse.Page,
			ElementsPerPage: extResponse.ElementsPerPage,
			TotalElements:   extResponse.TotalElements,
		},
		Agents: agents,
	}, nil
}

func (r *AgentRegistryRequester) GetAgent(agentCollectionAddress string, agentCollectionTokenId uint64) (*responses.AgentResponse, error) {
	extResponse, err := r.client.GetAgent(agentCollectionAddress, agentCollectionTokenId)
	if err != nil {
		return nil, err
	}
	if extResponse == nil {
		return nil, nil
	}

	return &responses.AgentResponse{
		ID:                extResponse.ID,
		Name:              extResponse.Name,
		Role:              extResponse.Role,
		Owner:             extResponse.Owner,
		Creator:           extResponse.Creator,
		Images:            extResponse.Images,
		CloneOf:           extResponse.CloneOf,
		Attributes:        extResponse.Attributes,
		EnsDomain:         extResponse.EnsDomain,
		WalletAddress:     extResponse.WalletAddress,
		ProposalId:        extResponse.ProposalId,
		CollectionTokenId: extResponse.CollectionTokenId,
		CollectionAddress: extResponse.CollectionAddress,
	}, nil
}

func (r *AgentRegistryRequester) ListAgentCollections() ([]*responses.AgentCollectionResponse, error) {
	extResponse, err := r.client.ListAgentCollections()
	if err != nil {
		return nil, err
	}

	out := make([]*responses.AgentCollectionResponse, len(extResponse))
	for i, c := range extResponse {
		out[i] = &responses.AgentCollectionResponse{
			Address: c.Address,
			Name:    c.Name,
			Symbol:  c.Symbol,
		}
	}

	return out, nil
}
