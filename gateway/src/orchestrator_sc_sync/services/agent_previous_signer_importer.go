package services

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/services/interfaces"
	"go.uber.org/zap"
)

type AgentPreviousSignerImporterInterface interface {
	Import(previousSignerDomain string, agent *models.Agent) error
}

type AgentPreviousSignerImporter struct {
	logger                          *zap.Logger
	agentRepository                 repositories.AgentRepositoryInterface
	agentPrivateKeyStorage          interfaces.AgentPrivateKeyStorageInterface
	externalAgentInformationFetcher interfaces.ExternalAgentInformationFetcherInterface
}

func NewPreviousAgentSignerImporter(
	logger *zap.Logger,
	agentRepository repositories.AgentRepositoryInterface,
	agentPrivateKeyStorage interfaces.AgentPrivateKeyStorageInterface,
	externalAgentInformationFetcher interfaces.ExternalAgentInformationFetcherInterface,
) *AgentPreviousSignerImporter {
	return &AgentPreviousSignerImporter{
		logger:                          logger,
		agentRepository:                 agentRepository,
		agentPrivateKeyStorage:          agentPrivateKeyStorage,
		externalAgentInformationFetcher: externalAgentInformationFetcher,
	}
}

func (asi *AgentPreviousSignerImporter) Import(previousSignerDomain string, agent *models.Agent) error {
	isInternalAgentForSigner, err := asi.externalAgentInformationFetcher.IsInternalAgent(previousSignerDomain, *agent.CollectionTokenId)
	if err != nil {
		return err
	}
	if !*isInternalAgentForSigner {
		return fmt.Errorf("agent is not internal for signer %s", previousSignerDomain)
	}

	asi.logger.Info(
		"agent founded",
		zap.Uint64("collectionTokenId", *agent.CollectionTokenId),
		zap.String("domain", previousSignerDomain),
	)

	agentImportInformation, err := asi.externalAgentInformationFetcher.ImportAgent(previousSignerDomain, *agent.CollectionTokenId)
	if err != nil {
		return err
	}

	asi.logger.Debug(
		"agent database information fetched",
		zap.Uint64("collectionTokenId", *agent.CollectionTokenId),
		zap.Any("database", *agentImportInformation),
	)

	err = asi.agentPrivateKeyStorage.Store(agent.ID, agentImportInformation.PrivateKey)
	if err != nil {
		return err
	}

	return nil
}
