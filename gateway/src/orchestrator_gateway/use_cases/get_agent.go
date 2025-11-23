package use_cases

import (
	"sync"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/responses"
)

type GetAgentInterface interface {
	Get(request *requests.AgentByCollectionTokenIdRequest) (*responses.AgentResponse, error)
}

type GetAgent struct {
	agentRepository             repositories.AgentRepositoryInterface
	llmConnectionRequester      interfaces.LLMConnectionRequesterInterface
	databaseConnectionRequester interfaces.DatabaseConnectionRequesterInterface
}

func NewGetAgent(
	agentRepository repositories.AgentRepositoryInterface,
	llmConnectionRequester interfaces.LLMConnectionRequesterInterface,
	databaseConnectionRequester interfaces.DatabaseConnectionRequesterInterface,
) *GetAgent {
	return &GetAgent{
		agentRepository:             agentRepository,
		llmConnectionRequester:      llmConnectionRequester,
		databaseConnectionRequester: databaseConnectionRequester,
	}
}

func (s *GetAgent) Get(request *requests.AgentByCollectionTokenIdRequest) (*responses.AgentResponse, error) {
	agent, err := s.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil {
		return nil, err
	}

	if agent == nil {
		return nil, nil // Agent not found
	}

	llmConnectionSetupStatus, databaseConnectionSetupStatus, err := getAgentSetupStatuses(
		agent,
		s.llmConnectionRequester,
		s.databaseConnectionRequester,
	)
	if err != nil {
		return nil, err
	}

	agentResponse := responses.AgentResponseFromModel(
		agent,
		llmConnectionSetupStatus.Status,
		databaseConnectionSetupStatus.Status,
		llmConnectionSetupStatus.IsCompleted,
		databaseConnectionSetupStatus.IsCompleted,
	)

	return agentResponse, nil
}

func getAgentSetupStatuses(
	agent *models.Agent,
	llmConnectionRequester interfaces.LLMConnectionRequesterInterface,
	databaseConnectionRequester interfaces.DatabaseConnectionRequesterInterface,
) (llmSetupStatus *interfaces.LLMConnectionStatus, databaseSetupStatus *interfaces.DatabaseConnectionStatus, err error) {
	var wg sync.WaitGroup

	var llmStatus *interfaces.LLMConnectionStatus = &interfaces.LLMConnectionStatus{
		IsCompleted: false,
		Status:      nil,
	}
	var dbStatus *interfaces.DatabaseConnectionStatus = &interfaces.DatabaseConnectionStatus{
		IsCompleted: false,
		Status:      nil,
	}

	var llmErr error
	var dbErr error

	if agent.Role != enums.AgentRoleHuman {
		wg.Add(1)
		go func() {
			defer wg.Done()
			llmStatus, llmErr = llmConnectionRequester.GetSetupStatus(agent.ID)
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		dbStatus, dbErr = databaseConnectionRequester.GetSetupStatus(agent.ID)
	}()

	wg.Wait()

	if llmErr != nil {
		return nil, nil, llmErr
	}
	if dbErr != nil {
		return nil, nil, dbErr
	}
	return llmStatus, dbStatus, nil
}
