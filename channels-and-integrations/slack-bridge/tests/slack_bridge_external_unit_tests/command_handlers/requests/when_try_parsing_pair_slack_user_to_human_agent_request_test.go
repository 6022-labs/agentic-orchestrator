package requests_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/stretchr/testify/assert"
)

func TestWhenTryParsingPairSlackUserToHumanAgentRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given an empty text", func(t *testing.T) {
		t.Parallel()

		t.Run("Should leave the fields empty", func(t *testing.T) {
			t.Parallel()

			request := &requests.PairSlackUserToHumanAgentRequest{}

			request.TryParse("")

			assert.Empty(t, request.SlackUserId)
			assert.Empty(t, request.HumanAgentCollectionAddress)
			assert.Zero(t, request.HumanAgentCollectionTokenId)
		})
	})

	t.Run("Given the token ID is missing", func(t *testing.T) {
		t.Parallel()

		t.Run("Should leave the fields empty", func(t *testing.T) {
			t.Parallel()

			request := &requests.PairSlackUserToHumanAgentRequest{}

			request.TryParse("U654321 " + common_tests.RandomEthAddress().Hex())

			assert.Empty(t, request.SlackUserId)
			assert.Empty(t, request.HumanAgentCollectionAddress)
			assert.Zero(t, request.HumanAgentCollectionTokenId)
		})
	})

	t.Run("Given too many arguments", func(t *testing.T) {
		t.Parallel()

		t.Run("Should leave the fields empty", func(t *testing.T) {
			t.Parallel()

			request := &requests.PairSlackUserToHumanAgentRequest{}

			request.TryParse("U654321 " + common_tests.RandomEthAddress().Hex() + " 1 2")

			assert.Empty(t, request.SlackUserId)
			assert.Empty(t, request.HumanAgentCollectionAddress)
			assert.Zero(t, request.HumanAgentCollectionTokenId)
		})
	})

	t.Run("Given the token ID is not a number", func(t *testing.T) {
		t.Parallel()

		t.Run("Should leave the token ID empty", func(t *testing.T) {
			t.Parallel()

			request := &requests.PairSlackUserToHumanAgentRequest{}

			request.TryParse("U654321 " + common_tests.RandomEthAddress().Hex() + " abc")

			assert.Equal(t, "U654321", request.SlackUserId)
			assert.Zero(t, request.HumanAgentCollectionTokenId)
		})
	})

	t.Run("Given the token ID is negative", func(t *testing.T) {
		t.Parallel()

		t.Run("Should leave the token ID empty", func(t *testing.T) {
			t.Parallel()

			request := &requests.PairSlackUserToHumanAgentRequest{}

			request.TryParse("U654321 " + common_tests.RandomEthAddress().Hex() + " -1")

			assert.Equal(t, "U654321", request.SlackUserId)
			assert.Zero(t, request.HumanAgentCollectionTokenId)
		})
	})

	t.Run("Given the token ID overflows uint64", func(t *testing.T) {
		t.Parallel()

		t.Run("Should leave the token ID empty", func(t *testing.T) {
			t.Parallel()

			request := &requests.PairSlackUserToHumanAgentRequest{}

			request.TryParse("U654321 " + common_tests.RandomEthAddress().Hex() + " 18446744073709551616")

			assert.Equal(t, "U654321", request.SlackUserId)
			assert.Zero(t, request.HumanAgentCollectionTokenId)
		})
	})

	t.Run("Given a slack user ID, a checksummed collection address and a token ID", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fill the fields as typed", func(t *testing.T) {
			t.Parallel()

			address := common_tests.RandomEthAddress().Hex()
			request := &requests.PairSlackUserToHumanAgentRequest{}

			request.TryParse("U654321 " + address + " 42")

			assert.Equal(t, "U654321", request.SlackUserId)
			assert.Equal(t, address, request.HumanAgentCollectionAddress)
			assert.Equal(t, uint64(42), request.HumanAgentCollectionTokenId)
		})
	})
}
