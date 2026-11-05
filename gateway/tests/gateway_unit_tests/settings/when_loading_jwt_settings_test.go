package settings_test

import (
	"testing"

	jwt_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	host_configurations "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_host/configurations"
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
			"jwt.hash":             "test-signing-hash-at-least-32-characters",
			"jwt.living_time":      3600,
			"jwt.nonce_ttl":        600,
			"jwt.refresh_ttl":      86400,
			"jwt.refresh_idle_ttl": 7200,
		}, "."), nil)

		t.Run("Should load them", func(t *testing.T) {
			t.Parallel()

			settings := jwt_settings.NewJwtSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "test-signing-hash-at-least-32-characters", settings.Hash)
			assert.Equal(t, 3600, settings.LivingTime)
			assert.Equal(t, 600, settings.NonceTtl)
			assert.Equal(t, 86400, settings.RefreshTtl)
			assert.Equal(t, 7200, settings.RefreshIdleTtl)
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
	t.Run("Given no refresh settings", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"jwt.hash":        "test-signing-hash-at-least-32-characters",
			"jwt.living_time": 3600,
		}, "."), nil)

		t.Run("Should fall back to a 7 day session with a 1 day idle timeout", func(t *testing.T) {
			t.Parallel()

			settings := jwt_settings.NewJwtSettings(zap.NewNop(), k)

			assert.Equal(t, 604800, settings.RefreshTtl)
			assert.Equal(t, 86400, settings.RefreshIdleTtl)
		})
	})

	t.Run("Given a non-positive refresh ttl", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"jwt.hash":        "test-signing-hash-at-least-32-characters",
			"jwt.living_time": 3600,
			"jwt.refresh_ttl": 0,
		}, "."), nil)

		t.Run("Should fall back to 604800", func(t *testing.T) {
			t.Parallel()

			settings := jwt_settings.NewJwtSettings(zap.NewNop(), k)

			assert.Equal(t, 604800, settings.RefreshTtl)
		})
	})

	t.Run("Given a non-positive refresh idle ttl", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"jwt.hash":             "test-signing-hash-at-least-32-characters",
			"jwt.living_time":      3600,
			"jwt.refresh_idle_ttl": 0,
		}, "."), nil)

		t.Run("Should fall back to 86400", func(t *testing.T) {
			t.Parallel()

			settings := jwt_settings.NewJwtSettings(zap.NewNop(), k)

			assert.Equal(t, 86400, settings.RefreshIdleTtl)
		})
	})

	t.Run("Given a refresh idle ttl under twice the living time", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"jwt.hash":             "test-signing-hash-at-least-32-characters",
			"jwt.living_time":      86400,
			"jwt.refresh_idle_ttl": 86400,
		}, "."), nil)

		t.Run("Should raise it to twice the living time", func(t *testing.T) {
			t.Parallel()

			settings := jwt_settings.NewJwtSettings(zap.NewNop(), k)

			assert.Equal(t, 172800, settings.RefreshIdleTtl)
		})
	})
}

func TestWhenLoadingJwtSettingsFromEnv(t *testing.T) {
	t.Run("Given an env override", func(t *testing.T) {
		t.Setenv("JWT__HASH", "test-signing-hash-at-least-32-characters")
		t.Setenv("JWT__LIVING_TIME", "3600")
		t.Setenv("JWT__REFRESH_TTL", "86400")
		t.Setenv("JWT__REFRESH_IDLE_TTL", "7200")

		k, err := host_configurations.LoadKoanfConfigFromFile("does-not-exist.json")
		assert.NoError(t, err)

		t.Run("Should win over the defaults", func(t *testing.T) {
			settings := jwt_settings.NewJwtSettings(zap.NewNop(), k)

			assert.Equal(t, 86400, settings.RefreshTtl)
			assert.Equal(t, 7200, settings.RefreshIdleTtl)
		})
	})
}
