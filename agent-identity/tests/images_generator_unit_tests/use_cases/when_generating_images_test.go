package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/images_generator/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/images_generator/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/images_generator/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/images_generator_mocks/services_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGeneratingImagesTestingSuite struct {
	sut *use_cases.GenerateImages

	mockDefaultImageGenerator *services_mocks.MockDefaultImageGeneratorInterface
	mockIconImageGenerator    *services_mocks.MockIconImageGeneratorInterface
}

func WhenGeneratingImagesBeforeEach(t *testing.T) *WhenGeneratingImagesTestingSuite {
	mockController := gomock.NewController(t)

	mockDefaultImageGenerator := services_mocks.NewMockDefaultImageGeneratorInterface(mockController)
	mockIconImageGenerator := services_mocks.NewMockIconImageGeneratorInterface(mockController)

	sut := use_cases.NewGenerateImages(mockDefaultImageGenerator, mockIconImageGenerator)

	return &WhenGeneratingImagesTestingSuite{
		sut:                       sut,
		mockDefaultImageGenerator: mockDefaultImageGenerator,
		mockIconImageGenerator:    mockIconImageGenerator,
	}
}

func TestWhenGeneratingImages(t *testing.T) {
	t.Parallel()

	t.Run("Given default image generation fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGeneratingImagesTestingSuite) {
			suite.mockDefaultImageGenerator.EXPECT().
				Generate(gomock.Any()).
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenGeneratingImagesBeforeEach(t)
			initSuite(basics)

			req := &requests.GenerateImagesRequest{
				Name:    "TestAgent",
				Role:    "expert",
				Creator: "0x1234",
			}
			resp, err := basics.sut.Generate(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.ErrorContains(t, err, "failed to generate default image")
		})
	})

	t.Run("Given icon image generation fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGeneratingImagesTestingSuite) {
			suite.mockDefaultImageGenerator.EXPECT().
				Generate(gomock.Any()).
				Return([]byte("default-image-bytes"), nil).
				Times(1)

			suite.mockIconImageGenerator.EXPECT().
				Generate(gomock.Any()).
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenGeneratingImagesBeforeEach(t)
			initSuite(basics)

			req := &requests.GenerateImagesRequest{
				Name:    "TestAgent",
				Role:    "expert",
				Creator: "0x1234",
			}
			resp, err := basics.sut.Generate(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.ErrorContains(t, err, "failed to generate icon image")
		})
	})

	t.Run("Given images are generated successfully", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the images response with default and icon keys", func(t *testing.T) {
			t.Parallel()
			basics := WhenGeneratingImagesBeforeEach(t)

			defaultImageBytes := []byte("default-image-bytes")
			iconImageBytes := []byte("icon-image-bytes")

			basics.mockDefaultImageGenerator.EXPECT().
				Generate(gomock.Any()).
				DoAndReturn(func(params services.GenerateImageParams) ([]byte, error) {
					assert.Equal(t, "TestAgent", params.Name)
					assert.Equal(t, "expert", params.Role)
					assert.Equal(t, "0x1234", params.Creator)
					return defaultImageBytes, nil
				}).
				Times(1)

			basics.mockIconImageGenerator.EXPECT().
				Generate(gomock.Any()).
				Return(iconImageBytes, nil).
				Times(1)

			req := &requests.GenerateImagesRequest{
				Name:    "TestAgent",
				Role:    "expert",
				Creator: "0x1234",
			}
			resp, err := basics.sut.Generate(req)

			assert.Nil(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, defaultImageBytes, resp.Images["default"])
			assert.Equal(t, iconImageBytes, resp.Images["icon"])
			assert.Len(t, resp.Images, 2)
		})
	})
}
