package settings_test

import (
	"testing"

	token_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_blockchain/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingToken6022Settings(t *testing.T) {
	address := "0x0000000000000000000000000000000000000001"

	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"token_6022.smart_contract_address": address,
		}, "."), nil)

		settings := token_settings.NewToken6022Settings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, address, settings.SmartContractAddress.Hex())
	})
}
