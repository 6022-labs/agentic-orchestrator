package mappings_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/mappings"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestWhenMappingExternalMessageResponse(t *testing.T) {
	t.Parallel()

	t.Run("Given a fully populated response", func(t *testing.T) {
		t.Parallel()

		messageId := uuid.New()
		conversationId := uuid.New()
		participantId := uuid.New()
		external := &mappings.ExternalMessageResponse{
			Id:             lo.ToPtr(messageId.String()),
			ConversationId: lo.ToPtr(conversationId.String()),
			ParticipantId:  lo.ToPtr(participantId.String()),
			Content:        lo.ToPtr("hello"),
			CreatedAt:      lo.ToPtr("2026-07-04T10:00:00Z"),
		}

		t.Run("Should map every field", func(t *testing.T) {
			t.Parallel()

			result := external.ToUseCaseResponse()

			assert.Equal(t, messageId, result.ID)
			assert.Equal(t, conversationId, result.ConversationId)
			assert.Equal(t, participantId, result.ParticipantId)
			assert.Equal(t, "hello", result.Content)
			assert.Equal(t, time.Date(2026, 7, 4, 10, 0, 0, 0, time.UTC), result.CreatedAt)
		})
	})

	t.Run("Given a nil response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			var external *mappings.ExternalMessageResponse

			assert.Nil(t, external.ToUseCaseResponse())
		})
	})

	t.Run("Given an empty response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should map zero values", func(t *testing.T) {
			t.Parallel()

			result := (&mappings.ExternalMessageResponse{}).ToUseCaseResponse()

			assert.Equal(t, uuid.Nil, result.ID)
			assert.Equal(t, uuid.Nil, result.ConversationId)
			assert.Equal(t, uuid.Nil, result.ParticipantId)
			assert.Empty(t, result.Content)
			assert.True(t, result.CreatedAt.IsZero())
		})
	})

	t.Run("Given invalid identifiers and timestamps", func(t *testing.T) {
		t.Parallel()

		external := &mappings.ExternalMessageResponse{
			Id:             lo.ToPtr("not-a-uuid"),
			ConversationId: lo.ToPtr("not-a-uuid"),
			ParticipantId:  lo.ToPtr("not-a-uuid"),
			CreatedAt:      lo.ToPtr("not-a-time"),
		}

		t.Run("Should map zero values", func(t *testing.T) {
			t.Parallel()

			result := external.ToUseCaseResponse()

			assert.Equal(t, uuid.Nil, result.ID)
			assert.Equal(t, uuid.Nil, result.ConversationId)
			assert.Equal(t, uuid.Nil, result.ParticipantId)
			assert.True(t, result.CreatedAt.IsZero())
		})
	})
}
