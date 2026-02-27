package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services/interfaces"
	gateway_interfaces "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/responses"
)

type ImportExternalAgentSecondStepInterface interface {
	Execute(request *requests.ImportExternalAgentSecondStepRequest) (*responses.ImportExternalAgentSecondStepResponse, error)
}

type ImportExternalAgentSecondStep struct {
	agentRepository                 repositories.AgentRepositoryInterface
	llmConnectionRequester          gateway_interfaces.LLMConnectionRequesterInterface
	databaseConnectionRequester     gateway_interfaces.DatabaseConnectionRequesterInterface
	externalAgentInformationFetcher interfaces.ExternalAgentInformationFetcherInterface
}

func NewImportExternalAgentSecondStep(
	agentRepository repositories.AgentRepositoryInterface,
	llmConnectionRequester gateway_interfaces.LLMConnectionRequesterInterface,
	databaseConnectionRequester gateway_interfaces.DatabaseConnectionRequesterInterface,
	externalAgentInformationFetcher interfaces.ExternalAgentInformationFetcherInterface,
) *ImportExternalAgentSecondStep {
	return &ImportExternalAgentSecondStep{
		agentRepository:                 agentRepository,
		llmConnectionRequester:          llmConnectionRequester,
		databaseConnectionRequester:     databaseConnectionRequester,
		externalAgentInformationFetcher: externalAgentInformationFetcher,
	}
}

func (i *ImportExternalAgentSecondStep) Execute(request *requests.ImportExternalAgentSecondStepRequest) (*responses.ImportExternalAgentSecondStepResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	agent, err := i.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil {
		return nil, err
	}

	if agent.IsInternal {
		return nil, fmt.Errorf("agent is internal, cannot be imported")
	}
	if agent.ImportationStatus == nil || *agent.ImportationStatus != enums.AgentImportationStatusFirstStep {
		return nil, fmt.Errorf("agent importation status is not first step")
	}
	if *agent.Owner != request.AuthenticatedWallet.Hex() {
		return nil, fmt.Errorf("agent cannot be imported, caller is not the owner")
	}

	err = i.llmConnectionRequester.Save(&gateway_interfaces.SaveLLMConnectionRequest{
		AgentId:  agent.ID,
		Model:    request.LLMModel,
		ApiKey:   request.LLMApiKey,
		BaseUrl:  request.LLMBaseUrl,
		Platform: request.LLMPlatform,
	})
	if err != nil {
		return nil, err
	}

	databaseConnectionResult, err := i.databaseConnectionRequester.Create(&gateway_interfaces.CreateDatabaseConnectionRequest{
		AgentId: agent.ID,
		Host:    request.DatabaseHost,
		Port:    request.DatabasePort,
		Driver:  request.DatabaseDriver,
	})
	if err != nil {
		return nil, err
	}

	agent.ImportationStatus = enums.AgentImportationStatusSecondStep.ToPointer()

	_, err = i.agentRepository.Save(agent)
	if err != nil {
		return nil, err
	}

	return &responses.ImportExternalAgentSecondStepResponse{
		DatabaseName:     databaseConnectionResult.DatabaseName,
		DatabaseUsername: databaseConnectionResult.DatabaseUsername,
		DatabasePassword: databaseConnectionResult.DatabasePassword,
	}, nil
}
