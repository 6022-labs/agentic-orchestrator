package services_test

import (
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-identity/images-generator/src/images_generator/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-identity/images-generator/src/images_generator/settings"
	"github.com/stretchr/testify/assert"
)

type WhenGeneratingAgentAnimatedSVGTestingSuiteSuite struct {
	sut *services.DefaultImageGenerator
}

func WhenGeneratingAgentAnimatedSVGBeforeEach() WhenGeneratingAgentAnimatedSVGTestingSuiteSuite {
	imagePath := "../../../src/orchestrator_dapp_api/images/nft_mark.png"
	agentImageGenerationSettings := &settings.AgentImageGenerationSettings{
		ImageMarkPath: &imagePath,
	}

	sut := services.NewDefaultImageGenerator(
		services.NewAgentSVGImage(
			agentImageGenerationSettings,
		),
		agentImageGenerationSettings,
	)

	return WhenGeneratingAgentAnimatedSVGTestingSuiteSuite{
		sut: sut,
	}
}

func generateValidGenerateImageParams() *services.GenerateImageParams {
	imageBytes := []byte("")
	resp, err := http.Get("https://ipfs.io/ipfs/bafybeia25xzo4k4nipuw6cyml4f4x3b3h6lneyvch6a4b3rl2q7e6hpgra")
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		imageBytes, _ = io.ReadAll(resp.Body)
	}

	originEnsDomain := "mars.agent.protocol6022.eth"
	originGeneration := uint(2)

	return &services.GenerateImageParams{
		Name:              "jupiter",
		CreatedAt:         time.Now(),
		Creator:           "0xd72a7503C30613A8F568de1f3ACF829199b66e2D",
		OriginEnsDomain:   &originEnsDomain,
		OriginGeneration:  &originGeneration,
		Stars:             8,
		Role:              "Facilitator",
		BaseImageBytes:    imageBytes,
		BaseImageMimeType: "image/png",
		Wallets:           []string{"0xbE4d33D1B5e0f5258eDe32644517D23Eb9E89232", "0x2A47A133eA3B8AceD8a646ca7950A9dEFB89E0e3"},
	}
}

func TestWhenGeneratingDefaultImage(t *testing.T) {
	t.Parallel()

	t.Run("Should generate a valid SVG image", func(t *testing.T) {
		t.Parallel()

		basics := WhenGeneratingAgentAnimatedSVGBeforeEach()

		params := generateValidGenerateImageParams()

		result, err := basics.sut.Generate(*params)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Contains(t, string(result), "<svg")
		assert.Contains(t, string(result), params.Name)
		assert.Contains(t, string(result), base64.StdEncoding.EncodeToString(params.BaseImageBytes))

		// Copy result into file:///Users/pierrefourny/Downloads/test.svg
		file, err := os.Create("/Users/pierrefourny/Downloads/test.svg")
		if err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
		defer file.Close()
		_, err = file.Write(result)
		if err != nil {
			t.Fatalf("Failed to write to file: %v", err)
		}
		// Close the file after writing
		err = file.Close()
		if err != nil {
			t.Fatalf("Failed to close file: %v", err)
		}
	})

	t.Run("Should handle empty border colors gracefully", func(t *testing.T) {
		t.Parallel()

		basics := WhenGeneratingAgentAnimatedSVGBeforeEach()

		params := generateValidGenerateImageParams()

		result, err := basics.sut.Generate(*params)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Contains(t, string(result), "<svg")
		assert.Contains(t, string(result), params.Name)
	})
}
