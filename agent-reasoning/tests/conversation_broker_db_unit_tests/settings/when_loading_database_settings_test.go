package settings_test

import (
	"testing"

	db_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_db/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingDatabaseSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given a configured database url", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"database.url": "postgres://user:pass@localhost:5432/broker",
		}, "."), nil)

		t.Run("Should load it", func(t *testing.T) {
			t.Parallel()

			settings := db_settings.NewConversationBrokerDatabaseSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "postgres://user:pass@localhost:5432/broker", settings.Url)
		})
	})
}
