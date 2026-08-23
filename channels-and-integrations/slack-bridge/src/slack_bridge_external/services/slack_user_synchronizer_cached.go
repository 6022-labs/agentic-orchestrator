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

const USER_SYNC_CACHE_KEY_TEMPLATE = "slack_user_sync:%s"

const slackUserCacheName = "slack_user"

type SlackUserSynchronizerCached struct {
	channelMetrics metrics_interfaces.ChannelMetricsInterface

	logger                *zap.Logger
	cacheSettings         *settings.CacheSettings
	cacheStorage          interfaces.CacheStorageInterface
	slackUserRepository   repositories.SlackUserRepositoryInterface
	slackUserSynchronizer external_interfaces.SlackUserSynchronizerInterface
}

func NewSlackUserSynchronizerCached(
	channelMetrics metrics_interfaces.ChannelMetricsInterface,
	logger *zap.Logger,
	cacheSettings *settings.CacheSettings,
	cacheStorage interfaces.CacheStorageInterface,
	slackUserRepository repositories.SlackUserRepositoryInterface,
	slackUserSynchronizer external_interfaces.SlackUserSynchronizerInterface,
) *SlackUserSynchronizerCached {
	return &SlackUserSynchronizerCached{
		channelMetrics: channelMetrics,

		logger:                logger,
		cacheStorage:          cacheStorage,
		cacheSettings:         cacheSettings,
		slackUserRepository:   slackUserRepository,
		slackUserSynchronizer: slackUserSynchronizer,
	}
}

func (s *SlackUserSynchronizerCached) Sync(ctx context.Context, id string) (*models.SlackUser, error) {
	inDatabaseUser, err := s.slackUserRepository.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting user %s from database: %w", id, err)
	}

	userSyncKey := fmt.Sprintf(USER_SYNC_CACHE_KEY_TEMPLATE, id)
	isUserSyncInCache, err := s.cacheStorage.ReadBool(ctx, userSyncKey)
	if err != nil {
		return nil, fmt.Errorf("error reading cache for user %s: %w", id, err)
	}

	s.channelMetrics.RecordCacheResult(ctx, slackUserCacheName,
		isUserSyncInCache != nil && *isUserSyncInCache && inDatabaseUser != nil)

	if isUserSyncInCache == nil || !*isUserSyncInCache || inDatabaseUser == nil {
		s.logger.Debug("User sync not found in cache, syncing from Slack", zap.String("id", id))

		inDatabaseUser, err = s.slackUserSynchronizer.Sync(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("error syncing user %s: %w", id, err)
		}

		s.cacheStorage.Store(ctx, userSyncKey, true, &s.cacheSettings.SlackUserSynchronizerTTLMinutes)
	}

	return inDatabaseUser, nil
}
