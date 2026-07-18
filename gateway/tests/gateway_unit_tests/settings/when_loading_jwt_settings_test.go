package settings_test

import (
	"testing"

	jwt_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingJwtSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given configured values", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"jwt.hash":        "test-signing-hash-at-least-32-characters",
			"jwt.living_time": 3600,
			"jwt.nonce_ttl":   600,
		}, "."), nil)

		t.Run("Should load them", func(t *testing.T) {
			t.Parallel()

			settings := jwt_settings.NewJwtSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "test-signing-hash-at-least-32-characters", settings.Hash)
			assert.Equal(t, 3600, settings.LivingTime)
			assert.Equal(t, 600, settings.NonceTtl)
		})
	})

	t.Run("Given no nonce ttl", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"jwt.hash":        "test-signing-hash-at-least-32-characters",
			"jwt.living_time": 3600,
		}, "."), nil)

		t.Run("Should fall back to 300", func(t *testing.T) {
			t.Parallel()

			settings := jwt_settings.NewJwtSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "test-signing-hash-at-least-32-characters", settings.Hash)
			assert.Equal(t, 3600, settings.LivingTime)
			assert.Equal(t, 300, settings.NonceTtl)
		})
	})

	t.Run("Given a non-positive nonce ttl", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"jwt.hash":        "test-signing-hash-at-least-32-characters",
			"jwt.living_time": 3600,
			"jwt.nonce_ttl":   0,
		}, "."), nil)

		t.Run("Should fall back to 300", func(t *testing.T) {
			t.Parallel()

			settings := jwt_settings.NewJwtSettings(zap.NewNop(), k)

			assert.Equal(t, 300, settings.NonceTtl)
		})
	})
}
