package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"gorm.io/gorm"
)

type SlackAgentSynchronizer struct {
	agentRequester       interfaces.AgentRequesterInterface
	slackAgentRepository repositories.SlackAgentRepositoryInterface
}

func NewSlackAgentSynchronizer(
	agentRequester interfaces.AgentRequesterInterface,
	slackAgentRepository repositories.SlackAgentRepositoryInterface,
) *SlackAgentSynchronizer {
	return &SlackAgentSynchronizer{
		agentRequester:       agentRequester,
		slackAgentRepository: slackAgentRepository,
	}
}

func (s *SlackAgentSynchronizer) Sync(collectionAddress string, collectionTokenId uint64) (*models.SlackAgent, error) {
	agentInOrchestrator, err := s.agentRequester.GetAgent(collectionAddress, collectionTokenId)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent from orchestrator: %w", err)
	}

	slackAgentInDatabase, err := s.slackAgentRepository.GetByCollectionAddressAndCollectionTokenId(collectionAddress, collectionTokenId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if slackAgentInDatabase == nil {
		slackAgentInDatabase = &models.SlackAgent{
			CollectionAddress: collectionAddress,
			CollectionTokenId: collectionTokenId,
			Name:              agentInOrchestrator.Name,
			EnsDomain:         agentInOrchestrator.EnsDomain,
			Role:              agentInOrchestrator.Role,
			IconImageCid:      agentInOrchestrator.IconImageCid,
			OwnerWallet:       agentInOrchestrator.OwnerWallet,
			IsProcessing:      false,
			LastProcessedAt:   nil,
		}

		slackAgentInDatabase, err = s.slackAgentRepository.Save(slackAgentInDatabase)
		if err != nil {
			return nil, err
		}

		return slackAgentInDatabase, nil
	}

	// Update only variable fields that come from the orchestrator
	if slackAgentInDatabase.OwnerWallet != agentInOrchestrator.OwnerWallet {
		slackAgentInDatabase.OwnerWallet = agentInOrchestrator.OwnerWallet
		err = s.slackAgentRepository.UpdateOwnerWallet(collectionAddress, collectionTokenId, agentInOrchestrator.OwnerWallet)
		if err != nil {
			return nil, err
		}
	}

	if (slackAgentInDatabase.IconImageCid == nil && agentInOrchestrator.IconImageCid != nil) ||
		(slackAgentInDatabase.IconImageCid != nil && agentInOrchestrator.IconImageCid != nil && *slackAgentInDatabase.IconImageCid != *agentInOrchestrator.IconImageCid) {
		slackAgentInDatabase.IconImageCid = agentInOrchestrator.IconImageCid
		err = s.slackAgentRepository.UpdateIconImageCid(collectionAddress, collectionTokenId, agentInOrchestrator.IconImageCid)
		if err != nil {
			return nil, err
		}
	}

	return slackAgentInDatabase, nil
}
