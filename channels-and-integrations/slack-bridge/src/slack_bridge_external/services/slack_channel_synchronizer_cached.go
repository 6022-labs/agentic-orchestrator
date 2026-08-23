package services

import (
	"context"

	"fmt"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	external_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/settings"
	"go.uber.org/zap"
)

const CHANNEL_SYNC_CACHE_KEY_TEMPLATE = "slack_channel_sync:%s"

const slackChannelCacheName = "slack_channel"

type SlackChannelSynchronizerCached struct {
	channelMetrics metrics_interfaces.ChannelMetricsInterface

	logger                   *zap.Logger
	cacheSettings            *settings.CacheSettings
	cacheStorage             interfaces.CacheStorageInterface
	slackChannelRepository   repositories.SlackChannelRepositoryInterface
	slackChannelSynchronizer external_interfaces.SlackChannelSynchronizerInterface
}

func NewSlackChannelSynchronizerCached(
	channelMetrics metrics_interfaces.ChannelMetricsInterface,
	logger *zap.Logger,
	cacheSettings *settings.CacheSettings,
	cacheStorage interfaces.CacheStorageInterface,
	slackChannelRepository repositories.SlackChannelRepositoryInterface,
	slackChannelSynchronizer external_interfaces.SlackChannelSynchronizerInterface,
) *SlackChannelSynchronizerCached {
	return &SlackChannelSynchronizerCached{
		channelMetrics: channelMetrics,

		logger:                   logger,
		cacheStorage:             cacheStorage,
		cacheSettings:            cacheSettings,
		slackChannelRepository:   slackChannelRepository,
		slackChannelSynchronizer: slackChannelSynchronizer,
	}
}

func (s *SlackChannelSynchronizerCached) Sync(ctx context.Context, channelId string) (*models.SlackChannel, error) {
	inDatabaseChannel, err := s.slackChannelRepository.GetById(ctx, channelId)
	if err != nil {
		return nil, fmt.Errorf("error getting channel %s from database: %w", channelId, err)
	}

	channelSyncKey := fmt.Sprintf(CHANNEL_SYNC_CACHE_KEY_TEMPLATE, channelId)
	isChannelSyncInCache, err := s.cacheStorage.ReadBool(ctx, channelSyncKey)
	if err != nil {
		return nil, fmt.Errorf("error reading cache for channel %s: %w", channelId, err)
	}

	s.channelMetrics.RecordCacheResult(ctx, slackChannelCacheName, isChannelSyncInCache != nil && *isChannelSyncInCache)

	if isChannelSyncInCache == nil || !*isChannelSyncInCache {
		s.logger.Debug("Channel sync not found in cache, syncing from Slack", zap.String("channelId", channelId))

		inDatabaseChannel, err = s.slackChannelSynchronizer.Sync(ctx, channelId)
		if err != nil {
			return nil, fmt.Errorf("error syncing channel %s: %w", channelId, err)
		}

		s.cacheStorage.Store(ctx, channelSyncKey, true, &s.cacheSettings.SlackChannelSynchronizerTTLMinutes)
	}

	return inDatabaseChannel, nil
}
