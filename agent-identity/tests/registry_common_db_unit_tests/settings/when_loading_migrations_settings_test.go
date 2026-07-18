package settings_test

import (
	"testing"

	migrations_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingMigrationsSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given a configured folder path that exists", func(t *testing.T) {
		t.Parallel()

		tempDir := t.TempDir()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"migrations.folder_path": tempDir,
		}, "."), nil)

		t.Run("Should load it", func(t *testing.T) {
			t.Parallel()

			settings := migrations_settings.NewMigrationsSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, tempDir, settings.MigrationsFolderPath)
		})
	})

	t.Run("Given a folder path with surrounding whitespace", func(t *testing.T) {
		t.Parallel()

		tempDir := t.TempDir()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"migrations.folder_path": "  " + tempDir + "  ",
		}, "."), nil)

		t.Run("Should trim it", func(t *testing.T) {
			t.Parallel()

			settings := migrations_settings.NewMigrationsSettings(zap.NewNop(), k)

			assert.Equal(t, tempDir, settings.MigrationsFolderPath)
		})
	})
}
