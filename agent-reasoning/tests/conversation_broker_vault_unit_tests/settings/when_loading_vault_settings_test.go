package settings_test

import (
	"testing"

	vault_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_vault/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingVaultSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given configured values", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"vault.address":          "http://127.0.0.1:8200",
			"vault.token":            "token",
			"vault.kv_v1_mount_path": "kv",
		}, "."), nil)

		t.Run("Should load them", func(t *testing.T) {
			t.Parallel()

			settings := vault_settings.NewVaultSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "http://127.0.0.1:8200", settings.Address)
			assert.Equal(t, "token", settings.Token)
			assert.Equal(t, "kv", settings.KvV1MountPath)
		})
	})

	t.Run("Given no mount path", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"vault.address": "http://127.0.0.1:8200",
			"vault.token":   "token",
		}, "."), nil)

		t.Run("Should default it to secret", func(t *testing.T) {
			t.Parallel()

			settings := vault_settings.NewVaultSettings(zap.NewNop(), k)

			assert.Equal(t, "secret", settings.KvV1MountPath)
		})
	})
}
