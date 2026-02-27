package settings_test

import (
	"testing"

	migrations_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingSlackBridgeMigrationsSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		folder := t.TempDir()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"migrations.folder_path": folder,
		}, "."), nil)

		settings := migrations_settings.NewMigrationsSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, folder, settings.MigrationsFolderPath)
	})
}
