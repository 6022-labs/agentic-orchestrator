package settings

import (
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const ProcessingSettingsKey = "processing"

type ProcessingSettings struct {
	AgentProcessingInterval                     uint64 `koanf:"agent_processing_interval"`
	FacilitatorTickInterval                     uint64 `koanf:"facilitator_tick_interval"`
	ConversationForceConcludeByFacilitatorAfter uint64 `koanf:"conversation_force_conclude_by_facilitator_after"`
	ConversationForceConcludeAfter              uint64 `koanf:"conversation_force_conclude_after"`
}

func NewProcessingSettings(
	logger *zap.Logger,
	k *koanf.Koanf,
) *ProcessingSettings {
	settings := ProcessingSettings{
		FacilitatorTickInterval:                     30,
		ConversationForceConcludeByFacilitatorAfter: 300,
		AgentProcessingInterval:                     20,
	}
	if err := k.Unmarshal(ProcessingSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal processing settings", zap.Error(err))
	}

	if settings.FacilitatorTickInterval == 0 {
		settings.FacilitatorTickInterval = 30
	}

	if settings.ConversationForceConcludeByFacilitatorAfter == 0 {
		settings.ConversationForceConcludeByFacilitatorAfter = 300
	}

	if settings.AgentProcessingInterval == 0 {
		settings.AgentProcessingInterval = 20
	}

	if settings.ConversationForceConcludeAfter == 0 {
		settings.ConversationForceConcludeAfter = settings.ConversationForceConcludeByFacilitatorAfter * 2
		logger.Info("processing.conversation_force_conclude_after is not set, using default value (x2 of processing.conversation_force_conclude_by_facilitator_after)",
			zap.Uint64("defaultConversationForceConcludeAfter", settings.ConversationForceConcludeAfter),
		)
	}

	if settings.FacilitatorTickInterval == 0 || settings.ConversationForceConcludeByFacilitatorAfter == 0 || settings.AgentProcessingInterval == 0 {
		logger.Fatal("processing settings must be greater than 0")
	}

	return &settings
}
