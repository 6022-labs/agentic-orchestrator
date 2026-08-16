package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	slack_bridge_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
)

type SlackAgentSynchronizer struct {
	agentRequester              interfaces.AgentRequesterInterface
	slackAgentRepository        repositories.SlackAgentRepositoryInterface
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface
}

func NewSlackAgentSynchronizer(
	agentRequester interfaces.AgentRequesterInterface,
	slackAgentRepository repositories.SlackAgentRepositoryInterface,
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface,
) *SlackAgentSynchronizer {
	return &SlackAgentSynchronizer{
		agentRequester:              agentRequester,
		slackAgentRepository:        slackAgentRepository,
		conversationBrokerRequester: conversationBrokerRequester,
	}
}

func (s *SlackAgentSynchronizer) Sync(collectionAddress string, collectionTokenId uint64) (*models.SlackAgent, error) {
	agentIdentity, err := s.agentRequester.GetAgent(collectionAddress, collectionTokenId)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent from orchestrator: %w", err)
	}

	slackAgentInDatabase, err := s.slackAgentRepository.GetByCollectionAddressAndCollectionTokenId(
		collectionAddress,
		collectionTokenId,
	)
	if err != nil {
		return nil, err
	}

	if slackAgentInDatabase == nil {
		brokerParticipantId, err := s.conversationBrokerRequester.RegisterAgent(collectionAddress, collectionTokenId)
		if err != nil {
			return nil, fmt.Errorf("failed to register agent participant with broker: %w", err)
		}

		slackAgentInDatabase = &models.SlackAgent{
			CollectionAddress:   collectionAddress,
			CollectionTokenId:   collectionTokenId,
			EnsDomain:           agentIdentity.EnsDomain,
			Role:                agentIdentity.Role,
			IconImageCid:        agentIdentity.IconImageCid,
			OwnerWallet:         agentIdentity.OwnerWallet,
			BrokerParticipantId: &brokerParticipantId,
		}

		slackAgentInDatabase, err = s.slackAgentRepository.Save(slackAgentInDatabase)
		if err != nil {
			return nil, err
		}

		return slackAgentInDatabase, nil
	}

	if slackAgentInDatabase.OwnerWallet != agentIdentity.OwnerWallet {
		slackAgentInDatabase.OwnerWallet = agentIdentity.OwnerWallet
		err = s.slackAgentRepository.UpdateOwnerWallet(collectionAddress, collectionTokenId, agentIdentity.OwnerWallet)
		if err != nil {
			return nil, err
		}
	}

	if (slackAgentInDatabase.IconImageCid == nil && agentIdentity.IconImageCid != nil) ||
		(slackAgentInDatabase.IconImageCid != nil && agentIdentity.IconImageCid != nil &&
			*slackAgentInDatabase.IconImageCid != *agentIdentity.IconImageCid) {
		slackAgentInDatabase.IconImageCid = agentIdentity.IconImageCid
		err = s.slackAgentRepository.UpdateIconImageCid(
			collectionAddress,
			collectionTokenId,
			agentIdentity.IconImageCid,
		)
		if err != nil {
			return nil, err
		}
	}

	return slackAgentInDatabase, nil
}
