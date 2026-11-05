package settings_test

import (
	"testing"

	gateway_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	host_configurations "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_host/configurations"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingGatewaySettings(t *testing.T) {
	t.Parallel()

	t.Run("Given a configured public prefix", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"gateway.public_prefix": " gateway/ ",
		}, "."), nil)

		t.Run("Should normalize it to a leading slash without a trailing one", func(t *testing.T) {
			t.Parallel()

			settings := gateway_settings.NewGatewaySettings(zap.NewNop(), k)

			assert.Equal(t, "/gateway", settings.PublicPrefix)
		})
	})

	t.Run("Given no public prefix", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")

		t.Run("Should fall back to /api", func(t *testing.T) {
			t.Parallel()

			settings := gateway_settings.NewGatewaySettings(zap.NewNop(), k)

			assert.Equal(t, "/api", settings.PublicPrefix)
		})
	})

	t.Run("Given an empty public prefix", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"gateway.public_prefix": "",
		}, "."), nil)

		t.Run("Should mount at the root", func(t *testing.T) {
			t.Parallel()

			settings := gateway_settings.NewGatewaySettings(zap.NewNop(), k)

			assert.Equal(t, "/", settings.PublicPrefix)
		})
	})
}

func TestWhenLoadingGatewaySettingsFromEnv(t *testing.T) {
	t.Run("Given an env override", func(t *testing.T) {
		t.Setenv("GATEWAY__PUBLIC_PREFIX", "/orchestrator")

		k, err := host_configurations.LoadKoanfConfigFromFile("does-not-exist.json")
		assert.NoError(t, err)

		t.Run("Should win over the default", func(t *testing.T) {
			settings := gateway_settings.NewGatewaySettings(zap.NewNop(), k)

			assert.Equal(t, "/orchestrator", settings.PublicPrefix)
		})
	})
}
