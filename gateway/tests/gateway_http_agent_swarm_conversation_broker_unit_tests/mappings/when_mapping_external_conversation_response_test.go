package mappings_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/mappings"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestWhenMappingExternalConversationResponse(t *testing.T) {
	t.Parallel()

	t.Run("Given a fully populated response", func(t *testing.T) {
		t.Parallel()

		conversationId := uuid.New()
		bridgeId := uuid.New()
		description := "demo"
		external := &mappings.ExternalConversationResponse{
			Id:          lo.ToPtr(conversationId.String()),
			BridgeId:    lo.ToPtr(bridgeId.String()),
			Description: &description,
			Status:      lo.ToPtr(clients.EnumsConversationStatus("concluded")),
			ConcludedAt: lo.ToPtr("2026-07-04T12:00:00Z"),
			CreatedAt:   lo.ToPtr("2026-07-04T10:00:00Z"),
			UpdatedAt:   lo.ToPtr("2026-07-04T11:00:00Z"),
		}

		t.Run("Should map every field", func(t *testing.T) {
			t.Parallel()

			result := external.ToUseCaseResponse()

			assert.Equal(t, conversationId, result.ID)
			assert.Equal(t, bridgeId, result.BridgeId)
			assert.Equal(t, &description, result.Description)
			assert.Equal(t, "concluded", result.Status)
			if assert.NotNil(t, result.ConcludedAt) {
				assert.Equal(t, time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC), *result.ConcludedAt)
			}
			assert.Equal(t, time.Date(2026, 7, 4, 10, 0, 0, 0, time.UTC), result.CreatedAt)
			assert.Equal(t, time.Date(2026, 7, 4, 11, 0, 0, 0, time.UTC), result.UpdatedAt)
		})
	})

	t.Run("Given a nil response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			var external *mappings.ExternalConversationResponse

			assert.Nil(t, external.ToUseCaseResponse())
		})
	})

	t.Run("Given an empty response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should map zero values", func(t *testing.T) {
			t.Parallel()

			result := (&mappings.ExternalConversationResponse{}).ToUseCaseResponse()

			assert.Equal(t, uuid.Nil, result.ID)
			assert.Equal(t, uuid.Nil, result.BridgeId)
			assert.Nil(t, result.Description)
			assert.Empty(t, result.Status)
			assert.Nil(t, result.ConcludedAt)
			assert.True(t, result.CreatedAt.IsZero())
			assert.True(t, result.UpdatedAt.IsZero())
		})
	})

	t.Run("Given invalid identifiers and timestamps", func(t *testing.T) {
		t.Parallel()

		external := &mappings.ExternalConversationResponse{
			Id:          lo.ToPtr("not-a-uuid"),
			BridgeId:    lo.ToPtr("not-a-uuid"),
			ConcludedAt: lo.ToPtr("not-a-time"),
			CreatedAt:   lo.ToPtr("not-a-time"),
			UpdatedAt:   lo.ToPtr("not-a-time"),
		}

		t.Run("Should map zero values", func(t *testing.T) {
			t.Parallel()

			result := external.ToUseCaseResponse()

			assert.Equal(t, uuid.Nil, result.ID)
			assert.Equal(t, uuid.Nil, result.BridgeId)
			assert.Nil(t, result.ConcludedAt)
			assert.True(t, result.CreatedAt.IsZero())
			assert.True(t, result.UpdatedAt.IsZero())
		})
	})
}
