package settings_test

import (
	"testing"

	remote_agent_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_remote_agent/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingRemoteAgentSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given configured values", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"remote_agent.request_timeout_seconds":            uint64(60),
			"remote_agent.well_known_request_timeout_seconds": uint64(5),
		}, "."), nil)

		t.Run("Should load them", func(t *testing.T) {
			t.Parallel()

			settings := remote_agent_settings.NewRemoteAgentSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, uint64(60), settings.RequestTimeoutSeconds)
			assert.Equal(t, uint64(5), settings.WellKnownRequestTimeoutSeconds)
		})
	})

	t.Run("Given no configuration", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")

		t.Run("Should fall back to defaults", func(t *testing.T) {
			t.Parallel()

			settings := remote_agent_settings.NewRemoteAgentSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, uint64(120), settings.RequestTimeoutSeconds)
			assert.Equal(t, uint64(10), settings.WellKnownRequestTimeoutSeconds)
		})
	})
}
