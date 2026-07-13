package settings_test

import (
	"testing"
	"time"

	registry_common_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	sync_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingSyncSettings(t *testing.T) {
	t.Run("Loads settings with block finalization delay", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"sync.start_block":              "12345678",
			"sync.block_finalization_delay": "3s",
		}, "."), nil)

		settings := sync_settings.NewSyncSettings(zap.NewNop(), k, &registry_common_settings.BlockchainSettings{ChainId: 80002})

		assert.NotNil(t, settings)
		assert.Equal(t, uint64(12345678), settings.StartBlock)
		assert.Equal(t, 3*time.Second, settings.BlockFinalizationDelay)
	})

	t.Run("Defaults the block finalization delay", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"sync.start_block": "12345678",
		}, "."), nil)

		settings := sync_settings.NewSyncSettings(zap.NewNop(), k, &registry_common_settings.BlockchainSettings{ChainId: 80002})

		assert.NotNil(t, settings)
		assert.Equal(t, uint64(12345678), settings.StartBlock)
		assert.Equal(t, 2*time.Second, settings.BlockFinalizationDelay)
	})

	t.Run("Falls back to the built-in start block when none is configured", func(t *testing.T) {
		k := koanf.New(".")

		settings := sync_settings.NewSyncSettings(zap.NewNop(), k, &registry_common_settings.BlockchainSettings{ChainId: 80002})

		assert.NotNil(t, settings)
		assert.Equal(t, uint64(39471293), settings.StartBlock)
	})
}
