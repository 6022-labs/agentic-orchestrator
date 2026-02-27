package settings_test

import (
	"testing"

	vault_settings "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_vault/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingDatabaseConnectorVaultSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"vault.address":          "http://localhost:8200",
			"vault.token":            "token-123",
			"vault.kv_v1_mount_path": "secret",
		}, "."), nil)

		settings := vault_settings.NewVaultSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "http://localhost:8200", settings.Address)
		assert.Equal(t, "token-123", settings.Token)
		assert.Equal(t, "secret", settings.KvV1MountPath)
	})
}
