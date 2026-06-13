package images_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services/interfaces"
)

// sampleImageParams builds the params shared by the image tests. The base image
// is a real PNG so tests that decode it (the icon) get valid input.
func sampleImageParams(t *testing.T) interfaces.GenerateImageParams {
	t.Helper()

	return interfaces.GenerateImageParams{
		Name:              "my-agent",
		EnsDomain:         "my-agent.agent.eth",
		Role:              "assistant",
		Creator:           "0xabc",
		BaseImageBytes:    sampleBasePNG(t),
		BaseImageMimeType: "image/png",
		Wallets:           []string{"0x1111", "0x2222"},
		CreatedAt:         time.Unix(0, 0),
		Stars:             3,
	}
}

// sampleBasePNG builds a throwaway gradient PNG so the embedded base image shows
// something recognisable without committing a binary fixture.
func sampleBasePNG(t *testing.T) []byte {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, 472, 728))
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, color.RGBA{R: uint8(x % 256), G: uint8(y % 256), B: 160, A: 255})
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("failed to encode sample png: %v", err)
	}
	return buf.Bytes()
}
