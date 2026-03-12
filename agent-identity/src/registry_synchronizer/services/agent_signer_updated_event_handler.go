package services

import (
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

type AgentSignerUpdatedEventHandler struct {
	logger *log.Logger

	orchestratorSettings settings.OrchestratorSettings

	agentRepository             repositories.AgentRepositoryInterface
	agentCollectionRequester    interfaces.AgentCollectionBlockchainRequesterInterface
	orchestratorDomainResolver  interfaces.OrchestratorDomainResolverInterface
	agentPreviousSignerImporter interfaces.AgentPreviousSignerImporterInterface
}

func NewAgentSignerUpdatedEventHandler(
	logger *log.Logger,
	orchestratorSettings settings.OrchestratorSettings,
	agentRepository repositories.AgentRepositoryInterface,
	agentCollectionRequester interfaces.AgentCollectionBlockchainRequesterInterface,
	orchestratorDomainResolver interfaces.OrchestratorDomainResolverInterface,
	agentPreviousSignerImporter interfaces.AgentPreviousSignerImporterInterface,
) *AgentSignerUpdatedEventHandler {
	return &AgentSignerUpdatedEventHandler{
		logger: logger,

		orchestratorSettings: orchestratorSettings,

		agentRepository:             agentRepository,
		agentCollectionRequester:    agentCollectionRequester,
		orchestratorDomainResolver:  orchestratorDomainResolver,
		agentPreviousSignerImporter: agentPreviousSignerImporter,
	}
}

func (handler *AgentSignerUpdatedEventHandler) Handle(event *abi.AgentCollectionV1AgentSignerUpdated) error {
	handler.logger.Printf("Handling AgentSignerUpdated event for collection %s, agentId %d, new signer %s",
		event.Raw.Address.Hex(),
		event.TokenId.Uint64(),
		event.NewSigner.Hex(),
	)

	agent, err := handler.agentRepository.GetByCollectionAddressAndCollectionTokenId(event.Raw.Address.Hex(), event.TokenId.Uint64())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("failed to get agent by collection and id: " + err.Error())
	}
	if agent == nil {
		return errors.New("agent not found for collection " + event.Raw.Address.Hex() + " and id " + event.TokenId.String())
	}

	isNewSignerInternal := event.NewSigner.Cmp(handler.orchestratorSettings.PublicKey) == 0
	if isNewSignerInternal && agent.ImportationStatus != nil && *agent.ImportationStatus == enums.AgentImportationStatusThirdStep {
		if agent.CollectionTokenId == nil {
			return errors.New("agent collection token id is missing for importation")
		}

		err = handler.importFromPreviousSigners(event.Raw.Address, *agent.CollectionTokenId, agent)
		if err != nil {
			return err
		}

		agent.ImportationStatus = enums.AgentImportationStatusFinished.ToPointer()
	}

	agent.IsInternal = isNewSignerInternal

	_, err = handler.agentRepository.Save(agent)
	if err != nil {
		return errors.New("failed to update agent: " + err.Error())
	}

	return nil
}

func (handler *AgentSignerUpdatedEventHandler) importFromPreviousSigners(
	collectionAddress common.Address,
	collectionTokenId uint64,
	agent *models.Agent,
) error {
	previousSigners, err := handler.agentCollectionRequester.GetAgentPreviousSigners(
		collectionAddress,
		*new(big.Int).SetUint64(collectionTokenId),
	)
	if err != nil {
		return fmt.Errorf("failed to get previous signers: %w", err)
	}
	if len(previousSigners) == 0 {
		return errors.New("agent has no previous signers")
	}

	for signerIndex := len(previousSigners) - 1; signerIndex >= 0; signerIndex-- {
		previousSigner := previousSigners[signerIndex]
		gatewayDomain, err := handler.orchestratorDomainResolver.ResolveGatewayDomain(previousSigner)
		if err != nil {
			handler.logger.Printf("failed to resolve signer domain for %s: %v", previousSigner.Hex(), err)
			continue
		}
		if gatewayDomain == nil || *gatewayDomain == "" {
			handler.logger.Printf("missing signer domain for %s", previousSigner.Hex())
			continue
		}

		err = handler.agentPreviousSignerImporter.Import(*gatewayDomain, agent)
		if err == nil {
			return nil
		}

		handler.logger.Printf("error while trying to import agent from signer %s: %v", previousSigner.Hex(), err)
	}

	return errors.New("failed to import agent from previous signers")
}
