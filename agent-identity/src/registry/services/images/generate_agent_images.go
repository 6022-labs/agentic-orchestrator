package images

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services/interfaces"
)

type GenerateAgentImages struct {
	defaultImageGenerator interfaces.DefaultImageGeneratorInterface
	iconImageGenerator    interfaces.IconImageGeneratorInterface
}

func NewGenerateAgentImages(
	defaultImageGenerator interfaces.DefaultImageGeneratorInterface,
	iconImageGenerator interfaces.IconImageGeneratorInterface,
) *GenerateAgentImages {
	return &GenerateAgentImages{
		defaultImageGenerator: defaultImageGenerator,
		iconImageGenerator:    iconImageGenerator,
	}
}

func (g *GenerateAgentImages) Generate(params interfaces.GenerateImageParams) (map[string][]byte, error) {
	defaultBytes, err := g.defaultImageGenerator.Generate(params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate default image: %w", err)
	}

	iconBytes, err := g.iconImageGenerator.Generate(params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate icon image: %w", err)
	}

	return map[string][]byte{
		"default": defaultBytes,
		"icon":    iconBytes,
	}, nil
}
