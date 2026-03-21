package settings_test

import (
	"testing"

	orchestrator_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingOrchestratorSettings(t *testing.T) {
	address := "0x0000000000000000000000000000000000000001"

	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"orchestrator.public_key": address,
		}, "."), nil)

		settings := orchestrator_settings.NewOrchestratorSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.True(t, common.HexToAddress(address) == settings.PublicKey)
	})
}
