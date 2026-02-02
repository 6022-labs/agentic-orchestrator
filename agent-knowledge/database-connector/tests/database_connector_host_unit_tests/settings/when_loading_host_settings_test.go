package settings_test

import (
	"testing"

	host_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_host/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingDatabaseConnectorHostSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("host", map[string]any{
			"api_port":         9001,
			"app_name":         "database-connector",
			"body_limit_bytes": 4096,
		})

		settings := host_settings.NewHostSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, 9001, settings.ApiPort)
		assert.Equal(t, "database-connector", settings.AppName)
		assert.Equal(t, 4096, settings.BodyLimitBytes)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("HOST__API_PORT", "9002")
		t.Setenv("HOST__APP_NAME", "database-connector")
		t.Setenv("HOST__BODY_LIMIT_BYTES", "8192")

		v := viper.New()
		settings := host_settings.NewHostSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, 9002, settings.ApiPort)
		assert.Equal(t, "database-connector", settings.AppName)
		assert.Equal(t, 8192, settings.BodyLimitBytes)
	})
}
