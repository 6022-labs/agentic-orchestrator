package enums_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/stretchr/testify/assert"
)

func TestWhenGettingAllowedPromptTemplateKeys(t *testing.T) {
	t.Parallel()

	t.Run("Given role is clone", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the clone keys", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, enums.AgentPromptTemplateKeyList{
				enums.AgentPromptTemplateKeyCloneSystem,
				enums.AgentPromptTemplateKeyCloneDefault,
			}, enums.ParticipantRoleClone.AllowedPromptTemplateKeys())
		})
	})

	t.Run("Given role is facilitator", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the facilitator keys", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, enums.AgentPromptTemplateKeyList{
				enums.AgentPromptTemplateKeyFacilitatorSystem,
				enums.AgentPromptTemplateKeyFacilitatorDefault,
				enums.AgentPromptTemplateKeyFacilitatorConclude,
			}, enums.ParticipantRoleFacilitator.AllowedPromptTemplateKeys())
		})
	})

	t.Run("Given role is expert, human or the zero value", func(t *testing.T) {
		t.Parallel()

		expected := enums.AgentPromptTemplateKeyList{
			enums.AgentPromptTemplateKeyExpertSystem,
			enums.AgentPromptTemplateKeyExpertDefault,
		}

		t.Run("Should return the expert keys for expert", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, expected, enums.ParticipantRoleExpert.AllowedPromptTemplateKeys())
		})

		t.Run("Should return the expert keys for human", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, expected, enums.ParticipantRoleHuman.AllowedPromptTemplateKeys())
		})

		t.Run("Should return the expert keys for the zero value", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, expected, enums.ParticipantRole("").AllowedPromptTemplateKeys())
		})
	})
}
