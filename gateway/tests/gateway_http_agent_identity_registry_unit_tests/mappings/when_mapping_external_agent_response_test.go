package mappings_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/mappings"
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
			Creator:           lo.ToPtr("0xcreator"),
			Images:            lo.ToPtr(map[string]string{"avatar": "ipfs://img"}),
			CloneOf:           lo.ToPtr(3),
			Attributes:        lo.ToPtr(map[string]string{"role": "expert"}),
			EnsDomain:         lo.ToPtr("ada.experts.80002.6022.eth"),
			Addresses:         lo.ToPtr(map[string]string{"evm": "0x1"}),
			ProposalId:        lo.ToPtr(9),
			CollectionTokenId: lo.ToPtr(7),
			CollectionAddress: lo.ToPtr("0xabc"),
		}

		t.Run("Should map every field", func(t *testing.T) {
			t.Parallel()

			result := external.ToUseCaseResponse()

			assert.Equal(t, agentId, result.ID)
			assert.Equal(t, "Ada", result.Name)
			assert.Equal(t, &owner, result.Owner)
			assert.Equal(t, "0xcreator", result.Creator)
			assert.Equal(t, map[string]string{"avatar": "ipfs://img"}, result.Images)
			assert.Equal(t, uint64(3), *result.CloneOf)
			assert.Equal(t, map[string]string{"role": "expert"}, result.Attributes)
			assert.Equal(t, "ada.experts.80002.6022.eth", result.EnsDomain)
			assert.Equal(t, map[string]string{"evm": "0x1"}, result.Addresses)
			assert.Equal(t, uint64(9), *result.ProposalId)
			assert.Equal(t, uint64(7), *result.CollectionTokenId)
			assert.Equal(t, "0xabc", result.CollectionAddress)
		})
	})

	t.Run("Given a nil response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			var external *mappings.ExternalAgentResponse

			assert.Nil(t, external.ToUseCaseResponse())
		})
	})

	t.Run("Given an empty response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should map zero values", func(t *testing.T) {
			t.Parallel()

			result := (&mappings.ExternalAgentResponse{}).ToUseCaseResponse()

			assert.Equal(t, uuid.Nil, result.ID)
			assert.Empty(t, result.Name)
			assert.Nil(t, result.Owner)
			assert.Nil(t, result.CloneOf)
			assert.Nil(t, result.ProposalId)
			assert.Nil(t, result.CollectionTokenId)
		})
	})

	t.Run("Given an invalid agent id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should map the nil uuid", func(t *testing.T) {
			t.Parallel()

			result := (&mappings.ExternalAgentResponse{Id: lo.ToPtr("not-a-uuid")}).ToUseCaseResponse()

			assert.Equal(t, uuid.Nil, result.ID)
		})
	})
}
