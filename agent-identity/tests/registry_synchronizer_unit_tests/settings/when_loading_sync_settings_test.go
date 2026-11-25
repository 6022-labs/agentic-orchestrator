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
	"go.uber.org/zap/zapcore"
)

func TestWhenLoadingSyncSettings(t *testing.T) {
	t.Parallel()

	blockchainSettings := &registry_common_settings.BlockchainSettings{ChainId: 80002}

	t.Run("Given a configured start block and block finalization delay", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"sync.start_block":                "12345678",
			"sync.block_finalization_delay":   "3s",
			"sync.block_events_poll_interval": "50ms",
		}, "."), nil)

		t.Run("Should load them", func(t *testing.T) {
			t.Parallel()

			settings := sync_settings.NewSyncSettings(zap.NewNop(), k, blockchainSettings)

			assert.NotNil(t, settings)
			assert.Equal(t, uint64(12345678), settings.StartBlock)
			assert.Equal(t, 3*time.Second, settings.BlockFinalizationDelay)
			assert.Equal(t, 50*time.Millisecond, settings.BlockEventsPollInterval)
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
			assert.Equal(t, 200*time.Millisecond, settings.BlockEventsPollInterval)
		})
	})

	t.Run("Given a configured gap window", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"sync.start_block":       "12345678",
			"sync.gap_window_blocks": "2000",
		}, "."), nil)

		t.Run("Should load it", func(t *testing.T) {
			t.Parallel()

			settings := sync_settings.NewSyncSettings(zap.NewNop(), k, blockchainSettings)

			assert.Equal(t, uint64(2000), settings.GapWindowBlocks)
		})
	})

	t.Run("Given a non-positive gap window", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"sync.start_block":       "12345678",
			"sync.gap_window_blocks": "0",
		}, "."), nil)

		t.Run("Should exit", func(t *testing.T) {
			t.Parallel()

			logger := zap.NewNop().WithOptions(zap.WithFatalHook(zapcore.WriteThenPanic))

			assert.Panics(t, func() { sync_settings.NewSyncSettings(logger, k, blockchainSettings) })
		})
	})

	t.Run("Given no gap window", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"sync.start_block": "12345678",
		}, "."), nil)

		t.Run("Should fall back to the default window", func(t *testing.T) {
			t.Parallel()

			settings := sync_settings.NewSyncSettings(zap.NewNop(), k, blockchainSettings)

			assert.Equal(t, uint64(10_000), settings.GapWindowBlocks)
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
