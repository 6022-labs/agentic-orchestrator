package mappings_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/mappings"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestWhenMappingExternalListAgentsResponse(t *testing.T) {
	t.Parallel()

	t.Run("Given a fully populated response", func(t *testing.T) {
		t.Parallel()

		agentId := uuid.New()
		external := &mappings.ExternalListAgentsResponse{
			Page:            lo.ToPtr(2),
			ElementsPerPage: lo.ToPtr(10),
			TotalElements:   lo.ToPtr(11),
			Agents: &[]clients.ResponsesAgentResponse{{
				Id:                lo.ToPtr(agentId.String()),
				Name:              lo.ToPtr("Ada"),
				CollectionAddress: lo.ToPtr("0xabc"),
				CollectionTokenId: lo.ToPtr(7),
			}},
		}

		t.Run("Should map every field", func(t *testing.T) {
			t.Parallel()

			result := external.ToUseCaseResponse()

			assert.Equal(t, uint(2), result.Page)
			assert.Equal(t, uint(10), result.ElementsPerPage)
			assert.Equal(t, uint(11), result.TotalElements)
			if assert.Len(t, result.Agents, 1) {
				assert.Equal(t, agentId, result.Agents[0].ID)
				assert.Equal(t, "Ada", result.Agents[0].Name)
				assert.Equal(t, "0xabc", result.Agents[0].CollectionAddress)
				assert.Equal(t, uint64(7), *result.Agents[0].CollectionTokenId)
			}
		})
	})

	t.Run("Given a nil response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			var external *mappings.ExternalListAgentsResponse

			assert.Nil(t, external.ToUseCaseResponse())
		})
	})

	t.Run("Given an empty response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should map zero values", func(t *testing.T) {
			t.Parallel()

			result := (&mappings.ExternalListAgentsResponse{}).ToUseCaseResponse()

			assert.Equal(t, uint(0), result.Page)
			assert.Equal(t, uint(0), result.ElementsPerPage)
			assert.Equal(t, uint(0), result.TotalElements)
			assert.Nil(t, result.Agents)
		})
	})
}
