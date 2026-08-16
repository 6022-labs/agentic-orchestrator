package mappings_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/mappings"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestWhenMappingExternalListConversationsResponse(t *testing.T) {
	t.Parallel()

	t.Run("Given a fully populated response", func(t *testing.T) {
		t.Parallel()

		conversationId := uuid.New()
		external := &mappings.ExternalListConversationsResponse{
			Page:            lo.ToPtr(1),
			ElementsPerPage: lo.ToPtr(20),
			TotalElements:   lo.ToPtr(1),
			Conversations: &[]clients.ResponsesConversationResponse{{
				Id:     lo.ToPtr(conversationId.String()),
				Status: lo.ToPtr(clients.EnumsConversationStatus("ongoing")),
			}},
		}

		t.Run("Should map every field", func(t *testing.T) {
			t.Parallel()

			result := external.ToUseCaseResponse()

			assert.Equal(t, uint(1), result.Page)
			assert.Equal(t, uint(20), result.ElementsPerPage)
			assert.Equal(t, uint(1), result.TotalElements)
			if assert.Len(t, result.Conversations, 1) {
				assert.Equal(t, conversationId, result.Conversations[0].ID)
				assert.Equal(t, "ongoing", result.Conversations[0].Status)
			}
		})
	})

	t.Run("Given a nil response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			var external *mappings.ExternalListConversationsResponse

			assert.Nil(t, external.ToUseCaseResponse())
		})
	})

	t.Run("Given an empty response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should map zero values", func(t *testing.T) {
			t.Parallel()

			result := (&mappings.ExternalListConversationsResponse{}).ToUseCaseResponse()

			assert.Equal(t, uint(0), result.Page)
			assert.Equal(t, uint(0), result.ElementsPerPage)
			assert.Equal(t, uint(0), result.TotalElements)
			assert.Nil(t, result.Conversations)
		})
	})
}
