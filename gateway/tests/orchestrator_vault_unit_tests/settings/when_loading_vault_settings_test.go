package settings_test

import (
	"testing"

	vault_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_vault/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingVaultSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"vault.address":          "http://127.0.0.1:8200",
			"vault.token":            "token",
			"vault.kv_v1_mount_path": "secret",
		}, "."), nil)

		settings := vault_settings.NewVaultSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "http://127.0.0.1:8200", settings.Address)
		assert.Equal(t, "token", settings.Token)
		assert.Equal(t, "secret", settings.KvV1MountPath)
	})
}
