package settings_test

import (
	"testing"

	host_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_host/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingLlmConnectorHostSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("host", map[string]any{
			"api_port":         9101,
			"app_name":         "llm-connector",
			"body_limit_bytes": 1024,
		})

		settings := host_settings.NewHostSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, 9101, settings.ApiPort)
		assert.Equal(t, "llm-connector", settings.AppName)
		assert.Equal(t, 1024, settings.BodyLimitBytes)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("HOST__API_PORT", "9102")
		t.Setenv("HOST__APP_NAME", "llm-connector")
		t.Setenv("HOST__BODY_LIMIT_BYTES", "2048")

		v := viper.New()
		settings := host_settings.NewHostSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, 9102, settings.ApiPort)
		assert.Equal(t, "llm-connector", settings.AppName)
		assert.Equal(t, 2048, settings.BodyLimitBytes)
	})
}
