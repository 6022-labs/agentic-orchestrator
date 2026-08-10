package models_test

import (
	"encoding/json"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_agent_identity_registry/models"
	"github.com/stretchr/testify/assert"
)

// The registry serves role inside attributes and the wallet inside addresses; it has no
// top-level field for either, so the mapping has to lift them out.
const registryAgentPayload = `{
	"id": "11111111-2222-3333-4444-555555555555",
	"name": "hermes",
	"ensDomain": "hermes.agents6022.eth",
	"attributes": {"role": "clone", "mood": "curious"},
	"addresses": {"evm": "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", "btc": "bc1qexample"},
	"collectionAddress": "0xabc",
	"collectionTokenId": 1
}`

func TestWhenMappingExternalAgentResponse(t *testing.T) {
	t.Parallel()

	t.Run("Given a response as the identity registry actually serves it", func(t *testing.T) {
		t.Parallel()

		decode := func(t *testing.T) *models.ExternalAgentResponse {
			var response models.ExternalAgentResponse
			assert.NoError(t, json.Unmarshal([]byte(registryAgentPayload), &response))
			return &response
		}

		t.Run("Should pass the attributes through, role included", func(t *testing.T) {
			t.Parallel()

			agent := decode(t).ToRegistryAgent()

			assert.Equal(t, enums.ParticipantRoleClone, enums.ParticipantRoleFromAttributes(agent.Attributes))
			assert.Equal(t, "curious", agent.Attributes["mood"])
		})

		t.Run("Should carry every published address, not just the EVM one", func(t *testing.T) {
			t.Parallel()

			agent := decode(t).ToRegistryAgent()

			assert.Equal(t, "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", agent.EvmAddress())
			assert.Equal(t, "bc1qexample", agent.Addresses["btc"])
		})

		t.Run("Should carry the identity fields through unchanged", func(t *testing.T) {
			t.Parallel()

			agent := decode(t).ToRegistryAgent()

			assert.Equal(t, "hermes", agent.Name)
			assert.Equal(t, "hermes.agents6022.eth", agent.EnsDomain)
			assert.Equal(t, "0xabc", agent.CollectionAddress)
			if assert.NotNil(t, agent.CollectionTokenId) {
				assert.Equal(t, uint64(1), *agent.CollectionTokenId)
			}
		})
	})

	t.Run("Given an agent with no attributes or addresses", func(t *testing.T) {
		t.Parallel()

		t.Run("Should map them empty instead of panicking on nil maps", func(t *testing.T) {
			t.Parallel()

			var response models.ExternalAgentResponse
			assert.NoError(t, json.Unmarshal([]byte(`{"name":"hermes"}`), &response))

			agent := response.ToRegistryAgent()

			assert.Empty(t, agent.EvmAddress())
			// Role is optional on-chain, so an absent attribute is not an error.
			assert.Equal(t, enums.ParticipantRoleExpert, enums.ParticipantRoleFromAttributes(agent.Attributes))
		})
	})

	t.Run("Given a nil response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should map to nil", func(t *testing.T) {
			t.Parallel()

			var response *models.ExternalAgentResponse
			assert.Nil(t, response.ToRegistryAgent())
		})
	})
}
