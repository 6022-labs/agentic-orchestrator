package settings_test

import (
	"testing"

	processing_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingProcessingSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"processing.participant_processing_interval":                  uint64(15),
			"processing.facilitator_tick_interval":                        uint64(45),
			"processing.conversation_force_conclude_by_facilitator_after": uint64(120),
			"processing.conversation_force_conclude_after":                uint64(900),
		}, "."), nil)

		settings := processing_settings.NewProcessingSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, uint64(15), settings.ParticipantProcessingInterval)
		assert.Equal(t, uint64(45), settings.FacilitatorTickInterval)
		assert.Equal(t, uint64(120), settings.ConversationForceConcludeByFacilitatorAfter)
		assert.Equal(t, uint64(900), settings.ConversationForceConcludeAfter)
	})

	t.Run("Defaults missing values", func(t *testing.T) {
		k := koanf.New(".")

		settings := processing_settings.NewProcessingSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, uint64(20), settings.ParticipantProcessingInterval)
		assert.Equal(t, uint64(30), settings.FacilitatorTickInterval)
		assert.Equal(t, uint64(300), settings.ConversationForceConcludeByFacilitatorAfter)
		assert.Equal(t, uint64(600), settings.ConversationForceConcludeAfter)
	})
}
