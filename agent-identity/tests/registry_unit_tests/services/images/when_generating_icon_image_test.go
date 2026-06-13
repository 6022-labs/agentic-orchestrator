package images_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services/images"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/settings"
	"github.com/stretchr/testify/assert"
)

// pngMagic is the 8-byte PNG signature.
var pngMagic = []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}

func TestWhenGeneratingIconImage(t *testing.T) {
	t.Parallel()

	t.Run("Given an agent frame to render", func(t *testing.T) {
		t.Parallel()

		t.Run("Should rasterize the static SVG into a PNG", func(t *testing.T) {
			t.Parallel()

			imageGenerationSettings := &settings.ImageGenerationSettings{ImageMarkPath: nil}
			sut := images.NewIconImageGenerator(imageGenerationSettings)

			result, err := sut.Generate(sampleImageParams(t))

			assert.NoError(t, err)
			assert.GreaterOrEqual(t, len(result), len(pngMagic))
			assert.Equal(t, pngMagic, result[:len(pngMagic)])
		})
	})
}
