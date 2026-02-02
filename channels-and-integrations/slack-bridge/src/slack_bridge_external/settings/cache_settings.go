package settings

import (
	"os"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const CacheSettingsKey = "cache"

type CacheSettings struct {
	SlackUserSynchronizerTTLMinutes    time.Duration
	SlackChannelSynchronizerTTLMinutes time.Duration
}

func NewCacheSettings(logger *zap.Logger, v *viper.Viper) *CacheSettings {
	v.SetDefault(CacheSettingsKey+".slack_user_synchronizer_ttl_minutes", uint64(60))
	v.SetDefault(CacheSettingsKey+".slack_channel_synchronizer_ttl_minutes", uint64(60))

	if val, ok := os.LookupEnv("CACHE__SLACK_USER_SYNCHRONIZER_TTL_MINUTES"); ok {
		v.Set(CacheSettingsKey+".slack_user_synchronizer_ttl_minutes", val)
	}
	if val, ok := os.LookupEnv("CACHE__SLACK_CHANNEL_SYNCHRONIZER_TTL_MINUTES"); ok {
		v.Set(CacheSettingsKey+".slack_channel_synchronizer_ttl_minutes", val)
	}

	slackUserSynchronizerTTLMinutes := v.GetUint64(CacheSettingsKey + ".slack_user_synchronizer_ttl_minutes")
	if slackUserSynchronizerTTLMinutes == 0 {
		logger.Info("Slack User Synchronizer TTL Minutes is set to 0, defaulting to 60 minutes")
		slackUserSynchronizerTTLMinutes = 60
	}

	slackChannelSynchronizerTTLMinutes := v.GetUint64(CacheSettingsKey + ".slack_channel_synchronizer_ttl_minutes")
	if slackChannelSynchronizerTTLMinutes == 0 {
		logger.Info("Slack Channel Synchronizer TTL Minutes is set to 0, defaulting to 60 minutes")
		slackChannelSynchronizerTTLMinutes = 60
	}

	settings := CacheSettings{
		SlackUserSynchronizerTTLMinutes:    time.Duration(slackUserSynchronizerTTLMinutes) * time.Minute,
		SlackChannelSynchronizerTTLMinutes: time.Duration(slackChannelSynchronizerTTLMinutes) * time.Minute,
	}

	return &settings
}
