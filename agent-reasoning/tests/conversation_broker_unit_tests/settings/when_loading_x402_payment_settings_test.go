package settings_test

import (
	"testing"

	x402_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingX402PaymentSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given configured values", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"x402_payment.signature_lifetime_seconds":    uint64(120),
			"x402_payment.clock_drift_tolerance_seconds": uint64(30),
			"x402_payment.challenge_cache_ttl_seconds":   uint64(90),
		}, "."), nil)

		t.Run("Should load them", func(t *testing.T) {
			t.Parallel()

			settings := x402_settings.NewX402PaymentSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, uint64(120), settings.SignatureLifetimeSeconds)
			assert.Equal(t, uint64(30), settings.ClockDriftToleranceSeconds)
			assert.Equal(t, uint64(90), settings.ChallengeCacheTtlSeconds)
		})
	})

	t.Run("Given no configuration", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")

		t.Run("Should fall back to defaults", func(t *testing.T) {
			t.Parallel()

			settings := x402_settings.NewX402PaymentSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, uint64(300), settings.SignatureLifetimeSeconds)
			assert.Equal(t, uint64(60), settings.ClockDriftToleranceSeconds)
			assert.Equal(t, uint64(300), settings.ChallengeCacheTtlSeconds)
		})
	})
}
