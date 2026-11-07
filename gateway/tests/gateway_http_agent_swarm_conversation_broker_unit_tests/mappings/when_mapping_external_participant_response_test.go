package mappings_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/mappings"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestWhenMappingExternalParticipantResponse(t *testing.T) {
	t.Parallel()

	t.Run("Given a fully populated response", func(t *testing.T) {
		t.Parallel()

		participantId := uuid.New()
		collectionAddress := "0xabc"
		ensDomain := "ada.experts.80002.6022.eth"
		displayName := "Ada"
		external := &mappings.ExternalParticipantResponse{
			Id:                lo.ToPtr(participantId.String()),
			CollectionAddress: &collectionAddress,
			CollectionTokenId: lo.ToPtr(7),
			EnsDomain:         &ensDomain,
			DisplayName:       &displayName,
		}

		t.Run("Should map every field", func(t *testing.T) {
			t.Parallel()

			result := external.ToUseCaseResponse()

			assert.Equal(t, participantId, result.ID)
			assert.Equal(t, &collectionAddress, result.CollectionAddress)
			assert.Equal(t, uint64(7), *result.CollectionTokenId)
			assert.Equal(t, &ensDomain, result.EnsDomain)
			assert.Equal(t, &displayName, result.DisplayName)
		})
	})

	t.Run("Given a nil response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			var external *mappings.ExternalParticipantResponse

			assert.Nil(t, external.ToUseCaseResponse())
		})
	})

	t.Run("Given an empty response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should map zero values", func(t *testing.T) {
			t.Parallel()

			result := (&mappings.ExternalParticipantResponse{}).ToUseCaseResponse()

			assert.Equal(t, uuid.Nil, result.ID)
			assert.Nil(t, result.CollectionAddress)
			assert.Nil(t, result.CollectionTokenId)
			assert.Nil(t, result.EnsDomain)
			assert.Nil(t, result.DisplayName)
		})
	})

	t.Run("Given invalid identifiers", func(t *testing.T) {
		t.Parallel()

		external := &mappings.ExternalParticipantResponse{
			Id: lo.ToPtr("not-a-uuid"),
		}

		t.Run("Should map the nil uuid", func(t *testing.T) {
			t.Parallel()

			result := external.ToUseCaseResponse()

			assert.Equal(t, uuid.Nil, result.ID)
		})
	})
}
