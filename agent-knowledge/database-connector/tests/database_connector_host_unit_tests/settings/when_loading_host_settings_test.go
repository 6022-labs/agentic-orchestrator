package settings_test

import (
	"testing"

	host_settings "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_host/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingDatabaseConnectorHostSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"host.api_port":         9001,
			"host.app_name":         "database-connector",
			"host.body_limit_bytes": 4096,
		}, "."), nil)

		settings := host_settings.NewHostSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, 9001, settings.ApiPort)
		assert.Equal(t, "database-connector", settings.AppName)
		assert.Equal(t, 4096, settings.BodyLimitBytes)
	})
}
