package mappings_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/mappings"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestWhenMappingExternalAgentCollectionResponse(t *testing.T) {
	t.Parallel()

	t.Run("Given a fully populated response", func(t *testing.T) {
		t.Parallel()

		external := &mappings.ExternalAgentCollectionResponse{
			Address: lo.ToPtr("0xabc"),
			Name:    lo.ToPtr("Experts"),
			Symbol:  lo.ToPtr("EXP"),
		}

		t.Run("Should map every field", func(t *testing.T) {
			t.Parallel()

			result := external.ToUseCaseResponse()

			assert.Equal(t, "0xabc", result.Address)
			assert.Equal(t, "Experts", result.Name)
			assert.Equal(t, "EXP", result.Symbol)
		})
	})

	t.Run("Given a nil response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			var external *mappings.ExternalAgentCollectionResponse

			assert.Nil(t, external.ToUseCaseResponse())
		})
	})

	t.Run("Given an empty response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should map zero values", func(t *testing.T) {
			t.Parallel()

			result := (&mappings.ExternalAgentCollectionResponse{}).ToUseCaseResponse()

			assert.Empty(t, result.Address)
			assert.Empty(t, result.Name)
			assert.Empty(t, result.Symbol)
		})
	})
}
