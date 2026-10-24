package requests_test

import (
	"strings"
	"testing"

	common_errors "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWhenValidatingAndSanitizingPostBrokerCompletionToSlackRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given an empty conversation ID", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.PostBrokerCompletionToSlackRequest{
				ConversationId:    uuid.Nil,
				CollectionAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
				CollectionTokenId: 1,
				Content:           "hello",
			}

			assert.ErrorIs(t, request.ValidateAndSanitize(), common_errors.ErrValidation)
		})
	})

	t.Run("Given an empty collection address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.PostBrokerCompletionToSlackRequest{
				ConversationId:    uuid.New(),
				CollectionAddress: "",
				CollectionTokenId: 1,
				Content:           "hello",
			}

			assert.ErrorIs(t, request.ValidateAndSanitize(), common_errors.ErrValidation)
		})
	})

	t.Run("Given a collection address without the 0x prefix", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.PostBrokerCompletionToSlackRequest{
				ConversationId:    uuid.New(),
				CollectionAddress: strings.TrimPrefix(strings.ToLower(common_tests.RandomEthAddress().Hex()), "0x"),
				CollectionTokenId: 1,
				Content:           "hello",
			}

			assert.ErrorIs(t, request.ValidateAndSanitize(), common_errors.ErrValidation)
		})
	})

	t.Run("Given a collection address that is not hex", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.PostBrokerCompletionToSlackRequest{
				ConversationId:    uuid.New(),
				CollectionAddress: "0xnot-an-address",
				CollectionTokenId: 1,
				Content:           "hello",
			}

			assert.ErrorIs(t, request.ValidateAndSanitize(), common_errors.ErrValidation)
		})
	})

	t.Run("Given a too short collection address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.PostBrokerCompletionToSlackRequest{
				ConversationId:    uuid.New(),
				CollectionAddress: "0x1234",
				CollectionTokenId: 1,
				Content:           "hello",
			}

			assert.ErrorIs(t, request.ValidateAndSanitize(), common_errors.ErrValidation)
		})
	})

	t.Run("Given an empty content", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.PostBrokerCompletionToSlackRequest{
				ConversationId:    uuid.New(),
				CollectionAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
				CollectionTokenId: 1,
				Content:           "",
			}

			assert.ErrorIs(t, request.ValidateAndSanitize(), common_errors.ErrValidation)
		})
	})

	t.Run("Given a whitespace-only content", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.PostBrokerCompletionToSlackRequest{
				ConversationId:    uuid.New(),
				CollectionAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
				CollectionTokenId: 1,
				Content:           "   ",
			}

			assert.ErrorIs(t, request.ValidateAndSanitize(), common_errors.ErrValidation)
		})
	})

	t.Run("Given a checksummed collection address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should lowercase it in place", func(t *testing.T) {
			t.Parallel()

			address := common_tests.RandomEthAddress().Hex()
			request := &requests.PostBrokerCompletionToSlackRequest{
				ConversationId:    uuid.New(),
				CollectionAddress: address,
				CollectionTokenId: 1,
				Content:           "hello",
			}

			assert.NoError(t, request.ValidateAndSanitize())
			assert.Equal(t, strings.ToLower(address), request.CollectionAddress)
		})
	})

	t.Run("Given a valid request", func(t *testing.T) {
		t.Parallel()

		t.Run("Should accept it", func(t *testing.T) {
			t.Parallel()

			request := &requests.PostBrokerCompletionToSlackRequest{
				ConversationId:    uuid.New(),
				CollectionAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
				CollectionTokenId: 1,
				Content:           "hello",
			}

			assert.NoError(t, request.ValidateAndSanitize())
		})
	})
}
