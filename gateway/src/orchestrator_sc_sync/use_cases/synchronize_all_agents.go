package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/services"
	"github.com/ethereum/go-ethereum/common"
)

type SynchronizeAllAgentsInterface interface {
	Synchronize() error
}

type SynchronizeAllAgents struct {
	agentSynchronizer         services.AgentSynchronizerInterface
	agentRepository           repositories.AgentRepositoryInterface
	agentCollectionRequester  interfaces.AgentCollectionRequesterInterface
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface
}

func NewSynchronizeAllAgents(
	agentSynchronizer services.AgentSynchronizerInterface,
	agentRepository repositories.AgentRepositoryInterface,
	agentCollectionRequester interfaces.AgentCollectionRequesterInterface,
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface,
) *SynchronizeAllAgents {
	return &SynchronizeAllAgents{
		agentRepository:           agentRepository,
		agentSynchronizer:         agentSynchronizer,
		agentCollectionRequester:  agentCollectionRequester,
		agentCollectionRepository: agentCollectionRepository,
	}
}

func (s *SynchronizeAllAgents) Synchronize() error {
	collections, err := s.agentCollectionRepository.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get all agent collections: %w", err)
	}

	for _, collection := range collections {
		collectionAddress := common.HexToAddress(collection.Address)
		nextTokenId, err := s.agentCollectionRequester.GetNextTokenId(collectionAddress)
		if err != nil {
			return fmt.Errorf("failed to get next token ID for collection %s: %w", collection.Address, err)
		}

		nextTokenIdInt := nextTokenId.Uint64()
		for tokenId := uint64(1); tokenId < nextTokenIdInt; tokenId++ {
			exists, err := s.agentRepository.ExistsByCollectionAddressAndCollectionTokenId(collection.Address, tokenId)
			if err != nil {
				return fmt.Errorf("failed to check existence of agent with collection %s and token ID %d: %w", collection.Address, tokenId, err)
			}
			if exists {
				continue
			}

			err = s.agentSynchronizer.Sync(collection.Address, tokenId)
			if err != nil {
				return fmt.Errorf("failed to synchronize agent with collection %s and token ID %d: %w", collection.Address, tokenId, err)
			}
		}
	}

	return nil
}
