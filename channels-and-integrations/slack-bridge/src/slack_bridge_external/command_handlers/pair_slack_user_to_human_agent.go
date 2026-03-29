package command_handlers

import (
	"errors"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PairSlackUserToHumanAgentInterface interface {
	Pair(request *requests.PairSlackUserToHumanAgentRequest) error
}

type PairSlackUserToHumanAgent struct {
	logger                      *zap.Logger
	slackUserRepository         repositories.SlackUserRepositoryInterface
	slackAgentSynchronizer      interfaces.SlackAgentSynchronizerInterface
	slackUserSynchronizerCached interfaces.SlackUserSynchronizerCachedInterface
}

func NewPairSlackUserToHumanAgent(
	logger *zap.Logger,
	slackUserRepository repositories.SlackUserRepositoryInterface,
	slackAgentSynchronizer interfaces.SlackAgentSynchronizerInterface,
	slackUserSynchronizerCached interfaces.SlackUserSynchronizerCachedInterface,
) *PairSlackUserToHumanAgent {
	return &PairSlackUserToHumanAgent{
		logger:                      logger,
		slackUserRepository:         slackUserRepository,
		slackAgentSynchronizer:      slackAgentSynchronizer,
		slackUserSynchronizerCached: slackUserSynchronizerCached,
	}
}

func (p *PairSlackUserToHumanAgent) Pair(request *requests.PairSlackUserToHumanAgentRequest) error {
	requesterUser, err := p.slackUserSynchronizerCached.Sync(request.RequesterSlackUserId)
	if err != nil {
		return fmt.Errorf("failed to sync requester user : %w", err)
	}
	if requesterUser.WalletAddress == nil {
		return fmt.Errorf("requester needs to verify their wallet first before pairing a human agent")
	}

	agent, err := p.slackAgentSynchronizer.Sync(request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId)
	if err != nil {
		return fmt.Errorf("failed to sync human agent: %w", err)
	}
	if agent.Role != enums.AgentRoleHuman {
		return fmt.Errorf("the specified agent %s %d is not a human agent, role is: %s", request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId, agent.Role)
	}
	if agent.OwnerWallet != *requesterUser.WalletAddress {
		return fmt.Errorf("requester does not own the human agent %s %d, owner is: %s", request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId, agent.OwnerWallet)
	}

	alreadyPairedUserForAgent, err := p.slackUserRepository.GetByPairedHumanAgent(
		request.HumanAgentCollectionAddress,
		request.HumanAgentCollectionTokenId,
	)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check if human agent is already paired: %w", err)
	}
	if alreadyPairedUserForAgent != nil {
		return fmt.Errorf("the specified human agent %s %d is already paired to another slack user", request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId)
	}

	slackUserToPair, err := p.slackUserSynchronizerCached.Sync(request.SlackUserId)
	if err != nil {
		return fmt.Errorf("failed to sync slack user to pair: %w", err)
	}

	slackUserToPair.PairedHumanAgentCollectionAddress = &request.HumanAgentCollectionAddress
	slackUserToPair.PairedHumanAgentCollectionTokenId = &request.HumanAgentCollectionTokenId

	_, err = p.slackUserRepository.Save(slackUserToPair)
	if err != nil {
		return fmt.Errorf("failed to save pairing: %w", err)
	}

	return nil
}
