package mappings_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_agent_identity_registry/mappings"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestWhenMappingExternalAgentResponse(t *testing.T) {
	t.Parallel()

	t.Run("Given a fully populated response", func(t *testing.T) {
		t.Parallel()

		agentId := uuid.New()
		owner := "0xowner"
		external := &mappings.ExternalAgentResponse{
			Id:                lo.ToPtr(agentId.String()),
			Name:              lo.ToPtr("Ada"),
			Owner:             &owner,
			EnsDomain:         lo.ToPtr("ada.eth"),
			Addresses:         lo.ToPtr(map[string]string{"evm": "0x1"}),
			Attributes:        lo.ToPtr(map[string]string{"role": "clone"}),
			CollectionAddress: lo.ToPtr("0xabc"),
			CollectionTokenId: lo.ToPtr(7),
		}

		t.Run("Should map every field", func(t *testing.T) {
			t.Parallel()

			result := external.ToRegistryAgent()

			assert.Equal(t, agentId, result.ID)
			assert.Equal(t, "Ada", result.Name)
			assert.Equal(t, &owner, result.Owner)
			assert.Equal(t, "ada.eth", result.EnsDomain)
			assert.Equal(t, map[string]string{"evm": "0x1"}, result.Addresses)
			assert.Equal(t, map[string]string{"role": "clone"}, result.Attributes)
			assert.Equal(t, "0xabc", result.CollectionAddress)
			assert.Equal(t, uint64(7), *result.CollectionTokenId)
		})
	})

	t.Run("Given a nil response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			var external *mappings.ExternalAgentResponse

			assert.Nil(t, external.ToRegistryAgent())
		})
	})

	t.Run("Given an empty response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should map zero values", func(t *testing.T) {
			t.Parallel()

			result := (&mappings.ExternalAgentResponse{}).ToRegistryAgent()

			assert.Equal(t, uuid.Nil, result.ID)
			assert.Empty(t, result.Name)
			assert.Nil(t, result.Owner)
			assert.Nil(t, result.CollectionTokenId)
		})
	})

	t.Run("Given an invalid agent id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should map the nil uuid", func(t *testing.T) {
			t.Parallel()

			result := (&mappings.ExternalAgentResponse{Id: lo.ToPtr("not-a-uuid")}).ToRegistryAgent()

			assert.Equal(t, uuid.Nil, result.ID)
		})
	})
}
