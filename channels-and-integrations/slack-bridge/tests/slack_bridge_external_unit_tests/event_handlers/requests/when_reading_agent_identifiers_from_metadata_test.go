package requests_test

import (
	"strings"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/stretchr/testify/assert"
)

func TestWhenReadingAgentIdentifiersFromMetadata(t *testing.T) {
	t.Parallel()

	t.Run("Given a collection token address without the 0x prefix", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			metadata := &requests.MessageCreatedMetadata{
				EventPayload: map[string]any{
					"collection_token_address": strings.TrimPrefix(common_tests.RandomEthAddress().Hex(), "0x"),
					"collection_token_id":      float64(7),
				},
			}

			_, _, err := metadata.ToAgentIdentifiers()

			assert.Error(t, err)
		})
	})

	t.Run("Given a collection token address that is not hex", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			metadata := &requests.MessageCreatedMetadata{
				EventPayload: map[string]any{
					"collection_token_address": "not-an-address",
					"collection_token_id":      float64(7),
				},
			}

			_, _, err := metadata.ToAgentIdentifiers()

			assert.Error(t, err)
		})
	})

	t.Run("Given a checksummed collection token address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return it lowercased with the token ID", func(t *testing.T) {
			t.Parallel()

			address := common_tests.RandomEthAddress().Hex()
			metadata := &requests.MessageCreatedMetadata{
				EventPayload: map[string]any{
					"collection_token_address": address,
					"collection_token_id":      float64(7),
				},
			}

			collectionAddress, collectionTokenId, err := metadata.ToAgentIdentifiers()

			assert.NoError(t, err)
			assert.Equal(t, strings.ToLower(address), *collectionAddress)
			assert.Equal(t, uint64(7), *collectionTokenId)
		})
	})
}
