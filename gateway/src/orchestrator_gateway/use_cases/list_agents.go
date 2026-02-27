package use_cases

import (
	"sort"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/responses"
)

type ListAgentsInterface interface {
	List(request *requests.ListAgentsRequest) (*responses.ListAgentsResponse, error)
}

type ListAgents struct {
	agentRepository             repositories.AgentRepositoryInterface
	llmConnectionRequester      interfaces.LLMConnectionRequesterInterface
	databaseConnectionRequester interfaces.DatabaseConnectionRequesterInterface
}

func NewListAgents(
	agentRepository repositories.AgentRepositoryInterface,
	llmConnectionRequester interfaces.LLMConnectionRequesterInterface,
	databaseConnectionRequester interfaces.DatabaseConnectionRequesterInterface,
) *ListAgents {
	return &ListAgents{
		agentRepository:             agentRepository,
		llmConnectionRequester:      llmConnectionRequester,
		databaseConnectionRequester: databaseConnectionRequester,
	}
}

func (s *ListAgents) List(request *requests.ListAgentsRequest) (*responses.ListAgentsResponse, error) {
	request.Sanitize()

	agents, err := s.agentRepository.Search(repositories.AgentSearchParamsWithPagination{
		AgentSearchParams: repositories.AgentSearchParams{
			IsInternal: request.IsInternal,
		},
		Page:            request.Page,
		SortBy:          request.SortBy.ToDatabaseField(),
		SortDirection:   request.SortDirection.ToDatabaseDirection(),
		ElementsPerPage: request.ElementsPerPage,
	})
	if err != nil {
		return nil, err
	}

	totalCount, err := s.agentRepository.SearchCount(repositories.AgentSearchParams{
		IsInternal: request.IsInternal,
	})
	if err != nil {
		return nil, err
	}

	type setupResult struct {
		index                         int
		llmConnectionSetupStatus      *interfaces.LLMConnectionStatus
		databaseConnectionSetupStatus *interfaces.DatabaseConnectionStatus
		err                           error
	}
	resultsCh := make(chan setupResult, len(agents))

	for index, agent := range agents {
		go func(idx int, ag *models.Agent) {
			llmConnectionSetupStatus, databaseConnectionSetupStatus, err := getAgentSetupStatuses(
				ag,
				s.llmConnectionRequester,
				s.databaseConnectionRequester,
			)
			resultsCh <- setupResult{
				index:                         idx,
				llmConnectionSetupStatus:      llmConnectionSetupStatus,
				databaseConnectionSetupStatus: databaseConnectionSetupStatus,
				err:                           err,
			}
		}(index, agent)
	}

	agentsResponse := make([]*responses.AgentResponse, 0)
	errs := make([]error, 0)
	for i := 0; i < len(agents); i++ {
		res := <-resultsCh
		errs = append(errs, res.err)

		if res.err != nil {
			continue
		}

		agentsResponse = append(agentsResponse, responses.AgentResponseFromModel(
			agents[res.index],
			res.llmConnectionSetupStatus.Status,
			res.databaseConnectionSetupStatus.Status,
			res.llmConnectionSetupStatus.IsCompleted,
			res.databaseConnectionSetupStatus.IsCompleted,
		))
	}

	for _, e := range errs {
		if e != nil {
			return nil, e
		}
	}

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
			Page:            *request.Page,
			ElementsPerPage: *request.ElementsPerPage,
			TotalElements:   totalCount,
		},
		Agents: agentsResponse,
	}, nil
}
