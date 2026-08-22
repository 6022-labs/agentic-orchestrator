package use_cases

import (
	"context"

	"sort"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type ListAgentsInterface interface {
	List(ctx context.Context, request *requests.ListAgentsRequest) (*responses.ListAgentsResponse, error)
}

type ListAgents struct {
	agentRepository        repositories.AgentRepositoryInterface
	agentAddressRepository repositories.AgentAddressRepositoryInterface
}

func NewListAgents(
	agentRepository repositories.AgentRepositoryInterface,
	agentAddressRepository repositories.AgentAddressRepositoryInterface,
) *ListAgents {
	return &ListAgents{
		agentRepository:        agentRepository,
		agentAddressRepository: agentAddressRepository,
	}
}

func (s *ListAgents) List(
	ctx context.Context,
	request *requests.ListAgentsRequest,
) (*responses.ListAgentsResponse, error) {
	request.Sanitize()

	searchParams := repositories.AgentSearchParams{IsMinted: request.IsMinted}

	agents, err := s.agentRepository.Search(ctx, repositories.AgentSearchParamsWithPagination{
		AgentSearchParams: searchParams,
		Page:              request.Page,
		SortBy:            request.SortBy.ToDatabaseField(),
		SortDirection:     request.SortDirection.ToDatabaseDirection(),
		ElementsPerPage:   request.ElementsPerPage,
	})
	if err != nil {
		return nil, err
	}

	totalCount, err := s.agentRepository.SearchCount(ctx, searchParams)
	if err != nil {
		return nil, err
	}

	agentIds := lo.Map(agents, func(agent *models.Agent, _ int) uuid.UUID {
		return agent.ID
	})
	agentAddresses, err := s.agentAddressRepository.GetAllByAgentIds(ctx, agentIds)
	if err != nil {
		return nil, err
	}
	addressesByAgentId := lo.GroupBy(agentAddresses, func(agentAddress *models.AgentAddress) uuid.UUID {
		return agentAddress.AgentId
	})

	agentsResponse := make([]*responses.AgentResponse, len(agents))
	for i, agent := range agents {
		agentsResponse[i] = responses.AgentResponseFromModel(agent, addressesByAgentId[agent.ID])
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
