package settings_test

import (
	"testing"

	processing_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingProcessingSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"processing.facilitator_tick_interval":                        uint64(10),
			"processing.conversation_force_conclude_by_facilitator_after": uint64(20),
			"processing.conversation_force_conclude_after":                uint64(30),
			"processing.agent_processing_interval":                        uint64(5),
		}, "."), nil)

		settings := processing_settings.NewProcessingSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, uint64(10), settings.FacilitatorTickInterval)
		assert.Equal(t, uint64(20), settings.ConversationForceConcludeByFacilitatorAfter)
		assert.Equal(t, uint64(30), settings.ConversationForceConcludeAfter)
		assert.Equal(t, uint64(5), settings.AgentProcessingInterval)
	})
}
