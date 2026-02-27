package settings

import (
	"time"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const CacheSettingsKey = "cache"

type CacheSettings struct {
	SlackUserSynchronizerTTLMinutes    time.Duration `koanf:"slack_user_synchronizer_ttl_minutes"`
	SlackChannelSynchronizerTTLMinutes time.Duration `koanf:"slack_channel_synchronizer_ttl_minutes"`
}

func NewCacheSettings(logger *zap.Logger, k *koanf.Koanf) *CacheSettings {
	settings := CacheSettings{
		SlackUserSynchronizerTTLMinutes:    60 * time.Minute,
		SlackChannelSynchronizerTTLMinutes: 60 * time.Minute,
	}

	if err := k.Unmarshal(CacheSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal cache settings", zap.Error(err))
	}

	if settings.SlackUserSynchronizerTTLMinutes == 0 {
		logger.Info("Slack User Synchronizer TTL Minutes is set to 0, defaulting to 60 minutes")
		settings.SlackUserSynchronizerTTLMinutes = 60 * time.Minute
	}

	if settings.SlackChannelSynchronizerTTLMinutes == 0 {
		logger.Info("Slack Channel Synchronizer TTL Minutes is set to 0, defaulting to 60 minutes")
		settings.SlackChannelSynchronizerTTLMinutes = 60 * time.Minute
	}

	return &settings
}
