package settings

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const ProcessingSettingsKey = "processing"

type ProcessingSettings struct {
	AgentProcessingInterval                     uint64 `mapstructure:"agent_processing_interval"`
	FacilitatorTickInterval                     uint64 `mapstructure:"facilitator_tick_interval"`
	ConversationForceConcludeByFacilitatorAfter uint64 `mapstructure:"conversation_force_conclude_by_facilitator_after"`
	ConversationForceConcludeAfter              uint64 `mapstructure:"conversation_force_conclude_after"`
}

func NewProcessingSettings(
	logger *zap.Logger,
	v *viper.Viper,
) *ProcessingSettings {
	v.SetDefault(ProcessingSettingsKey+".facilitator_tick_interval", uint64(30))
	v.SetDefault(ProcessingSettingsKey+".conversation_force_conclude_by_facilitator_after", uint64(300))
	v.SetDefault(ProcessingSettingsKey+".agent_processing_interval", uint64(20))

	if val, ok := os.LookupEnv("PROCESSING__FACILITATOR_TICK_INTERVAL"); ok {
		v.Set(ProcessingSettingsKey+".facilitator_tick_interval", val)
	}
	if val, ok := os.LookupEnv("PROCESSING__CONVERSATION_FORCE_CONCLUDE_BY_FACILITATOR_AFTER"); ok {
		v.Set(ProcessingSettingsKey+".conversation_force_conclude_by_facilitator_after", val)
	}
	if val, ok := os.LookupEnv("PROCESSING__CONVERSATION_FORCE_CONCLUDE_AFTER"); ok {
		v.Set(ProcessingSettingsKey+".conversation_force_conclude_after", val)
	}
	if val, ok := os.LookupEnv("PROCESSING__AGENT_PROCESSING_INTERVAL"); ok {
		v.Set(ProcessingSettingsKey+".agent_processing_interval", val)
	}

	var settings ProcessingSettings
	err := v.UnmarshalKey(ProcessingSettingsKey, &settings)
	if err != nil {
		logger.Fatal("failed to unmarshal processing settings", zap.Error(err))
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
