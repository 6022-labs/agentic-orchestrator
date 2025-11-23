package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type AgentImageGenerationSettings struct {
	ImageMarkPath *string
}

func NewAgentImageGenerationSettings(logger *zap.Logger) *AgentImageGenerationSettings {
	var imageMarkPath *string = nil
	imageMarkPathEnv := os.Getenv("AGENT_IMAGE_GENERATION_IMAGE_MARK_PATH")

	if len(strings.TrimSpace(imageMarkPathEnv)) != 0 {
		if _, err := os.Stat(imageMarkPathEnv); os.IsNotExist(err) {
			logger.Warn("The file specified in AGENT_IMAGE_GENERATION_IMAGE_MARK_PATH does not exist", zap.String("path", imageMarkPathEnv))
		} else {
			imageMarkPath = &imageMarkPathEnv
		}
	}

	if imageMarkPath == nil {
		logger.Info("AGENT_IMAGE_GENERATION_IMAGE_MARK_PATH is not set, nft generation will not be marked")
	}

	return &AgentImageGenerationSettings{
		ImageMarkPath: imageMarkPath,
	}
}
