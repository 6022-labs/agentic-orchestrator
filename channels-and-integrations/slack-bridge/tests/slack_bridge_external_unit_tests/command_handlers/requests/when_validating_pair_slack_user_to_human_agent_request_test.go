package requests_test

import (
	"strings"
	"testing"

	common_errors "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/stretchr/testify/assert"
)

func TestWhenValidatingAndSanitizingPairSlackUserToHumanAgentRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given an empty slack user ID", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.PairSlackUserToHumanAgentRequest{
				HumanAgentCollectionAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
				HumanAgentCollectionTokenId: 1,
			}

			assert.ErrorIs(t, request.ValidateAndSanitize(), common_errors.ErrValidation)
		})
	})

	t.Run("Given a collection address without the 0x prefix", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.PairSlackUserToHumanAgentRequest{
				SlackUserId: "U654321",
				HumanAgentCollectionAddress: strings.TrimPrefix(
					strings.ToLower(common_tests.RandomEthAddress().Hex()),
					"0x",
				),
				HumanAgentCollectionTokenId: 1,
			}

			assert.ErrorIs(t, request.ValidateAndSanitize(), common_errors.ErrValidation)
		})
	})

	t.Run("Given a collection address that is not hex", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.PairSlackUserToHumanAgentRequest{
				SlackUserId:                 "U654321",
				HumanAgentCollectionAddress: "not-an-address",
				HumanAgentCollectionTokenId: 1,
			}

			assert.ErrorIs(t, request.ValidateAndSanitize(), common_errors.ErrValidation)
		})
	})

	t.Run("Given a zero collection token ID", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.PairSlackUserToHumanAgentRequest{
				SlackUserId:                 "U654321",
				HumanAgentCollectionAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
			}

			assert.ErrorIs(t, request.ValidateAndSanitize(), common_errors.ErrValidation)
		})
	})

	t.Run("Given a checksummed collection address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should lowercase it in place", func(t *testing.T) {
			t.Parallel()

			address := common_tests.RandomEthAddress().Hex()
			request := &requests.PairSlackUserToHumanAgentRequest{
				SlackUserId:                 "U654321",
				HumanAgentCollectionAddress: address,
				HumanAgentCollectionTokenId: 1,
			}

			assert.NoError(t, request.ValidateAndSanitize())
			assert.Equal(t, strings.ToLower(address), request.HumanAgentCollectionAddress)
		})
	})

	t.Run("Given a slack user ID and a hex collection address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should accept it", func(t *testing.T) {
			t.Parallel()

			request := &requests.PairSlackUserToHumanAgentRequest{
				SlackUserId:                 "U654321",
				HumanAgentCollectionAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
				HumanAgentCollectionTokenId: 1,
			}

			assert.NoError(t, request.ValidateAndSanitize())
		})
	})
}
