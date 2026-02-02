package settings_test

import (
	"testing"

	db_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_agent_knowledge_database_connector/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentDatabaseConnectorSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("agent_database_connector", map[string]any{
			"base_url": "https://db.example.com/",
		})

		settings := db_settings.NewAgentDatabaseConnectorSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://db.example.com", settings.BaseUrl)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("AGENT_DATABASE_CONNECTOR__BASE_URL", "https://db.env/")

		v := viper.New()
		settings := db_settings.NewAgentDatabaseConnectorSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://db.env", settings.BaseUrl)
	})
}
