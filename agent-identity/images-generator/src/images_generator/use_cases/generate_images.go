package use_cases

import (
	"fmt"
	"time"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-identity/images-generator/src/images_generator/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-identity/images-generator/src/images_generator/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-identity/images-generator/src/images_generator/use_cases/responses"
)

type GenerateImagesInterface interface {
	Generate(request *requests.GenerateImagesRequest) (*responses.GenerateImagesResponse, error)
}

type GenerateImages struct {
	defaultImageGenerator services.DefaultImageGeneratorInterface
	iconImageGenerator    services.IconImageGeneratorInterface
}

func NewGenerateImages(
	defaultImageGenerator services.DefaultImageGeneratorInterface,
	iconImageGenerator services.IconImageGeneratorInterface,
) *GenerateImages {
	return &GenerateImages{
		defaultImageGenerator: defaultImageGenerator,
		iconImageGenerator:    iconImageGenerator,
	}
}

func (g *GenerateImages) Generate(request *requests.GenerateImagesRequest) (*responses.GenerateImagesResponse, error) {
	params := services.GenerateImageParams{
		Name:              request.Name,
		Role:              request.Role,
		Creator:           request.Creator,
		EnsDomain:         request.EnsDomain,
		OriginEnsDomain:   request.OriginEnsDomain,
		OriginGeneration:  request.OriginGeneration,
		BaseImageBytes:    request.BaseImageBytes,
		BaseImageMimeType: request.BaseImageMimeType,
		Wallets:           request.Wallets,
		CreatedAt:         time.Now(),
		Stars:             request.Stars,
	}

	defaultImageBytes, err := g.defaultImageGenerator.Generate(params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate default image: %w", err)
	}

	iconImageBytes, err := g.iconImageGenerator.Generate(params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate icon image: %w", err)
	}

	return &responses.GenerateImagesResponse{
		Images: map[string][]byte{
			"default": defaultImageBytes,
			"icon":    iconImageBytes,
		},
	}, nil
}
