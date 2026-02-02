package settings_test

import (
	"os"
	"testing"

	image_generation_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-identity/images-generator/src/images_generator/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentImageGenerationSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "image-mark-*")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())
		_ = tempFile.Close()

		v := viper.New()
		v.Set("image_generation", map[string]any{
			"image_mark_path": tempFile.Name(),
		})

		settings := image_generation_settings.NewImageGenerationSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		if assert.NotNil(t, settings.ImageMarkPath) {
			assert.Equal(t, tempFile.Name(), *settings.ImageMarkPath)
		}
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "image-mark-*")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())
		_ = tempFile.Close()

		t.Setenv("IMAGE_GENERATION__IMAGE_MARK_PATH", tempFile.Name())

		v := viper.New()
		settings := image_generation_settings.NewImageGenerationSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		if assert.NotNil(t, settings.ImageMarkPath) {
			assert.Equal(t, tempFile.Name(), *settings.ImageMarkPath)
		}
	})
}
