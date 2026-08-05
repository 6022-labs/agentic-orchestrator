package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"go.uber.org/zap"
)

type SlackChannelSynchronizer struct {
	logger                 *zap.Logger
	slackChannelRequester  interfaces.SlackChannelRequesterInterface
	slackChannelRepository repositories.SlackChannelRepositoryInterface
	slackUserRepository    repositories.SlackUserRepositoryInterface
	slackUserSynchronizer  interfaces.SlackUserSynchronizerInterface
}

func NewSlackChannelSynchronizer(
	logger *zap.Logger,
	slackChannelRequester interfaces.SlackChannelRequesterInterface,
	slackChannelRepository repositories.SlackChannelRepositoryInterface,
	slackUserRepository repositories.SlackUserRepositoryInterface,
	slackUserSynchronizer interfaces.SlackUserSynchronizerInterface,
) *SlackChannelSynchronizer {
	return &SlackChannelSynchronizer{
		logger:                 logger,
		slackChannelRequester:  slackChannelRequester,
		slackChannelRepository: slackChannelRepository,
		slackUserRepository:    slackUserRepository,
		slackUserSynchronizer:  slackUserSynchronizer,
	}
}

func (s *SlackChannelSynchronizer) Sync(channelId string) (*models.SlackChannel, error) {
	s.logger.Debug("Starting channel synchronization", zap.String("channelId", channelId))

	existingChannelInDatabase, err := s.slackChannelRepository.GetById(channelId)
	if err != nil {
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
	existingChannelInDatabase.Description = existingChannelInSlack.Description

	existingChannelInDatabase, err = s.slackChannelRepository.Save(existingChannelInDatabase)
	if err != nil {
		return nil, err
	}

	if err := s.reconcileMembership(channelId, existingChannelInSlack.UserIds); err != nil {
		return nil, err
	}

	return existingChannelInDatabase, nil
}

// reconcileMembership attaches every Slack member to the channel via the
// slack_channel_users join table. Known users get an immediate synchronous
// link insert. Unknown users (no slack_users row yet) are handled in
// goroutines — the per-user synchronizer hits Slack's users.info and
// registers the user with the broker, then writes the membership row. The
// caller returns without waiting on the fan-out: agents can @mention those
// users as soon as the background work lands, not before.
func (s *SlackChannelSynchronizer) reconcileMembership(channelId string, userIds []string) error {
	if len(userIds) == 0 {
		return nil
	}

	existing, err := s.slackUserRepository.ListByIds(userIds)
	if err != nil {
		return fmt.Errorf("failed to load existing slack users: %w", err)
	}
	known := make(map[string]struct{}, len(existing))
	for _, u := range existing {
		known[u.Id] = struct{}{}
	}

	for _, userId := range userIds {
		if _, ok := known[userId]; ok {
			if err := s.slackChannelRepository.AddUser(channelId, userId); err != nil {
				return fmt.Errorf("failed to link known user %s to channel %s: %w", userId, channelId, err)
			}
			continue
		}

		go s.syncUnknownMember(channelId, userId)
	}

	return nil
}

func (s *SlackChannelSynchronizer) syncUnknownMember(channelId, userId string) {
	user, err := s.slackUserSynchronizer.Sync(userId)
	if err != nil {
		s.logger.Warn("Failed to sync unknown slack user during channel membership reconciliation",
			zap.String("channelId", channelId),
			zap.String("userId", userId),
			zap.Error(err),
		)
		return
	}
	if user == nil {
		return
	}

	err = s.slackChannelRepository.AddUser(channelId, user.Id)
	if err != nil {
		s.logger.Warn("Failed to link freshly-synced user to channel",
			zap.String("channelId", channelId),
			zap.String("userId", userId),
			zap.Error(err),
		)
	}
}
