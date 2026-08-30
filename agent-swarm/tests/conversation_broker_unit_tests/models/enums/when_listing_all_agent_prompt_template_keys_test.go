package enums_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/stretchr/testify/assert"
)

func TestWhenListingAllAgentPromptTemplateKeys(t *testing.T) {
	t.Parallel()

	t.Run("Given the full key list", func(t *testing.T) {
		t.Parallel()

		t.Run("Should carry every key a role may customize, so Set accepts it", func(t *testing.T) {
			t.Parallel()

			all := enums.GetAllAgentPromptTemplateKeys()

			roles := []enums.ParticipantRole{
				enums.ParticipantRoleClone,
				enums.ParticipantRoleExpert,
				enums.ParticipantRoleFacilitator,
				enums.ParticipantRoleHuman,
			}
			for _, role := range roles {
				for _, key := range role.AllowedPromptTemplateKeys() {
					assert.True(t, all.Contains(key), "%s is allowed for %s but missing from the full list", key, role)
				}
			}
		})

		t.Run("Should let Set parse every listed key", func(t *testing.T) {
			t.Parallel()

			for _, key := range enums.GetAllAgentPromptTemplateKeys() {
				var parsed enums.AgentPromptTemplateKey
				assert.NoError(t, parsed.Set(key.String()))
				assert.Equal(t, key, parsed)
			}
		})
	})
}
