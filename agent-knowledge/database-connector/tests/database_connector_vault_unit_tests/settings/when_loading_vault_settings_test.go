package settings_test

import (
	"testing"

	vault_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_vault/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingDatabaseConnectorVaultSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("vault", map[string]any{
			"address":          "http://localhost:8200",
			"token":            "token-123",
			"kv_v1_mount_path": "secret",
		})

		settings := vault_settings.NewVaultSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "http://localhost:8200", settings.Address)
		assert.Equal(t, "token-123", settings.Token)
		assert.Equal(t, "secret", settings.KvV1MountPath)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("VAULT__ADDRESS", "http://localhost:8200")
		t.Setenv("VAULT__TOKEN", "token-456")
		t.Setenv("VAULT__KV_V1_MOUNT_PATH", "secret-env")

		v := viper.New()
		settings := vault_settings.NewVaultSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "http://localhost:8200", settings.Address)
		assert.Equal(t, "token-456", settings.Token)
		assert.Equal(t, "secret-env", settings.KvV1MountPath)
	})
}
