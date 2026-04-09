package use_cases

import (
	"sort"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
)

type ListAgentsInterface interface {
	List(request *requests.ListAgentsRequest) (*responses.ListAgentsResponse, error)
}

type ListAgents struct {
	agentRepository repositories.AgentRepositoryInterface
}

func NewListAgents(agentRepository repositories.AgentRepositoryInterface) *ListAgents {
	return &ListAgents{agentRepository: agentRepository}
}

func (s *ListAgents) List(request *requests.ListAgentsRequest) (*responses.ListAgentsResponse, error) {
	request.Sanitize()

	agents, err := s.agentRepository.Search(repositories.AgentSearchParamsWithPagination{
		AgentSearchParams: repositories.AgentSearchParams{IsInternal: request.IsInternal},
		Page:              request.Page,
		SortBy:            request.SortBy.ToDatabaseField(),
		SortDirection:     request.SortDirection.ToDatabaseDirection(),
		ElementsPerPage:   request.ElementsPerPage,
	})
	if err != nil {
		return nil, err
	}

	totalCount, err := s.agentRepository.SearchCount(repositories.AgentSearchParams{IsInternal: request.IsInternal})
	if err != nil {
		return nil, err
	}

	agentsResponse := make([]*responses.AgentResponse, len(agents))
	for i, agent := range agents {
		agentsResponse[i] = responses.AgentResponseFromModel(agent)
	}

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
			Page:            *request.Page,
			ElementsPerPage: *request.ElementsPerPage,
			TotalElements:   totalCount,
		},
		Agents: agentsResponse,
	}, nil
}
