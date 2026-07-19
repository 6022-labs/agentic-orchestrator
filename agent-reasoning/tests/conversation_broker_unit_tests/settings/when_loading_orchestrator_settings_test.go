package settings_test

import (
	"testing"

	orchestrator_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingOrchestratorSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given a configured public key", func(t *testing.T) {
		t.Parallel()

		address := "0x0000000000000000000000000000000000000001"
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"orchestrator.public_key": address,
		}, "."), nil)

		t.Run("Should load it as an address", func(t *testing.T) {
			t.Parallel()

			settings := orchestrator_settings.NewOrchestratorSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, common.HexToAddress(address), settings.PublicKey)
			assert.True(t, settings.HasPublicKey())
		})
	})

	t.Run("Given no configuration", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")

		t.Run("Should report no public key", func(t *testing.T) {
			t.Parallel()

			settings := orchestrator_settings.NewOrchestratorSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.False(t, settings.HasPublicKey())
		})
	})
}
