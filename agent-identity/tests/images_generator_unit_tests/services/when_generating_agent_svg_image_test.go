package services_test

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/images_generator/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/images_generator/settings"
	"github.com/stretchr/testify/assert"
)

type WhenGeneratingAgentSVGImageTestingSuite struct {
	sut *services.AgentSVGImage
}

func WhenGeneratingAgentSVGImageBeforeEach(t *testing.T) *WhenGeneratingAgentSVGImageTestingSuite {
	imageGenSettings := &settings.ImageGenerationSettings{
		ImageMarkPath: nil,
	}

	sut := services.NewAgentSVGImage(imageGenSettings)

	return &WhenGeneratingAgentSVGImageTestingSuite{
		sut: sut,
	}
}

func TestWhenGeneratingAgentSVGImage(t *testing.T) {
	t.Parallel()

	t.Run("Given RenderBorderGradient is called", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an SVG linear gradient element", func(t *testing.T) {
			t.Parallel()
			basics := WhenGeneratingAgentSVGImageBeforeEach(t)

			result := basics.sut.RenderBorderGradient("0x1234abcd")

			assert.Contains(t, result, `<linearGradient id="borderGradient"`)
			assert.Contains(t, result, `</linearGradient>`)
			assert.Contains(t, result, `<stop offset=`)
		})

		t.Run("Should produce deterministic output for the same creator", func(t *testing.T) {
			t.Parallel()
			basics := WhenGeneratingAgentSVGImageBeforeEach(t)

			result1 := basics.sut.RenderBorderGradient("0xaabbcc")
			result2 := basics.sut.RenderBorderGradient("0xaabbcc")

			assert.Equal(t, result1, result2)
		})

		t.Run("Should produce different output for different creators", func(t *testing.T) {
			t.Parallel()
			basics := WhenGeneratingAgentSVGImageBeforeEach(t)

			result1 := basics.sut.RenderBorderGradient("0xaaa")
			result2 := basics.sut.RenderBorderGradient("0xbbb")

			assert.NotEqual(t, result1, result2)
		})
	})

	t.Run("Given RenderImageData is called", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a base64 data URI", func(t *testing.T) {
			t.Parallel()
			basics := WhenGeneratingAgentSVGImageBeforeEach(t)

			imageBytes := []byte("test-image-data")
			mimeType := "image/png"

			result := basics.sut.RenderImageData(imageBytes, mimeType)

			expectedB64 := base64.StdEncoding.EncodeToString(imageBytes)
			expected := fmt.Sprintf("data:%s;base64,%s", mimeType, expectedB64)
			assert.Equal(t, expected, result)
		})
	})

	t.Run("Given RenderAgentInformation is called", func(t *testing.T) {
		t.Parallel()

		t.Run("Should render badges for role, created, and creator", func(t *testing.T) {
			t.Parallel()
			basics := WhenGeneratingAgentSVGImageBeforeEach(t)

			params := services.GenerateImageParams{
				Name:      "TestAgent",
				EnsDomain: "test.6022.eth",
				Role:      "expert",
				Creator:   "0x1234",
				CreatedAt: time.Unix(1000000, 0),
			}

			result := basics.sut.RenderAgentInformation(params)

			assert.Contains(t, result, "Role: expert")
			assert.Contains(t, result, "Created:")
			assert.Contains(t, result, "Creator: 0x1234")
		})

		t.Run("Should render origin badge when OriginEnsDomain is set", func(t *testing.T) {
			t.Parallel()
			basics := WhenGeneratingAgentSVGImageBeforeEach(t)

			originDomain := "origin.6022.eth"
			originGen := uint(2)
			params := services.GenerateImageParams{
				Name:             "CloneAgent",
				EnsDomain:        "clone.6022.eth",
				Role:             "human",
				Creator:          "0x5678",
				OriginEnsDomain:  &originDomain,
				OriginGeneration: &originGen,
				CreatedAt:        time.Unix(1000000, 0),
			}

			result := basics.sut.RenderAgentInformation(params)

			assert.Contains(t, result, "Origin:")
			assert.Contains(t, result, "origin.6022.eth")
		})
	})

	t.Run("Given RenderStars is called", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return empty string for zero stars", func(t *testing.T) {
			t.Parallel()
			basics := WhenGeneratingAgentSVGImageBeforeEach(t)

			result := basics.sut.RenderStars(0)

			assert.Empty(t, result)
		})

		t.Run("Should render the correct number of stars", func(t *testing.T) {
			t.Parallel()
			basics := WhenGeneratingAgentSVGImageBeforeEach(t)

			result := basics.sut.RenderStars(3)

			assert.Contains(t, result, "<text")
			assert.Equal(t, 3, strings.Count(result, "\u2605"))
		})

		t.Run("Should cap stars at 8", func(t *testing.T) {
			t.Parallel()
			basics := WhenGeneratingAgentSVGImageBeforeEach(t)

			result := basics.sut.RenderStars(10)

			assert.Equal(t, 8, strings.Count(result, "\u2605"))
		})
	})

	t.Run("Given RenderMark is called with no mark path", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return empty string", func(t *testing.T) {
			t.Parallel()
			basics := WhenGeneratingAgentSVGImageBeforeEach(t)

			result, err := basics.sut.RenderMark()

			assert.Nil(t, err)
			assert.Empty(t, result)
		})
	})

	t.Run("Given EstimateStrLen is called", func(t *testing.T) {
		t.Parallel()

		t.Run("Should estimate string length based on character width", func(t *testing.T) {
			t.Parallel()
			basics := WhenGeneratingAgentSVGImageBeforeEach(t)

			result := basics.sut.EstimateStrLen("hello", 13)

			charWidth := float32(13) * 0.6
			expected := int(float32(5) * charWidth)
			assert.Equal(t, expected, result)
		})
	})
}
