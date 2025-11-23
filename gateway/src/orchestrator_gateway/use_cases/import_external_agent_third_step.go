package use_cases

import (
	"fmt"
	"math/big"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/responses"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type ImportExternalAgentThirdStepInterface interface {
	Execute(request *requests.AgentByCollectionTokenIdRequest) (*responses.ImportExternalAgentThirdStepResponse, error)
}

type ImportExternalAgentThirdStep struct {
	orchestratorSettings    *settings.OrchestratorSettings
	agentRepository         repositories.AgentRepositoryInterface
	updateAgentSignerSigner interfaces.UpdateAgentSignerSignerInterface
}

func NewImportExternalAgentThirdStep(
	orchestratorSettings *settings.OrchestratorSettings,
	agentRepository repositories.AgentRepositoryInterface,
	updateAgentSignerSigner interfaces.UpdateAgentSignerSignerInterface,
) *ImportExternalAgentThirdStep {
	return &ImportExternalAgentThirdStep{
		agentRepository:         agentRepository,
		orchestratorSettings:    orchestratorSettings,
		updateAgentSignerSigner: updateAgentSignerSigner,
	}
}

func (f *ImportExternalAgentThirdStep) Execute(request *requests.AgentByCollectionTokenIdRequest) (*responses.ImportExternalAgentThirdStepResponse, error) {
	agent, err := f.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil {
		return nil, err
	}

	if agent == nil {
		return nil, errors.NewNotFoundError("Agent not found")
	}
	if agent.IsInternal {
		return nil, fmt.Errorf("agent is already internal")
	}
	if agent.ImportationStatus == nil ||
		(*agent.ImportationStatus != enums.AgentImportationStatusSecondStep &&
			*agent.ImportationStatus != enums.AgentImportationStatusThirdStep) {
		return nil, fmt.Errorf("agent importation status is not second step")
	}
	if *agent.Owner != request.AuthenticatedWallet.Hex() {
		return nil, fmt.Errorf("agent cannot be imported, caller is not the owner")
	}

	nonce := uuid.New()
	collectionAddress := common.HexToAddress(request.AgentCollectionAddress)
	collectionTokenId := big.NewInt(0).SetUint64(*agent.CollectionTokenId)
	signature, err := f.updateAgentSignerSigner.Sign(
		*collectionTokenId,
		f.orchestratorSettings.PublicKey.Hex(),
		nonce,
		collectionAddress,
	)
	if err != nil {
		return nil, err
	}

	agent.ImportationStatus = enums.AgentImportationStatusThirdStep.ToPointer()

	_, err = f.agentRepository.Save(agent)
	if err != nil {
		return nil, err
	}

	return &responses.ImportExternalAgentThirdStepResponse{
		Signature: *signature,
		Signer:    f.orchestratorSettings.PublicKey.Hex(),
	}, nil
}
