package configurations

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-identity/images-generator/src/images_generator/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-identity/images-generator/src/images_generator/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-identity/images-generator/src/images_generator/use_cases"
	"go.uber.org/dig"
)

func AddImagesGeneratorConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(settings.NewAgentImageGenerationSettings)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		services.NewAgentSVGImage,
		dig.As(new(services.AgentSVGImageInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewDefaultImageGenerator,
		dig.As(new(services.DefaultImageGeneratorInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewIconImageGenerator,
		dig.As(new(services.IconImageGeneratorInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Use Cases
	err = container.Provide(
		use_cases.NewGenerateImages,
		dig.As(new(use_cases.GenerateImagesInterface)),
	)
	if err != nil {
		panic(err)
	}
}
