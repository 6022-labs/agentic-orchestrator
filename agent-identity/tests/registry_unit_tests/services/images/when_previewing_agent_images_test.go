package images_test

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services/images"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/settings"
	"github.com/stretchr/testify/assert"
)

// Render sample agent frames into <repo-root>/tmp/agent-images. Run with:
//
//	go generate ./tests/registry_unit_tests/services/images/
//
//go:generate sh -c "AGENT_IMAGE_PREVIEW_DIR=$(git rev-parse --show-toplevel)/tmp/agent-images go test -run TestPreviewAgentImages -count=1 -v ."

// TestPreviewAgentImages is a visual-inspection helper, not an assertion test.
// It is skipped unless AGENT_IMAGE_PREVIEW_DIR is set — `go generate` (above)
// sets it for you. You can also run it directly with the env var on the same
// line:
//
//	AGENT_IMAGE_PREVIEW_DIR=/tmp/agent-images go test ./tests/registry_unit_tests/services/images/ -run TestPreviewAgentImages -count=1 -v
//
// It drops default.svg (open in a browser) + icon.png (the rasterized frame).
func TestPreviewAgentImages(t *testing.T) {
	outDir := os.Getenv("AGENT_IMAGE_PREVIEW_DIR")
	if outDir == "" {
		t.Skip("set AGENT_IMAGE_PREVIEW_DIR to dump preview images for visual inspection")
	}

	base := sampleBasePNG(t)
	creator := "0x4bca9E17a8326d283d1fC68A725B46b42c2D9146"

	cases := []struct {
		name         string
		agentName    string
		ensDomain    string
		wallets      []string
		stars        uint8
		baseImageURL string // optional; falls back to the generated gradient
	}{
		{
			name:      "baseline",
			agentName: "avogadro",
			ensDomain: "avogadro.agent.protocol6022.eth",
			wallets:   []string{"0xf5cfC4723d64C0B53b3aF49C3a184f16E7cACA40"},
			stars:     7,
		},
		{
			// A real fetched image rather than the synthetic gradient.
			name:         "real-image",
			agentName:    "avogadro",
			ensDomain:    "avogadro.agent.protocol6022.eth",
			wallets:      []string{"0xf5cfC4723d64C0B53b3aF49C3a184f16E7cACA40"},
			stars:        6,
			baseImageURL: "https://ipfs.io/ipfs/bafybeia25xzo4k4nipuw6cyml4f4x3b3h6lneyvch6a4b3rl2q7e6hpgra",
		},
		{
			// Both the name and the ENS domain are long enough to switch their
			// labels from static text to the scrolling carousel.
			name:      "long-name-and-ens",
			agentName: "avogadro-the-incredibly-verbose-research-agent",
			ensDomain: "avogadro-the-incredibly-verbose-research-agent.sub.agent.protocol6022.eth",
			wallets:   []string{"0xf5cfC4723d64C0B53b3aF49C3a184f16E7cACA40"},
			stars:     5,
		},
		{
			// Several wallets — a multi-value scrolling carousel on the default
			// frame; the icon shows only the first.
			name:      "multi-wallet",
			agentName: "avogadro",
			ensDomain: "avogadro.agent.protocol6022.eth",
			wallets: []string{
				"0xf5cfC4723d64C0B53b3aF49C3a184f16E7cACA40",
				"0x1111111111111111111111111111111111111111",
				"0x2222222222222222222222222222222222222222",
				"0x3333333333333333333333333333333333333333",
			},
			stars: 8,
		},
	}

	markPath := repoMarkPath()
	imageGenerationSettings := &settings.ImageGenerationSettings{ImageMarkPath: &markPath}
	generator := images.NewGenerateAgentImages(
		images.NewDefaultImageGenerator(imageGenerationSettings),
		images.NewIconImageGenerator(imageGenerationSettings),
	)

	// The default frame is SVG; the icon is a rasterized PNG.
	ext := map[string]string{"default": "svg", "icon": "png"}

	for _, c := range cases {
		baseBytes, baseMime := base, "image/png"
		if c.baseImageURL != "" {
			baseBytes, baseMime = fetchImage(t, c.baseImageURL)
		}

		result, err := generator.Generate(interfaces.GenerateImageParams{
			Name:              c.agentName,
			EnsDomain:         c.ensDomain,
			Role:              "expert",
			Creator:           creator,
			BaseImageBytes:    baseBytes,
			BaseImageMimeType: baseMime,
			Wallets:           c.wallets,
			CreatedAt:         time.Now(),
			Stars:             c.stars,
		})
		assert.NoError(t, err)

		caseDir := filepath.Join(outDir, c.name)
		if err := os.MkdirAll(caseDir, 0o755); err != nil {
			t.Fatalf("failed to create case dir %s: %v", caseDir, err)
		}

		for kind, data := range result {
			path := filepath.Join(caseDir, kind+"."+ext[kind])
			if err := os.WriteFile(path, data, 0o644); err != nil {
				t.Fatalf("failed to write %s: %v", path, err)
			}
			t.Logf("wrote %s (%d bytes)", path, len(data))
		}
	}
}

// fetchImage downloads an image and reports its bytes + mime type, so the
// preview can render a real agent picture instead of the synthetic gradient.
func fetchImage(t *testing.T, url string) ([]byte, string) {
	t.Helper()

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("failed to fetch %s: %v", url, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read %s: %v", url, err)
	}

	mime := resp.Header.Get("Content-Type")
	if mime == "" {
		mime = http.DetectContentType(data)
	}
	return data, mime
}

// repoMarkPath resolves the watermark relative to this test file, so the
// preview shows the mark regardless of the working directory.
func repoMarkPath() string {
	_, thisFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..", "src", "registry_host", "images", "nft_mark.png")
}
