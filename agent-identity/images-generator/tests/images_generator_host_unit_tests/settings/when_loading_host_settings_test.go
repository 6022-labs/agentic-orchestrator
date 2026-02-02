package settings_test

import (
	"testing"

	host_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-identity/images-generator/src/images_generator_host/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingImagesGeneratorHostSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("host", map[string]any{
			"api_port":         8001,
			"app_name":         "images-generator",
			"body_limit_bytes": 2048,
		})

		settings := host_settings.NewHostSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, 8001, settings.ApiPort)
		assert.Equal(t, "images-generator", settings.AppName)
		assert.Equal(t, 2048, settings.BodyLimitBytes)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("HOST__API_PORT", "8002")
		t.Setenv("HOST__APP_NAME", "images-generator")
		t.Setenv("HOST__BODY_LIMIT_BYTES", "4096")

		v := viper.New()
		settings := host_settings.NewHostSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, 8002, settings.ApiPort)
		assert.Equal(t, "images-generator", settings.AppName)
		assert.Equal(t, 4096, settings.BodyLimitBytes)
	})
}
