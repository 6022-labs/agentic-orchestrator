package settings_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	host_configurations "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_host/configurations"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestWhenLoadingConsoleSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given no configuration", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fail fast", func(t *testing.T) {
			t.Parallel()

			logger := zap.NewNop().WithOptions(zap.WithFatalHook(zapcore.WriteThenPanic))

			assert.Panics(t, func() {
				settings.NewConsoleSettings(logger, koanf.New("."))
			})
		})
	})

	t.Run("Given a relative url", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{"console.url": "orchestrator.example.com"}, "."), nil)

		t.Run("Should fail fast", func(t *testing.T) {
			t.Parallel()

			logger := zap.NewNop().WithOptions(zap.WithFatalHook(zapcore.WriteThenPanic))

			assert.Panics(t, func() {
				settings.NewConsoleSettings(logger, k)
			})
		})
	})

	t.Run("Given a url with a path", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(
			confmap.Provider(map[string]any{"console.url": " https://orchestrator.example.com:8443/console "}, "."),
			nil,
		)

		t.Run("Should derive the origin and the domain", func(t *testing.T) {
			t.Parallel()

			result := settings.NewConsoleSettings(zap.NewNop(), k)

			assert.Equal(t, "https://orchestrator.example.com:8443/console", result.Url)
			assert.Equal(t, "https://orchestrator.example.com:8443", result.Origin)
			assert.Equal(t, "orchestrator.example.com:8443", result.Domain)
		})
	})
}

func TestWhenLoadingConsoleSettingsFromEnv(t *testing.T) {
	t.Run("Given an env override", func(t *testing.T) {
		t.Setenv("CONSOLE__URL", "https://orchestrator.example.com")

		k, err := host_configurations.LoadKoanfConfigFromFile("does-not-exist.json")
		assert.NoError(t, err)

		t.Run("Should win over the default", func(t *testing.T) {
			result := settings.NewConsoleSettings(zap.NewNop(), k)

			assert.Equal(t, "https://orchestrator.example.com", result.Url)
		})
	})
}
