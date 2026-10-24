package requests_test

import (
	"strings"
	"testing"

	common_errors "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/stretchr/testify/assert"
)

func TestWhenValidatingAndSanitizingRemoveAgentFromChannelRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given an empty collection address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.RemoveAgentFromChannelRequest{
				AgentCommandArguments: requests.AgentCommandArguments{CollectionAddress: "", CollectionTokenId: 1},
			}

			var validationError *common_errors.ValidationError
			assert.ErrorAs(t, request.ValidateAndSanitize(), &validationError)
			assert.Equal(t, "collectionAddress", validationError.Field)
		})
	})

	t.Run("Given a collection address without the 0x prefix", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.RemoveAgentFromChannelRequest{
				AgentCommandArguments: requests.AgentCommandArguments{
					CollectionAddress: strings.TrimPrefix(strings.ToLower(common_tests.RandomEthAddress().Hex()), "0x"),
					CollectionTokenId: 1,
				},
			}

			var validationError *common_errors.ValidationError
			assert.ErrorAs(t, request.ValidateAndSanitize(), &validationError)
			assert.Equal(t, "collectionAddress", validationError.Field)
		})
	})

	t.Run("Given a collection address that is not hex", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.RemoveAgentFromChannelRequest{
				AgentCommandArguments: requests.AgentCommandArguments{
					CollectionAddress: "not-an-address",
					CollectionTokenId: 1,
				},
			}

			var validationError *common_errors.ValidationError
			assert.ErrorAs(t, request.ValidateAndSanitize(), &validationError)
			assert.Equal(t, "collectionAddress", validationError.Field)
		})
	})

	t.Run("Given a too short collection address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.RemoveAgentFromChannelRequest{
				AgentCommandArguments: requests.AgentCommandArguments{
					CollectionAddress: "0x1234",
					CollectionTokenId: 1,
				},
			}

			var validationError *common_errors.ValidationError
			assert.ErrorAs(t, request.ValidateAndSanitize(), &validationError)
			assert.Equal(t, "collectionAddress", validationError.Field)
		})
	})

	t.Run("Given a zero collection token ID", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.RemoveAgentFromChannelRequest{
				AgentCommandArguments: requests.AgentCommandArguments{
					CollectionAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
				},
			}

			var validationError *common_errors.ValidationError
			assert.ErrorAs(t, request.ValidateAndSanitize(), &validationError)
			assert.Equal(t, "collectionTokenId", validationError.Field)
		})
	})

	t.Run("Given a checksummed collection address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should lowercase it in place", func(t *testing.T) {
			t.Parallel()

			address := common_tests.RandomEthAddress().Hex()
			request := &requests.RemoveAgentFromChannelRequest{
				AgentCommandArguments: requests.AgentCommandArguments{CollectionAddress: address, CollectionTokenId: 1},
			}

			assert.NoError(t, request.ValidateAndSanitize())
			assert.Equal(t, strings.ToLower(address), request.CollectionAddress)
		})
	})

	t.Run("Given a hex collection address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should accept it", func(t *testing.T) {
			t.Parallel()

			request := &requests.RemoveAgentFromChannelRequest{
				AgentCommandArguments: requests.AgentCommandArguments{
					CollectionAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
					CollectionTokenId: 1,
				},
			}

			assert.NoError(t, request.ValidateAndSanitize())
		})
	})
}
