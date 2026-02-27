package use_cases

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/compute"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/responses"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type ListAgentContextKeysInterface interface {
	List(request *requests.ListAgentContextKeysRequest) (*responses.ListAgentContextKeysResponse, error)
}

type ListAgentContextKeys struct {
	agentRepository       repositories.AgentRepositoryInterface
	agentContextRequester interfaces.AgentContextRequesterInterface
}

func NewListAgentContextKeys(
	agentRepository repositories.AgentRepositoryInterface,
	agentContextRequester interfaces.AgentContextRequesterInterface,
) *ListAgentContextKeys {
	return &ListAgentContextKeys{
		agentRepository:       agentRepository,
		agentContextRequester: agentContextRequester,
	}
}

func (c *ListAgentContextKeys) List(request *requests.ListAgentContextKeysRequest) (*responses.ListAgentContextKeysResponse, error) {
	err := request.EnsureAuthenticatedWallet()
	if err != nil {
		return nil, err
	}
	request.Sanitize()

	agent, err := c.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if agent == nil {
		return nil, errors.NewValidationError("agentCollectionTokenId", "agent not found")
	}
	if !agent.IsInternal {
		return nil, errors.NewValidationError("agentCollectionTokenId", "Agent is not internal")
	}
	if *agent.Owner != request.AuthenticatedWallet.Hex() {
		return nil, errors.NewValidationError("agentCollectionTokenId", "You are not the owner of this agent")
	}

	agentContexts, err := c.agentContextRequester.GetAgentContexts(
		agent.ID,
		request.Search,
		(*string)(request.SortBy),
		request.SortDirection.ToDatabaseDirection(),
	)
	if err != nil {
		return nil, err
	}

	agentContextsKeys := lo.Map(agentContexts, func(agentContext *models.AgentContext, _ int) string {
		return agentContext.Key
	})

	// Exclude context discovery contexts
	agentContextsKeys = lo.Filter(agentContextsKeys, func(key string, _ int) bool {
		return !strings.HasPrefix(key, compute.ConversationContextDiscoveryKeyPrefix) &&
			key != compute.LegacyContextDiscoveryKey
	})

	totalElements := uint(len(agentContextsKeys))

	page := *request.Page
	elementsPerPage := *request.ElementsPerPage

	start := int((page - 1) * elementsPerPage)
	end := start + int(elementsPerPage)
	if start > len(agentContextsKeys) {
		start = len(agentContextsKeys)
	}
	if end > len(agentContextsKeys) {
		end = len(agentContextsKeys)
	}

	paginatedKeys := lo.Slice(agentContextsKeys, start, end)

	return &responses.ListAgentContextKeysResponse{
		PaginateResponse: responses.PaginateResponse{
			Page:            page,
			ElementsPerPage: elementsPerPage,
			TotalElements:   totalElements,
		},
		ContextKeys: paginatedKeys,
	}, nil
}
