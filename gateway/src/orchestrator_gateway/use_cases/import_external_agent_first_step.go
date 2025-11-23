package use_cases

import (
	"fmt"
	"math/big"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/responses"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type ImportExternalAgentFirstStepInterface interface {
	Execute(request *requests.AgentByCollectionTokenIdRequest) (*responses.ImportExternalAgentFirstStepResponse, error)
}

type ImportExternalAgentFirstStep struct {
	logger                          *zap.Logger
	agentRepository                 repositories.AgentRepositoryInterface
	agentCollectionRequester        interfaces.AgentCollectionRequesterInterface
	orchestratorDomainResolver      interfaces.OrchestratorDomainResolverInterface
	externalAgentInformationFetcher interfaces.ExternalAgentInformationFetcherInterface
}

func NewImportExternalAgentFirstStep(
	logger *zap.Logger,
	agentRepository repositories.AgentRepositoryInterface,
	agentCollectionRequester interfaces.AgentCollectionRequesterInterface,
	orchestratorDomainResolver interfaces.OrchestratorDomainResolverInterface,
	externalAgentInformationFetcher interfaces.ExternalAgentInformationFetcherInterface,
) *ImportExternalAgentFirstStep {
	return &ImportExternalAgentFirstStep{
		logger:                          logger,
		agentRepository:                 agentRepository,
		agentCollectionRequester:        agentCollectionRequester,
		orchestratorDomainResolver:      orchestratorDomainResolver,
		externalAgentInformationFetcher: externalAgentInformationFetcher,
	}
}

func (i *ImportExternalAgentFirstStep) Execute(request *requests.AgentByCollectionTokenIdRequest) (*responses.ImportExternalAgentFirstStepResponse, error) {
	agent, err := i.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil {
		return nil, err
	}
	if agent.IsInternal {
		return nil, fmt.Errorf("agent is already internal, cannot be imported")
	}
	if agent.ImportationStatus != nil && *agent.ImportationStatus != enums.AgentImportationStatusFirstStep {
		return nil, fmt.Errorf("agent importation status is not first step")
	}

	collectionAddress := common.HexToAddress(request.AgentCollectionAddress)
	collectionTokenId := new(big.Int).SetUint64(*agent.CollectionTokenId)
	owner, err := i.agentCollectionRequester.GetAgentOwner(collectionAddress, *collectionTokenId)
	if err != nil {
		return nil, err
	}
	if *owner != *request.AuthenticatedWallet {
		return nil, fmt.Errorf("authenticated wallet is not the owner of the agent")
	}

	previousSigners, err := i.agentCollectionRequester.GetAgentPreviousSigners(collectionAddress, *collectionTokenId)
	if err != nil {
		return nil, err
	}
	if len(previousSigners) == 0 {
		return nil, fmt.Errorf("agent has no previous signers")
	}

	metadata, err := i.fetchAgentMetadata(previousSigners, *agent.CollectionTokenId)
	if err != nil {
		return nil, err
	}

	newOwner := owner.Hex()
	agent.Owner = &newOwner // Update the owner in case of a de-sync
	agent.ImportationStatus = enums.AgentImportationStatusFirstStep.ToPointer()

	_, err = i.agentRepository.Save(agent)
	if err != nil {
		return nil, err
	}

	return &responses.ImportExternalAgentFirstStepResponse{
		LLMModel:    metadata.LLMModel,
		LLMPlatform: metadata.LLMPlatform,
	}, nil
}

func (i *ImportExternalAgentFirstStep) fetchAgentMetadata(previousSigners []common.Address, collectionTokenId uint64) (*interfaces.AgentMetadataInformation, error) {
	var err error
	var metadata *interfaces.AgentMetadataInformation

	previousSignerIndex := len(previousSigners) - 1

	for metadata == nil && previousSignerIndex >= 0 {
		previousSigner := previousSigners[previousSignerIndex]

		metadata, err = i.fetchAgentMetadataFromSigner(previousSigner, collectionTokenId)
		if err == nil {
			continue
		}

		i.logger.Warn(
			"cannot import agent",
			zap.String("signer", previousSigner.Hex()),
			zap.Error(err),
		)
		previousSignerIndex--
	}

	if metadata == nil {
		return nil, fmt.Errorf("cannot fetch metadata from previous signers")
	}

	return metadata, nil
}

func (i *ImportExternalAgentFirstStep) fetchAgentMetadataFromSigner(previousSigner common.Address, collectionTokenId uint64) (*interfaces.AgentMetadataInformation, error) {
	signerAIDomain, err := i.orchestratorDomainResolver.ResolveAIDomain(previousSigner)
	if err != nil {
		return nil, err
	}

	i.logger.Info(
		"signer domain found",
		zap.String("signer", previousSigner.Hex()),
		zap.String("domain", *signerAIDomain),
	)

	isInternalAgentForSigner, err := i.externalAgentInformationFetcher.IsInternalAgent(*signerAIDomain, collectionTokenId)
	if err != nil {
		return nil, err
	}
	if !*isInternalAgentForSigner {
		return nil, fmt.Errorf("agent is not known by the signer")
	}

	return i.externalAgentInformationFetcher.FetchAgentMetadataInformation(*signerAIDomain, collectionTokenId)
}
