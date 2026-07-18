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
	t.Parallel()

	blockchainSettings := &registry_common_settings.BlockchainSettings{ChainId: 80002}

	t.Run("Given a configured start block and block finalization delay", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"sync.start_block":              "12345678",
			"sync.block_finalization_delay": "3s",
		}, "."), nil)

		t.Run("Should load them", func(t *testing.T) {
			t.Parallel()

			settings := sync_settings.NewSyncSettings(zap.NewNop(), k, blockchainSettings)

			assert.NotNil(t, settings)
			assert.Equal(t, uint64(12345678), settings.StartBlock)
			assert.Equal(t, 3*time.Second, settings.BlockFinalizationDelay)
		})
	})

	t.Run("Given no block finalization delay", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"sync.start_block": "12345678",
		}, "."), nil)

		t.Run("Should fall back to the default delay", func(t *testing.T) {
			t.Parallel()

			settings := sync_settings.NewSyncSettings(zap.NewNop(), k, blockchainSettings)

			assert.NotNil(t, settings)
			assert.Equal(t, uint64(12345678), settings.StartBlock)
			assert.Equal(t, 2*time.Second, settings.BlockFinalizationDelay)
		})
	})

	t.Run("Given no configured start block", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")

		t.Run("Should fall back to the built-in start block for the chain", func(t *testing.T) {
			t.Parallel()

			settings := sync_settings.NewSyncSettings(zap.NewNop(), k, blockchainSettings)

			assert.NotNil(t, settings)
			assert.Positive(t, settings.StartBlock)
		})
	})
}
