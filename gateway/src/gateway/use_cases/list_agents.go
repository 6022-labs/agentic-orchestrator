package use_cases

import (
	"context"

	"sort"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type ListAgentsInterface interface {
	List(ctx context.Context, request *requests.ListAgentsRequest) (*responses.ListAgentsResponse, error)
}

type ListAgents struct {
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface
}

func NewListAgents(agentRegistryRequester interfaces.AgentRegistryRequesterInterface) *ListAgents {
	return &ListAgents{agentRegistryRequester: agentRegistryRequester}
}

func (s *ListAgents) List(
	ctx context.Context,
	request *requests.ListAgentsRequest,
) (*responses.ListAgentsResponse, error) {
	request.Sanitize()

	registryResponse, err := s.agentRegistryRequester.ListAgents(ctx, &interfaces.ListRegistryAgentsRequest{
		Page:            request.Page,
		SortBy:          request.SortBy.ToDatabaseField(),
		SortDirection:   request.SortDirection.ToDatabaseDirection(),
		ElementsPerPage: request.ElementsPerPage,
		IsMinted:        request.IsMinted,
	})
	if err != nil {
		return nil, err
	}

	agentsResponse := registryResponse.Agents

	// Order by collection address ascending and then by collection token ID ascending
	sort.Slice(agentsResponse, func(i, j int) bool {
		if agentsResponse[i].CollectionAddress != agentsResponse[j].CollectionAddress {
			return agentsResponse[i].CollectionAddress < agentsResponse[j].CollectionAddress
		}

		if agentsResponse[i].CollectionTokenId == nil {
			return true
		}

		if agentsResponse[j].CollectionTokenId == nil {
			return false
		}

		return *agentsResponse[i].CollectionTokenId < *agentsResponse[j].CollectionTokenId
	})

	return &responses.ListAgentsResponse{
		PaginateResponse: responses.PaginateResponse{
			Page:            registryResponse.Page,
			ElementsPerPage: registryResponse.ElementsPerPage,
			TotalElements:   registryResponse.TotalElements,
		},
		Agents: agentsResponse,
	}, nil
}
