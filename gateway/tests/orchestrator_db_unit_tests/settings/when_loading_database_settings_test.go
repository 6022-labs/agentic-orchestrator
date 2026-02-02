package settings_test

import (
	"testing"

	db_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_db/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingDatabaseSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("database", map[string]any{
			"url": "postgres://user:pass@localhost:5432/db",
		})

		settings := db_settings.NewDatabaseSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "postgres://user:pass@localhost:5432/db", settings.Url)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("DATABASE__URL", "postgres://env:env@localhost:5432/db")

		v := viper.New()
		settings := db_settings.NewDatabaseSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "postgres://env:env@localhost:5432/db", settings.Url)
	})
}
