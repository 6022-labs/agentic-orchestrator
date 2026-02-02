package settings_test

import (
	"testing"

	images_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_agent_identity_images_generator/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingImagesGeneratorSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("images_generator", map[string]any{
			"base_url": "https://images.example.com/",
		})

		settings := images_settings.NewImagesGeneratorSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://images.example.com", settings.BaseUrl)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("IMAGES_GENERATOR__BASE_URL", "https://images.env/")

		v := viper.New()
		settings := images_settings.NewImagesGeneratorSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://images.env", settings.BaseUrl)
	})
}
