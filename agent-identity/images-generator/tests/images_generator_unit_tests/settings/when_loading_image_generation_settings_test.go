package settings_test

import (
	"os"
	"testing"

	image_generation_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/images-generator/src/images_generator/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentImageGenerationSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "image-mark-*")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())
		_ = tempFile.Close()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"image_generation.image_mark_path": tempFile.Name(),
		}, "."), nil)

		settings := image_generation_settings.NewImageGenerationSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		if assert.NotNil(t, settings.ImageMarkPath) {
			assert.Equal(t, tempFile.Name(), *settings.ImageMarkPath)
		}
	})
}
