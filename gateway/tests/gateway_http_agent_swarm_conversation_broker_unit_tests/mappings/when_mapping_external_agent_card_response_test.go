package mappings_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/mappings"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestWhenMappingExternalAgentCardResponse(t *testing.T) {
	t.Parallel()

	t.Run("Given a fully populated response", func(t *testing.T) {
		t.Parallel()

		external := &mappings.ExternalAgentCardResponse{
			Protocol:        lo.ToPtr("6022"),
			ProtocolVersion: lo.ToPtr("1"),
			Name:            lo.ToPtr("Ada"),
			EnsDomain:       lo.ToPtr("ada.experts.80002.6022.eth"),
			Endpoints: &clients.ResponsesRemoteAgentCardEndpointsResponse{
				Responses:       lo.ToPtr("https://ada.example/responses"),
				ChatCompletions: lo.ToPtr("https://ada.example/chat"),
			},
			PaymentMethods: lo.ToPtr([]string{"x402"}),
		}

		t.Run("Should map every field", func(t *testing.T) {
			t.Parallel()

			result := external.ToUseCaseResponse()

			assert.Equal(t, "6022", result.Protocol)
			assert.Equal(t, "1", result.ProtocolVersion)
			assert.Equal(t, "Ada", result.Name)
			assert.Equal(t, "ada.experts.80002.6022.eth", result.EnsDomain)
			assert.Equal(t, "https://ada.example/responses", result.Endpoints.Responses)
			assert.Equal(t, "https://ada.example/chat", result.Endpoints.ChatCompletions)
			assert.Equal(t, []string{"x402"}, result.PaymentMethods)
		})
	})

	t.Run("Given a nil response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			var external *mappings.ExternalAgentCardResponse

			assert.Nil(t, external.ToUseCaseResponse())
		})
	})

	t.Run("Given an empty response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should map zero values", func(t *testing.T) {
			t.Parallel()

			result := (&mappings.ExternalAgentCardResponse{}).ToUseCaseResponse()

			assert.Empty(t, result.Protocol)
			assert.Empty(t, result.ProtocolVersion)
			assert.Empty(t, result.Name)
			assert.Empty(t, result.EnsDomain)
			assert.Empty(t, result.Endpoints.Responses)
			assert.Empty(t, result.Endpoints.ChatCompletions)
			assert.Nil(t, result.PaymentMethods)
		})
	})
}
