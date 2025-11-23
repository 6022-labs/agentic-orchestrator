package services

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	external_interfaces "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/settings"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const CHANNEL_SYNC_CACHE_KEY_TEMPLATE = "slack_channel_sync:%s"

type SlackChannelSynchronizerCached struct {
	logger                   *zap.Logger
	cacheSettings            *settings.CacheSettings
	cacheStorage             interfaces.CacheStorageInterface
	slackChannelRepository   repositories.SlackChannelRepositoryInterface
	slackChannelSynchronizer external_interfaces.SlackChannelSynchronizerInterface
}

func NewSlackChannelSynchronizerCached(
	logger *zap.Logger,
	cacheSettings *settings.CacheSettings,
	cacheStorage interfaces.CacheStorageInterface,
	slackChannelRepository repositories.SlackChannelRepositoryInterface,
	slackChannelSynchronizer external_interfaces.SlackChannelSynchronizerInterface,
) *SlackChannelSynchronizerCached {
	return &SlackChannelSynchronizerCached{
		logger:                   logger,
		cacheStorage:             cacheStorage,
		cacheSettings:            cacheSettings,
		slackChannelRepository:   slackChannelRepository,
		slackChannelSynchronizer: slackChannelSynchronizer,
	}
}

func (s *SlackChannelSynchronizerCached) Sync(channelId string) (*models.SlackChannel, error) {
	inDatabaseChannel, err := s.slackChannelRepository.GetById(channelId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("error getting channel %s from database: %w", channelId, err)
	}

	channelSyncKey := fmt.Sprintf(CHANNEL_SYNC_CACHE_KEY_TEMPLATE, channelId)
	isChannelSyncInCache, err := s.cacheStorage.ReadBool(channelSyncKey)
	if err != nil {
		return nil, fmt.Errorf("error reading cache for channel %s: %w", channelId, err)
	}

	if isChannelSyncInCache == nil || !*isChannelSyncInCache {
		s.logger.Debug("Channel sync not found in cache, syncing from Slack", zap.String("channelId", channelId))

		inDatabaseChannel, err = s.slackChannelSynchronizer.Sync(channelId)
		if err != nil {
			return nil, fmt.Errorf("error syncing channel %s: %w", channelId, err)
		}

		s.cacheStorage.Store(channelSyncKey, true, &s.cacheSettings.SlackChannelSynchronizerTTLMinutes)
	}

	return inDatabaseChannel, nil
}
