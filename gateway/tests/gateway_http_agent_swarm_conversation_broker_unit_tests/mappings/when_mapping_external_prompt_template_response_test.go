package mappings_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/mappings"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestWhenMappingExternalPromptTemplateResponse(t *testing.T) {
	t.Parallel()

	t.Run("Given a fully populated response", func(t *testing.T) {
		t.Parallel()

		external := &mappings.ExternalPromptTemplateResponse{
			CollectionAddress: lo.ToPtr("0xabc"),
			CollectionTokenId: lo.ToPtr(7),
			Key:               lo.ToPtr(clients.EnumsAgentPromptTemplateKey("clone-system")),
			Content:           lo.ToPtr("You are {{name}}"),
		}

		t.Run("Should map every field", func(t *testing.T) {
			t.Parallel()

			result := external.ToUseCaseResponse()

			assert.Equal(t, "0xabc", result.CollectionAddress)
			assert.Equal(t, uint64(7), result.CollectionTokenId)
			assert.Equal(t, "clone-system", result.Key)
			assert.Equal(t, "You are {{name}}", result.Content)
		})
	})

	t.Run("Given a nil response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			var external *mappings.ExternalPromptTemplateResponse

			assert.Nil(t, external.ToUseCaseResponse())
		})
	})

	t.Run("Given an empty response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should map zero values", func(t *testing.T) {
			t.Parallel()

			result := (&mappings.ExternalPromptTemplateResponse{}).ToUseCaseResponse()

			assert.Empty(t, result.CollectionAddress)
			assert.Equal(t, uint64(0), result.CollectionTokenId)
			assert.Empty(t, result.Key)
			assert.Empty(t, result.Content)
		})
	})
}
