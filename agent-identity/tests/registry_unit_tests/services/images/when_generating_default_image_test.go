package images_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services/images"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/settings"
	"github.com/stretchr/testify/assert"
)

func TestWhenGeneratingDefaultImage(t *testing.T) {
	t.Parallel()

	t.Run("Given an agent with a base image, wallets and stars", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an animated SVG embedding the image, gradient and stars", func(t *testing.T) {
			t.Parallel()

			imageGenerationSettings := &settings.ImageGenerationSettings{ImageMarkPath: nil}
			sut := images.NewDefaultImageGenerator(imageGenerationSettings)

			result, err := sut.Generate(sampleImageParams(t))

			assert.NoError(t, err)
			svg := string(result)
			assert.Contains(t, svg, "<svg")
			assert.Contains(t, svg, "data:image/png;base64,")
			assert.Contains(t, svg, "url(#borderGradient)")
			assert.Contains(t, svg, "★★★")
			assert.Contains(t, svg, "<animate")
		})
	})
}
