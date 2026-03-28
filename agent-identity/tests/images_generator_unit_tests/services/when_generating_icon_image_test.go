package services_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/images_generator/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/images_generator/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/images_generator_mocks/services_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGeneratingIconImageTestingSuite struct {
	sut *services.IconImageGenerator

	mockAgentSVGImage            *services_mocks.MockAgentSVGImageInterface
	agentImageGenerationSettings *settings.ImageGenerationSettings
}

func WhenGeneratingIconImageBeforeEach(t *testing.T) *WhenGeneratingIconImageTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentSVGImage := services_mocks.NewMockAgentSVGImageInterface(mockController)
	agentImageGenerationSettings := &settings.ImageGenerationSettings{
		ImageMarkPath: nil,
	}

	sut := services.NewIconImageGenerator(mockAgentSVGImage, agentImageGenerationSettings)

	return &WhenGeneratingIconImageTestingSuite{
		sut:                          sut,
		mockAgentSVGImage:            mockAgentSVGImage,
		agentImageGenerationSettings: agentImageGenerationSettings,
	}
}

func TestWhenGeneratingIconImage(t *testing.T) {
	t.Parallel()

	t.Run("Given image generation is requested with valid parameters", func(t *testing.T) {
		t.Parallel()

		t.Run("Should call AgentSVGImage methods and produce image bytes", func(t *testing.T) {
			t.Parallel()
			basics := WhenGeneratingIconImageBeforeEach(t)

			params := services.GenerateImageParams{
				Name:              "TestAgent",
				EnsDomain:         "test.6022.eth",
				Role:              "expert",
				Creator:           "0x1234",
				BaseImageBytes:    []byte("base-image"),
				BaseImageMimeType: "image/png",
				Wallets:           []string{"0xwallet1"},
				CreatedAt:         time.Now(),
				Stars:             3,
			}

			basics.mockAgentSVGImage.EXPECT().
				RenderBorderGradient(params.Creator).
				Return(`<linearGradient id="borderGradient"></linearGradient>`).
				Times(1)

			basics.mockAgentSVGImage.EXPECT().
				RenderImageData(params.BaseImageBytes, params.BaseImageMimeType).
				Return("data:image/png;base64,test").
				Times(1)

			basics.mockAgentSVGImage.EXPECT().
				RenderAgentInformation(gomock.Any()).
				Return(`<g>agent info</g>`).
				Times(1)

			basics.mockAgentSVGImage.EXPECT().
				RenderStars(params.Stars).
				Return(`<text>stars</text>`).
				Times(1)

			result, err := basics.sut.Generate(params)

			assert.Nil(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result)
		})
	})

	t.Run("Given image mark path is configured and RenderMark fails", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error from RenderMark", func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)

			mockAgentSVGImage := services_mocks.NewMockAgentSVGImageInterface(mockController)
			markPath := "/tmp/test-mark.png"
			agentImageGenerationSettings := &settings.ImageGenerationSettings{
				ImageMarkPath: &markPath,
			}

			sut := services.NewIconImageGenerator(mockAgentSVGImage, agentImageGenerationSettings)

			params := services.GenerateImageParams{
				Name:              "TestAgent",
				EnsDomain:         "test.6022.eth",
				Role:              "expert",
				Creator:           "0x1234",
				BaseImageBytes:    []byte("base-image"),
				BaseImageMimeType: "image/png",
				Wallets:           []string{},
				CreatedAt:         time.Now(),
				Stars:             0,
			}

			mockAgentSVGImage.EXPECT().
				RenderBorderGradient(params.Creator).
				Return(`<linearGradient id="borderGradient"></linearGradient>`).
				Times(1)

			mockAgentSVGImage.EXPECT().
				RenderImageData(params.BaseImageBytes, params.BaseImageMimeType).
				Return("data:image/png;base64,test").
				Times(1)

			mockAgentSVGImage.EXPECT().
				RenderAgentInformation(gomock.Any()).
				Return(`<g>agent info</g>`).
				Times(1)

			mockAgentSVGImage.EXPECT().
				RenderStars(params.Stars).
				Return("").
				Times(1)

			mockAgentSVGImage.EXPECT().
				RenderMark().
				Return("", assert.AnError).
				Times(1)

			_, err := sut.Generate(params)

			assert.NotNil(t, err)
			assert.ErrorContains(t, err, "failed to add image mark")
		})
	})
}
