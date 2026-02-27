package settings_test

import (
	"testing"
	"time"

	cache_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingCacheSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"cache.slack_user_synchronizer_ttl_minutes":    "5m",
			"cache.slack_channel_synchronizer_ttl_minutes": "10m",
		}, "."), nil)

		settings := cache_settings.NewCacheSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, 5*time.Minute, settings.SlackUserSynchronizerTTLMinutes)
		assert.Equal(t, 10*time.Minute, settings.SlackChannelSynchronizerTTLMinutes)
	})
}
