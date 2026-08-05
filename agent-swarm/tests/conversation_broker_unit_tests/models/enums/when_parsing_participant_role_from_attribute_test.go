package enums_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/stretchr/testify/assert"
)

func TestWhenParsingParticipantRoleFromAttribute(t *testing.T) {
	t.Parallel()

	t.Run("Given the attribute is a known role", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the matching role", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, enums.ParticipantRoleFacilitator, enums.ParticipantRoleFromAttribute("facilitator"))
			assert.Equal(t, enums.ParticipantRoleClone, enums.ParticipantRoleFromAttribute("Clone"))
		})
	})

	t.Run("Given the attribute is unknown or empty", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fall back to expert", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, enums.ParticipantRoleExpert, enums.ParticipantRoleFromAttribute("banana"))
			assert.Equal(t, enums.ParticipantRoleExpert, enums.ParticipantRoleFromAttribute(""))
		})
	})
}
