package settings_test

import (
	"testing"

	processing_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_internal/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingProcessingSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("processing", map[string]any{
			"facilitator_tick_interval":                        uint64(10),
			"conversation_force_conclude_by_facilitator_after": uint64(20),
			"conversation_force_conclude_after":                uint64(30),
			"agent_processing_interval":                        uint64(5),
		})

		settings := processing_settings.NewProcessingSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, uint64(10), settings.FacilitatorTickInterval)
		assert.Equal(t, uint64(20), settings.ConversationForceConcludeByFacilitatorAfter)
		assert.Equal(t, uint64(30), settings.ConversationForceConcludeAfter)
		assert.Equal(t, uint64(5), settings.AgentProcessingInterval)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("PROCESSING__FACILITATOR_TICK_INTERVAL", "11")
		t.Setenv("PROCESSING__CONVERSATION_FORCE_CONCLUDE_BY_FACILITATOR_AFTER", "22")
		t.Setenv("PROCESSING__CONVERSATION_FORCE_CONCLUDE_AFTER", "33")
		t.Setenv("PROCESSING__AGENT_PROCESSING_INTERVAL", "6")

		v := viper.New()
		settings := processing_settings.NewProcessingSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, uint64(11), settings.FacilitatorTickInterval)
		assert.Equal(t, uint64(22), settings.ConversationForceConcludeByFacilitatorAfter)
		assert.Equal(t, uint64(33), settings.ConversationForceConcludeAfter)
		assert.Equal(t, uint64(6), settings.AgentProcessingInterval)
	})
}
