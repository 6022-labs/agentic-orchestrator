package services

import (
	"errors"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SlackChannelSynchronizer struct {
	logger                 *zap.Logger
	slackChannelRequester  interfaces.SlackChannelRequesterInterface
	slackChannelRepository repositories.SlackChannelRepositoryInterface
}

func NewSlackChannelSynchronizer(
	logger *zap.Logger,
	slackChannelRequester interfaces.SlackChannelRequesterInterface,
	slackChannelRepository repositories.SlackChannelRepositoryInterface,
) *SlackChannelSynchronizer {
	return &SlackChannelSynchronizer{
		logger:                 logger,
		slackChannelRequester:  slackChannelRequester,
		slackChannelRepository: slackChannelRepository,
	}
}

func (s *SlackChannelSynchronizer) Sync(channelId string) (*models.SlackChannel, error) {
	s.logger.Debug("Starting channel synchronization", zap.String("channelId", channelId))

	existingChannelInDatabase, err := s.slackChannelRepository.GetById(channelId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	existingChannelInSlack, err := s.slackChannelRequester.GetChannelInformation(channelId)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel information from Slack: %w", err)
	}
	if existingChannelInSlack == nil {
		return nil, fmt.Errorf("channel with ID %s not found in Slack", channelId)
	}

	if existingChannelInDatabase == nil {
		existingChannelInDatabase = &models.SlackChannel{
			Id: channelId,
		}

		if !existingChannelInSlack.IsMember {
			err = s.slackChannelRequester.JoinChannel(channelId)
			if err != nil {
				return nil, err
			}
		}
	}

	existingChannelInDatabase.Name = existingChannelInSlack.Name
	existingChannelInDatabase.Members = existingChannelInSlack.Members
	existingChannelInDatabase.Description = existingChannelInSlack.Description

	existingChannelInDatabase, err = s.slackChannelRepository.Save(existingChannelInDatabase)
	if err != nil {
		return nil, err
	}

	return existingChannelInDatabase, nil
}
