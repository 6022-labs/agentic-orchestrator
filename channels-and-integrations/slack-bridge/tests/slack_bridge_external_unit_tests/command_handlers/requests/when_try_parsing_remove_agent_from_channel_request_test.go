package requests_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/stretchr/testify/assert"
)

func TestWhenTryParsingRemoveAgentFromChannelRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given an empty text", func(t *testing.T) {
		t.Parallel()

		t.Run("Should leave the fields empty", func(t *testing.T) {
			t.Parallel()

			request := &requests.RemoveAgentFromChannelRequest{}

			request.TryParse("")

			assert.Empty(t, request.CollectionAddress)
			assert.Zero(t, request.CollectionTokenId)
		})
	})

	t.Run("Given the token ID is missing", func(t *testing.T) {
		t.Parallel()

		t.Run("Should leave the fields empty", func(t *testing.T) {
			t.Parallel()

			request := &requests.RemoveAgentFromChannelRequest{}

			request.TryParse(common_tests.RandomEthAddress().Hex())

			assert.Empty(t, request.CollectionAddress)
			assert.Zero(t, request.CollectionTokenId)
		})
	})

	t.Run("Given too many arguments", func(t *testing.T) {
		t.Parallel()

		t.Run("Should leave the fields empty", func(t *testing.T) {
			t.Parallel()

			request := &requests.RemoveAgentFromChannelRequest{}

			request.TryParse(common_tests.RandomEthAddress().Hex() + " 1 2")

			assert.Empty(t, request.CollectionAddress)
			assert.Zero(t, request.CollectionTokenId)
		})
	})

	t.Run("Given the token ID is not a number", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fill the address and leave the token ID empty", func(t *testing.T) {
			t.Parallel()

			address := common_tests.RandomEthAddress().Hex()
			request := &requests.RemoveAgentFromChannelRequest{}

			request.TryParse(address + " abc")

			assert.Equal(t, address, request.CollectionAddress)
			assert.Zero(t, request.CollectionTokenId)
		})
	})

	t.Run("Given the token ID is negative", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fill the address and leave the token ID empty", func(t *testing.T) {
			t.Parallel()

			address := common_tests.RandomEthAddress().Hex()
			request := &requests.RemoveAgentFromChannelRequest{}

			request.TryParse(address + " -1")

			assert.Equal(t, address, request.CollectionAddress)
			assert.Zero(t, request.CollectionTokenId)
		})
	})

	t.Run("Given the token ID overflows uint64", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fill the address and leave the token ID empty", func(t *testing.T) {
			t.Parallel()

			address := common_tests.RandomEthAddress().Hex()
			request := &requests.RemoveAgentFromChannelRequest{}

			request.TryParse(address + " 18446744073709551616")

			assert.Equal(t, address, request.CollectionAddress)
			assert.Zero(t, request.CollectionTokenId)
		})
	})

	t.Run("Given a checksummed collection address and a token ID", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fill the fields as typed", func(t *testing.T) {
			t.Parallel()

			address := common_tests.RandomEthAddress().Hex()
			request := &requests.RemoveAgentFromChannelRequest{}

			request.TryParse(address + " 42")

			assert.Equal(t, address, request.CollectionAddress)
			assert.Equal(t, uint64(42), request.CollectionTokenId)
		})
	})

	t.Run("Given arguments separated by extra whitespace", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fill the fields", func(t *testing.T) {
			t.Parallel()

			address := common_tests.RandomEthAddress().Hex()
			request := &requests.RemoveAgentFromChannelRequest{}

			request.TryParse("  " + address + "   7 ")

			assert.Equal(t, address, request.CollectionAddress)
			assert.Equal(t, uint64(7), request.CollectionTokenId)
		})
	})
}
