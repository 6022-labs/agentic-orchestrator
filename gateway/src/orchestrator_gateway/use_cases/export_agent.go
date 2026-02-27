package use_cases

import (
	"fmt"
	"math/big"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services/interfaces"
	gateway_interfaces "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/responses"
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

type ExportAgentInterface interface {
	Export(requests *requests.AgentByCollectionTokenIdRequest) (*responses.ExportAgentResponse, error)
}

type ExportAgent struct {
	agentRepository             repositories.AgentRepositoryInterface
	agentPrivateKeyStorage      interfaces.AgentPrivateKeyStorageInterface
	agentCollectionRequester    interfaces.AgentCollectionRequesterInterface
	llmConnectionRequester      gateway_interfaces.LLMConnectionRequesterInterface
	databaseConnectionRequester gateway_interfaces.DatabaseConnectionRequesterInterface
}

func NewExportAgent(
	agentRepository repositories.AgentRepositoryInterface,
	agentPrivateKeyStorage interfaces.AgentPrivateKeyStorageInterface,
	agentCollectionRequester interfaces.AgentCollectionRequesterInterface,
	llmConnectionRequester gateway_interfaces.LLMConnectionRequesterInterface,
	databaseConnectionRequester gateway_interfaces.DatabaseConnectionRequesterInterface,
) *ExportAgent {
	return &ExportAgent{
		agentRepository:             agentRepository,
		agentPrivateKeyStorage:      agentPrivateKeyStorage,
		llmConnectionRequester:      llmConnectionRequester,
		agentCollectionRequester:    agentCollectionRequester,
		databaseConnectionRequester: databaseConnectionRequester,
	}
}

func (e *ExportAgent) Export(requests *requests.AgentByCollectionTokenIdRequest) (*responses.ExportAgentResponse, error) {
	agent, err := e.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		requests.AgentCollectionAddress,
		requests.AgentCollectionTokenId,
	)
	if err != nil {
		return nil, err
	}

	collectionAddress := common.HexToAddress(requests.AgentCollectionAddress)
	collectionTokenId := big.NewInt(0).SetUint64(requests.AgentCollectionTokenId)
	agentSigner, err := e.agentCollectionRequester.GetAgentSigner(collectionAddress, *collectionTokenId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if agent == nil {
		return nil, errors.NewValidationError("agentCollectionTokenId", "agent not found")
	}
	if !agent.IsInternal {
		return nil, errors.NewValidationError("agentCollectionTokenId", "Agent is not internal")
	}
	if agentSigner.Hex() != requests.AuthenticatedWallet.Hex() {
		return nil, errors.NewValidationError("agentCollectionTokenId", "You are not the owner of this agent")
	}

	privateKey, err := e.agentPrivateKeyStorage.Get(agent.ID)
	if err != nil {
		return nil, err
	}
	if privateKey == nil {
		return nil, fmt.Errorf("private key not found")
	}

	err = e.agentPrivateKeyStorage.Delete(agent.ID)
	if err != nil {
		return nil, err
	}

	err = e.databaseConnectionRequester.Remove(agent.ID)
	if err != nil {
		return nil, err
	}

	err = e.llmConnectionRequester.Remove(agent.ID)
	if err != nil {
		return nil, err
	}

	return &responses.ExportAgentResponse{
		PrivateKey: *privateKey,
	}, nil
}
