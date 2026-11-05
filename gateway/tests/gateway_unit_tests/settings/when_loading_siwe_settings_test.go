package settings_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingSiweSettings(t *testing.T) {
	t.Parallel()

	consoleSettings := &settings.ConsoleSettings{
		Url:    "https://orchestrator.example.com:8443",
		Domain: "orchestrator.example.com:8443",
	}

	t.Run("Given an empty configuration", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fall back to the defaults", func(t *testing.T) {
			t.Parallel()

			result := settings.NewSiweSettings(zap.NewNop(), koanf.New("."), consoleSettings)

			assert.Equal(t, "Sign in to the 6022 orchestrator.", result.Statement)
			assert.Equal(t, 60, result.SignedRequestMaxAgeSeconds)
		})
	})

	t.Run("Given a non-positive signed request max age", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"siwe.signed_request_max_age_seconds": 0,
		}, "."), nil)

		t.Run("Should fall back to the default", func(t *testing.T) {
			t.Parallel()

			result := settings.NewSiweSettings(zap.NewNop(), k, consoleSettings)

			assert.Equal(t, 60, result.SignedRequestMaxAgeSeconds)
		})
	})

	t.Run("Given a complete configuration", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"siwe.statement":                      "Sign in.",
			"siwe.signed_request_max_age_seconds": 30,
		}, "."), nil)

		t.Run("Should load it and take the uri and domain from the console", func(t *testing.T) {
			t.Parallel()

			result := settings.NewSiweSettings(zap.NewNop(), k, consoleSettings)

			assert.Equal(t, "Sign in.", result.Statement)
			assert.Equal(t, 30, result.SignedRequestMaxAgeSeconds)
			assert.Equal(t, "https://orchestrator.example.com:8443", result.Uri)
			assert.Equal(t, "orchestrator.example.com:8443", result.Domain)
		})
	})
}
