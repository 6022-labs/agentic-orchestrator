package settings_test

import (
	"testing"

	host_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/images_generator_host/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingImagesGeneratorHostSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"host.api_port":         8001,
			"host.app_name":         "images-generator",
			"host.body_limit_bytes": 2048,
		}, "."), nil)

		settings := host_settings.NewHostSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, 8001, settings.ApiPort)
		assert.Equal(t, "images-generator", settings.AppName)
		assert.Equal(t, 2048, settings.BodyLimitBytes)
	})
}
