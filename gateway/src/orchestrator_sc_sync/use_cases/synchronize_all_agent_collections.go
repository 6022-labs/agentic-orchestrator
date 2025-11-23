package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/services/interfaces"
)

type SynchronizeAllAgentCollectionsInterface interface {
	Synchronize() error
}

type SynchronizeAllAgentCollections struct {
	agentCollectionSynchronizer      services.AgentCollectionSynchronizerInterface
	agentCollectionRepository        repositories.AgentCollectionRepositoryInterface
	agentCollectionsManagerRequester interfaces.AgentCollectionsManagerRequesterInterface
}

func NewSynchronizeAllAgentCollections(
	agentCollectionSynchronizer services.AgentCollectionSynchronizerInterface,
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface,
	agentCollectionsManagerRequester interfaces.AgentCollectionsManagerRequesterInterface,
) *SynchronizeAllAgentCollections {
	return &SynchronizeAllAgentCollections{
		agentCollectionRepository:        agentCollectionRepository,
		agentCollectionSynchronizer:      agentCollectionSynchronizer,
		agentCollectionsManagerRequester: agentCollectionsManagerRequester,
	}
}

func (s *SynchronizeAllAgentCollections) Synchronize() error {
	addresses, err := s.agentCollectionsManagerRequester.GetCollectionAddresses()
	if err != nil {
		return err
	}

	for _, address := range addresses {
		collectionAlreadyExistsInPersistence, err := s.agentCollectionRepository.ExistsByCollectionAddress(address.Hex())
		if err != nil {
			return fmt.Errorf("failed to check if collection %s exists in persistence: %w", address.Hex(), err)
		}

		if !collectionAlreadyExistsInPersistence {
			err = s.agentCollectionSynchronizer.SyncByCollectionAddress(address.Hex())
			if err != nil {
				return fmt.Errorf("failed to synchronize collection %s: %w", address.Hex(), err)
			}
		}
	}

	return nil
}
