package settings

import (
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type CacheSettings struct {
	SlackUserSynchronizerTTLMinutes    time.Duration
	SlackChannelSynchronizerTTLMinutes time.Duration
}

func NewCacheSettings(logger *zap.Logger) *CacheSettings {
	slackUserSynchronizerTTLMinutesStr := os.Getenv("SLACK_USER_SYNCHRONIZER_TTL_MINUTES")
	slackUserSynchronizerTTLMinutes, err := strconv.ParseUint(slackUserSynchronizerTTLMinutesStr, 10, 64)
	if err != nil {
		logger.Info("SLACK_USER_SYNCHRONIZER_TTL_MINUTES is not set, using default value of 60 minutes", zap.Error(err))
		slackUserSynchronizerTTLMinutes = 60
	}

	slackChannelSynchronizerTTLMinutesStr := os.Getenv("SLACK_CHANNEL_SYNCHRONIZER_TTL_MINUTES")
	slackChannelSynchronizerTTLMinutes, err := strconv.ParseUint(slackChannelSynchronizerTTLMinutesStr, 10, 64)
	if err != nil {
		logger.Info("SLACK_CHANNEL_SYNCHRONIZER_TTL_MINUTES is not set, using default value of 60 minutes", zap.Error(err))
		slackChannelSynchronizerTTLMinutes = 60
	}

	return &CacheSettings{
		SlackUserSynchronizerTTLMinutes:    time.Duration(slackUserSynchronizerTTLMinutes) * time.Minute,
		SlackChannelSynchronizerTTLMinutes: time.Duration(slackChannelSynchronizerTTLMinutes) * time.Minute,
	}
}
