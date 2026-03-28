package settings_test

import (
	"testing"

	images_generator_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_image_generator/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingImagesGeneratorSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"images_generator.base_url": "https://images.example.com/",
		}, "."), nil)

		settings := images_generator_settings.NewImagesGeneratorSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://images.example.com", settings.BaseUrl)
	})

	t.Run("Trims trailing slash from base url", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"images_generator.base_url": "https://images.example.com///",
		}, "."), nil)

		settings := images_generator_settings.NewImagesGeneratorSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://images.example.com", settings.BaseUrl)
	})
}
