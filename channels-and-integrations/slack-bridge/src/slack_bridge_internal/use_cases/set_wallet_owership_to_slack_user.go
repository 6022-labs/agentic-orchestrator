package use_cases

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_internal/use_cases/requests"
)

type SetWalletOwnershipToSlackUserInterface interface {
	Set(request *requests.SetWalletOwnershipToSlackUserRequest) error
}

type SetWalletOwnershipToSlackUser struct {
	slackUserRepository      repositories.SlackUserRepositoryInterface
	walletOwnershipChallenge interfaces.WalletOwnershipChallengeInterface
}

func NewSetWalletOwnershipToSlackUser(
	slackUserRepository repositories.SlackUserRepositoryInterface,
	walletOwnershipChallenge interfaces.WalletOwnershipChallengeInterface,
) *SetWalletOwnershipToSlackUser {
	return &SetWalletOwnershipToSlackUser{
		slackUserRepository:      slackUserRepository,
		walletOwnershipChallenge: walletOwnershipChallenge,
	}
}

func (uc *SetWalletOwnershipToSlackUser) Set(request *requests.SetWalletOwnershipToSlackUserRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	userId, err := uc.walletOwnershipChallenge.ConsumeChallenge(&request.Challenge)
	if err != nil {
		return err
	}

	slackUser, err := uc.slackUserRepository.GetById(*userId)
	if err != nil {
		return err
	}

	slackUser.WalletAddress = &request.WalletAddress

	_, err = uc.slackUserRepository.Save(slackUser)
	if err != nil {
		return err
	}

	return nil
}
