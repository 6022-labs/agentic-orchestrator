package services

import (
	"context"
	"fmt"
	"strings"

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

func (s *SlackAgentSynchronizer) Sync(
	ctx context.Context,
	collectionAddress string,
	collectionTokenId uint64,
) (*models.SlackAgent, error) {
	collectionAddress = strings.ToLower(collectionAddress)

	agentIdentity, err := s.agentRequester.GetAgent(ctx, collectionAddress, collectionTokenId)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent from orchestrator: %w", err)
	}
	ownerWallet := strings.ToLower(agentIdentity.OwnerWallet)

	slackAgentInDatabase, err := s.slackAgentRepository.GetByCollectionAddressAndCollectionTokenId(ctx,
		collectionAddress,
		collectionTokenId,
	)
	if err != nil {
		return nil, err
	}

	if slackAgentInDatabase == nil {
		brokerParticipantId, err := s.conversationBrokerRequester.RegisterAgent(
			ctx,
			collectionAddress,
			collectionTokenId,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to register agent participant with broker: %w", err)
		}

		slackAgentInDatabase = &models.SlackAgent{
			CollectionAddress:   collectionAddress,
			CollectionTokenId:   collectionTokenId,
			EnsDomain:           agentIdentity.EnsDomain,
			Role:                agentIdentity.Role,
			IconImageCid:        agentIdentity.IconImageCid,
			OwnerWallet:         ownerWallet,
			BrokerParticipantId: &brokerParticipantId,
		}

		slackAgentInDatabase, err = s.slackAgentRepository.Save(ctx, slackAgentInDatabase)
		if err != nil {
			return nil, err
		}

		return slackAgentInDatabase, nil
	}

	if slackAgentInDatabase.OwnerWallet != ownerWallet {
		slackAgentInDatabase.OwnerWallet = ownerWallet
		err = s.slackAgentRepository.UpdateOwnerWallet(
			ctx,
			collectionAddress,
			collectionTokenId,
			ownerWallet,
		)
		if err != nil {
			return nil, err
		}
	}

	if (slackAgentInDatabase.IconImageCid == nil && agentIdentity.IconImageCid != nil) ||
		(slackAgentInDatabase.IconImageCid != nil && agentIdentity.IconImageCid != nil &&
			*slackAgentInDatabase.IconImageCid != *agentIdentity.IconImageCid) {
		slackAgentInDatabase.IconImageCid = agentIdentity.IconImageCid
		err = s.slackAgentRepository.UpdateIconImageCid(ctx,
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
