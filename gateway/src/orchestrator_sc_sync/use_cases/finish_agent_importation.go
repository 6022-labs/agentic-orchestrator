package use_cases

import (
	"fmt"
	"math/big"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services/interfaces"
	sc_sync_services "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/use_cases/requests"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type FinishAgentImportationInterface interface {
	Finish(request *requests.FinishAgentImportationRequest) error
}

type FinishAgentImportation struct {
	logger                      *zap.Logger
	agentRepository             repositories.AgentRepositoryInterface
	agentCollectionRequester    interfaces.AgentCollectionRequesterInterface
	orchestratorDomainResolver  interfaces.OrchestratorDomainResolverInterface
	agentPreviousSignerImporter sc_sync_services.AgentPreviousSignerImporterInterface
}

func NewFinishAgentImportation(
	logger *zap.Logger,
	agentRepository repositories.AgentRepositoryInterface,
	agentCollectionRequester interfaces.AgentCollectionRequesterInterface,
	orchestratorDomainResolver interfaces.OrchestratorDomainResolverInterface,
	agentPreviousSignerImporter sc_sync_services.AgentPreviousSignerImporterInterface,
) *FinishAgentImportation {
	return &FinishAgentImportation{
		logger:                      logger,
		agentRepository:             agentRepository,
		orchestratorDomainResolver:  orchestratorDomainResolver,
		agentPreviousSignerImporter: agentPreviousSignerImporter,
		agentCollectionRequester:    agentCollectionRequester,
	}
}

func (i *FinishAgentImportation) Finish(request *requests.FinishAgentImportationRequest) error {
	i.logger.Info("importing external agent", zap.Uint64("agent_collection_token_id", request.AgentCollectionTokenId))

	agent, err := i.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil {
		return err
	}
	if agent == nil {
		return fmt.Errorf("agent not found")
	}

	if agent.IsInternal {
		return fmt.Errorf("agent is already internal")
	}

	if agent.ImportationStatus == nil || *agent.ImportationStatus != enums.AgentImportationStatusThirdStep {
		return fmt.Errorf("agent is not in the correct importation status")
	}

	collectionAddress := common.HexToAddress(request.AgentCollectionAddress)
	collectionTokenId := new(big.Int).SetUint64(*agent.CollectionTokenId)
	previousSigners, err := i.agentCollectionRequester.GetAgentPreviousSigners(collectionAddress, *collectionTokenId)
	if err != nil {
		return err
	}
	if len(previousSigners) == 0 {
		return fmt.Errorf("agent has no previous signers")
	}

	importSucceed := false
	previousSignerIndex := len(previousSigners) - 1

	for !importSucceed && previousSignerIndex >= 0 {
		previousSigner := previousSigners[previousSignerIndex]

		err = i.tryImportAgentFromPreviousSigner(previousSigner, agent)
		if err == nil {
			importSucceed = true
			continue
		}

		i.logger.Warn("error while trying to import agent", zap.Error(err))
		previousSignerIndex--
	}

	if !importSucceed {
		return fmt.Errorf("failed to import agent from previous signers")
	}

	agent.IsInternal = true
	agent.ImportationStatus = enums.AgentImportationStatusFinished.ToPointer()

	_, err = i.agentRepository.Save(agent)
	if err != nil {
		return err
	}

	return nil
}

func (i *FinishAgentImportation) tryImportAgentFromPreviousSigner(previousSigner common.Address, agent *models.Agent) error {
	signerAIDomain, err := i.orchestratorDomainResolver.ResolveAIDomain(previousSigner)
	if err != nil {
		return err
	}

	i.logger.Info(
		"signer domain found",
		zap.String("signer", previousSigner.Hex()),
		zap.String("domain", *signerAIDomain),
	)

	return i.agentPreviousSignerImporter.Import(*signerAIDomain, agent)
}
