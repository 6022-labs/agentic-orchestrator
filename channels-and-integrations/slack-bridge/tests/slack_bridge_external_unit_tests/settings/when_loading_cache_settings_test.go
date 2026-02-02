package settings_test

import (
	"testing"
	"time"

	cache_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingCacheSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("cache", map[string]any{
			"slack_user_synchronizer_ttl_minutes":    uint64(5),
			"slack_channel_synchronizer_ttl_minutes": uint64(10),
		})

		settings := cache_settings.NewCacheSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, 5*time.Minute, settings.SlackUserSynchronizerTTLMinutes)
		assert.Equal(t, 10*time.Minute, settings.SlackChannelSynchronizerTTLMinutes)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("CACHE__SLACK_USER_SYNCHRONIZER_TTL_MINUTES", "3")
		t.Setenv("CACHE__SLACK_CHANNEL_SYNCHRONIZER_TTL_MINUTES", "7")

		v := viper.New()
		settings := cache_settings.NewCacheSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, 3*time.Minute, settings.SlackUserSynchronizerTTLMinutes)
		assert.Equal(t, 7*time.Minute, settings.SlackChannelSynchronizerTTLMinutes)
	})
}
