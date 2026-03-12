package settings_test

import (
	"testing"

	sync_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingSyncSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"sync.start_block": "12345678",
		}, "."), nil)

		settings := sync_settings.NewSyncSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, uint64(12345678), settings.StartBlock)
	})
}
