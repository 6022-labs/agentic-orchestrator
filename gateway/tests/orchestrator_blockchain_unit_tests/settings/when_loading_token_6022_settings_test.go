package settings_test

import (
	"testing"

	token_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_blockchain/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingToken6022Settings(t *testing.T) {
	address := "0x0000000000000000000000000000000000000001"

	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("token_6022", map[string]any{
			"smart_contract_address": address,
		})

		settings := token_settings.NewToken6022Settings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, address, settings.SmartContractAddress.Hex())
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("TOKEN_6022__SMART_CONTRACT_ADDRESS", address)

		v := viper.New()
		settings := token_settings.NewToken6022Settings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, address, settings.SmartContractAddress.Hex())
	})
}
