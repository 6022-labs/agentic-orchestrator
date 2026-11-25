package settings

import (
	"time"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const ProcessingSettingsKey = "processing"

type ProcessingSettings struct {
	ParticipantProcessingInterval               uint64
	FacilitatorTickInterval                     uint64
	ConversationForceConcludeByFacilitatorAfter uint64
	ConversationForceConcludeAfter              uint64
	// LeaseDuration is how long a crashed broker keeps a participant; a running one renews it.
	LeaseDuration time.Duration
	// MaxConcurrentParticipants bounds the agents one broker processes at once.
	MaxConcurrentParticipants int
	// MaxConcurrentConversations bounds the conversations one agent answers at once.
	MaxConcurrentConversations int
}

type processingSettingsRaw struct {
	ParticipantProcessingInterval               uint64 `koanf:"participant_processing_interval"`
	FacilitatorTickInterval                     uint64 `koanf:"facilitator_tick_interval"`
	ConversationForceConcludeByFacilitatorAfter uint64 `koanf:"conversation_force_conclude_by_facilitator_after"`
	ConversationForceConcludeAfter              uint64 `koanf:"conversation_force_conclude_after"`
	LeaseDurationSeconds                        uint64 `koanf:"lease_duration_seconds"`
	MaxConcurrentParticipants                   int    `koanf:"max_concurrent_participants"`
	MaxConcurrentConversations                  int    `koanf:"max_concurrent_conversations"`
}

func NewProcessingSettings(logger *zap.Logger, k *koanf.Koanf) *ProcessingSettings {
	settings := processingSettingsRaw{
		ParticipantProcessingInterval:               20,
		FacilitatorTickInterval:                     30,
		ConversationForceConcludeByFacilitatorAfter: 300,
		LeaseDurationSeconds:                        60,
		MaxConcurrentParticipants:                   200,
		MaxConcurrentConversations:                  50,
	}
	if err := k.Unmarshal(ProcessingSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal processing settings", zap.Error(err))
	}

	if settings.ParticipantProcessingInterval == 0 {
		settings.ParticipantProcessingInterval = 20
	}
	if settings.FacilitatorTickInterval == 0 {
		settings.FacilitatorTickInterval = 30
	}
	if settings.ConversationForceConcludeByFacilitatorAfter == 0 {
		settings.ConversationForceConcludeByFacilitatorAfter = 300
	}
	if settings.ConversationForceConcludeAfter == 0 {
		settings.ConversationForceConcludeAfter = settings.ConversationForceConcludeByFacilitatorAfter * 2
		logger.Info(
			"processing.conversation_force_conclude_after is not set, "+
				"defaulting to 2x conversation_force_conclude_by_facilitator_after",
			zap.Uint64("defaultConversationForceConcludeAfter", settings.ConversationForceConcludeAfter),
		)
	}

	if settings.LeaseDurationSeconds == 0 {
		settings.LeaseDurationSeconds = 60
	}
	if settings.MaxConcurrentParticipants <= 0 {
		settings.MaxConcurrentParticipants = 200
	}
	if settings.MaxConcurrentConversations <= 0 {
		settings.MaxConcurrentConversations = 50
	}

	return &ProcessingSettings{
		ParticipantProcessingInterval:               settings.ParticipantProcessingInterval,
		FacilitatorTickInterval:                     settings.FacilitatorTickInterval,
		ConversationForceConcludeByFacilitatorAfter: settings.ConversationForceConcludeByFacilitatorAfter,
		ConversationForceConcludeAfter:              settings.ConversationForceConcludeAfter,
		LeaseDuration:                               time.Duration(settings.LeaseDurationSeconds) * time.Second,
		MaxConcurrentParticipants:                   settings.MaxConcurrentParticipants,
		MaxConcurrentConversations:                  settings.MaxConcurrentConversations,
	}
}
