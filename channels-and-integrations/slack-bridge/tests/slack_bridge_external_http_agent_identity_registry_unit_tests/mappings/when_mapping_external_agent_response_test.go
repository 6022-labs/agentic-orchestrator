package mappings_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_http_agent_identity_registry/mappings"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestWhenMappingExternalAgentResponse(t *testing.T) {
	t.Parallel()

	t.Run("Given a response with a role attribute and an icon image", func(t *testing.T) {
		t.Parallel()

		external := &mappings.ExternalAgentResponse{
			Owner:             lo.ToPtr("0xowner"),
			EnsDomain:         lo.ToPtr("ada.eth"),
			Images:            lo.ToPtr(map[string]string{"icon": "cid123"}),
			Attributes:        lo.ToPtr(map[string]string{"role": "expert"}),
			CollectionAddress: lo.ToPtr("0xabc"),
			CollectionTokenId: lo.ToPtr(7),
		}

		t.Run("Should map every field", func(t *testing.T) {
			t.Parallel()

			result, err := external.ToSlackAgent()

			assert.NoError(t, err)
			assert.Equal(t, enums.AgentRoleExpert, result.Role)
			assert.Equal(t, "0xowner", result.OwnerWallet)
			assert.Equal(t, "ada.eth", result.EnsDomain)
			assert.Equal(t, uint64(7), result.CollectionTokenId)
			assert.Equal(t, "0xabc", result.CollectionAddress)
			assert.Equal(t, "cid123", *result.IconImageCid)
		})
	})

	t.Run("Given a response without an icon image", func(t *testing.T) {
		t.Parallel()

		external := &mappings.ExternalAgentResponse{
			Attributes: lo.ToPtr(map[string]string{"role": "clone"}),
		}

		t.Run("Should leave the icon cid nil", func(t *testing.T) {
			t.Parallel()

			result, err := external.ToSlackAgent()

			assert.NoError(t, err)
			assert.Nil(t, result.IconImageCid)
		})
	})

	t.Run("Given a nil response", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			var external *mappings.ExternalAgentResponse

			result, err := external.ToSlackAgent()

			assert.NoError(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Given a response without a role attribute", func(t *testing.T) {
		t.Parallel()

		external := &mappings.ExternalAgentResponse{
			Attributes: lo.ToPtr(map[string]string{}),
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			result, err := external.ToSlackAgent()

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
}
