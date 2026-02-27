package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services/interfaces"
	"github.com/ethereum/go-ethereum/common"
)

type AgentCollectionSynchronizerInterface interface {
	Sync(collectionAddress, collectionName, collectionSymbol string) error
	SyncByCollectionAddress(collectionAddress string) error
}

type AgentCollectionSynchronizer struct {
	agentCollectionRequester  interfaces.AgentCollectionRequesterInterface
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface
}

func NewAgentCollectionSynchronizer(
	agentCollectionRequester interfaces.AgentCollectionRequesterInterface,
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface,
) *AgentCollectionSynchronizer {
	return &AgentCollectionSynchronizer{
		agentCollectionRequester:  agentCollectionRequester,
		agentCollectionRepository: agentCollectionRepository,
	}
}

func (u *AgentCollectionSynchronizer) Sync(collectionAddress, collectionName, collectionSymbol string) error {
	agentCollection := &models.AgentCollection{
		Address: collectionAddress,
		Name:    collectionName,
		Symbol:  collectionSymbol,
	}

	_, err := u.agentCollectionRepository.Save(agentCollection)
	if err != nil {
		return err
	}

	return nil
}

func (u *AgentCollectionSynchronizer) SyncByCollectionAddress(collectionAddress string) error {
	if !common.IsHexAddress(collectionAddress) {
		return fmt.Errorf("invalid collection address: %s", collectionAddress)
	}

	collectionAddressHex := common.HexToAddress(collectionAddress)

	collectionName, err := u.agentCollectionRequester.GetName(collectionAddressHex)
	if err != nil {
		return err
	}

	collectionSymbol, err := u.agentCollectionRequester.GetSymbol(collectionAddressHex)
	if err != nil {
		return err
	}

	return u.Sync(collectionAddress, *collectionName, *collectionSymbol)
}
