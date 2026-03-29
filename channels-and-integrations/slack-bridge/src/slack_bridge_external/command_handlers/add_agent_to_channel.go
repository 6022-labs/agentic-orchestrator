package command_handlers

import (
	"errors"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type AddAgentToChannelInterface interface {
	Add(request *requests.AddAgentToChannelRequest) error
}

type AddAgentToChannel struct {
	slackAgentSynchronizer   interfaces.SlackAgentSynchronizerInterface
	slackChannelSynchronizer interfaces.SlackChannelSynchronizerInterface
	slackChannelRepository   repositories.SlackChannelRepositoryInterface
}

func NewAddAgentToChannel(
	slackAgentSynchronizer interfaces.SlackAgentSynchronizerInterface,
	slackChannelRepository repositories.SlackChannelRepositoryInterface,
	slackChannelSynchronizer interfaces.SlackChannelSynchronizerInterface,
) *AddAgentToChannel {
	return &AddAgentToChannel{
		slackAgentSynchronizer:   slackAgentSynchronizer,
		slackChannelRepository:   slackChannelRepository,
		slackChannelSynchronizer: slackChannelSynchronizer,
	}
}

func (h *AddAgentToChannel) Add(request *requests.AddAgentToChannelRequest) error {
	slackChannel, err := h.slackChannelRepository.GetById(request.ChannelID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to get channel with ID %s: %w", request.ChannelID, err)
	}
	if slackChannel != nil {
		channelContainAgent := lo.ContainsBy(slackChannel.Agents, func(agent models.SlackAgent) bool {
			return agent.CollectionAddress == request.CollectionAddress && agent.CollectionTokenId == request.CollectionTokenId
		})

		if channelContainAgent { // Do not add the agent if it's already in the channel
			return nil
		}
	}

	slackChannel, err = h.slackChannelSynchronizer.Sync(request.ChannelID)
	if err != nil {
		return err
	}

	slackAgent, err := h.slackAgentSynchronizer.Sync(request.CollectionAddress, request.CollectionTokenId)
	if err != nil {
		return fmt.Errorf("failed to sync agent with collection address %s and token ID %d: %w",
			request.CollectionAddress,
			request.CollectionTokenId,
			err,
		)
	}

	_, err = h.slackChannelRepository.AddAgent(slackChannel, slackAgent)
	if err != nil {
		return fmt.Errorf("failed to add agent to channel: %w", err)
	}

	return nil
}
