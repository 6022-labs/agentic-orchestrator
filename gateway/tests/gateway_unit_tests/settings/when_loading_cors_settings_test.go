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

func TestWhenLoadingCorsSettings(t *testing.T) {
	t.Parallel()

	consoleSettings := &settings.ConsoleSettings{Origin: "https://console.example.com"}

	t.Run("Given no configuration", func(t *testing.T) {
		t.Parallel()

		t.Run("Should allow only the console origin", func(t *testing.T) {
			t.Parallel()

			result := settings.NewCorsSettings(zap.NewNop(), koanf.New("."), consoleSettings)

			assert.Equal(t, []string{"https://console.example.com"}, result.AllowOrigins)
		})
	})

	t.Run("Given an origin with a path", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"cors.allow_origins": []string{"https://app.example.com/console"},
		}, "."), nil)

		t.Run("Should fail fast", func(t *testing.T) {
			t.Parallel()

			logger := zap.NewNop().WithOptions(zap.WithFatalHook(zapcore.WriteThenPanic))

			assert.Panics(t, func() {
				settings.NewCorsSettings(logger, k, consoleSettings)
			})
		})
	})

	t.Run("Given a wildcard", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{"cors.allow_origins": []string{"*"}}, "."), nil)

		t.Run("Should keep it", func(t *testing.T) {
			t.Parallel()

			result := settings.NewCorsSettings(zap.NewNop(), k, consoleSettings)

			assert.Equal(t, []string{"*"}, result.AllowOrigins)
		})
	})

	t.Run("Given a list of origins", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"cors.allow_origins": []string{"https://a.example.com/", " https://b.example.com:8443"},
		}, "."), nil)

		t.Run("Should trim them and drop trailing slashes", func(t *testing.T) {
			t.Parallel()

			result := settings.NewCorsSettings(zap.NewNop(), k, consoleSettings)

			assert.Equal(t, []string{"https://a.example.com", "https://b.example.com:8443"}, result.AllowOrigins)
		})
	})
}

func TestWhenLoadingCorsSettingsFromEnv(t *testing.T) {
	t.Run("Given a comma-separated env override", func(t *testing.T) {
		t.Setenv("CORS__ALLOW_ORIGINS", "https://a.example.com, https://b.example.com")

		k, err := host_configurations.LoadKoanfConfigFromFile("does-not-exist.json")
		assert.NoError(t, err)

		t.Run("Should split it into origins", func(t *testing.T) {
			result := settings.NewCorsSettings(zap.NewNop(), k, &settings.ConsoleSettings{})

			assert.Equal(t, []string{"https://a.example.com", "https://b.example.com"}, result.AllowOrigins)
		})
	})
}
