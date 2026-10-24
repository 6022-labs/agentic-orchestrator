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

func TestWhenValidatingAndSanitizingSetWalletOwnershipToSlackUserRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given an empty challenge", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.SetWalletOwnershipToSlackUserRequest{
				Challenge:     uuid.Nil,
				WalletAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
			}

			assert.ErrorIs(t, request.ValidateAndSanitize(), common_errors.ErrValidation)
		})
	})

	t.Run("Given an empty wallet address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.SetWalletOwnershipToSlackUserRequest{
				Challenge:     uuid.New(),
				WalletAddress: "",
			}

			assert.ErrorIs(t, request.ValidateAndSanitize(), common_errors.ErrValidation)
		})
	})

	t.Run("Given a wallet address without the 0x prefix", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.SetWalletOwnershipToSlackUserRequest{
				Challenge:     uuid.New(),
				WalletAddress: strings.TrimPrefix(strings.ToLower(common_tests.RandomEthAddress().Hex()), "0x"),
			}

			assert.ErrorIs(t, request.ValidateAndSanitize(), common_errors.ErrValidation)
		})
	})

	t.Run("Given a wallet address that is not hex", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.SetWalletOwnershipToSlackUserRequest{
				Challenge:     uuid.New(),
				WalletAddress: "0xnot-an-address",
			}

			assert.ErrorIs(t, request.ValidateAndSanitize(), common_errors.ErrValidation)
		})
	})

	t.Run("Given a too short wallet address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.SetWalletOwnershipToSlackUserRequest{
				Challenge:     uuid.New(),
				WalletAddress: "0x1234",
			}

			assert.ErrorIs(t, request.ValidateAndSanitize(), common_errors.ErrValidation)
		})
	})

	t.Run("Given a checksummed wallet address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should lowercase it in place", func(t *testing.T) {
			t.Parallel()

			address := common_tests.RandomEthAddress().Hex()
			request := &requests.SetWalletOwnershipToSlackUserRequest{
				Challenge:     uuid.New(),
				WalletAddress: address,
			}

			assert.NoError(t, request.ValidateAndSanitize())
			assert.Equal(t, strings.ToLower(address), request.WalletAddress)
		})
	})

	t.Run("Given a valid request", func(t *testing.T) {
		t.Parallel()

		t.Run("Should accept it", func(t *testing.T) {
			t.Parallel()

			request := &requests.SetWalletOwnershipToSlackUserRequest{
				Challenge:     uuid.New(),
				WalletAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
			}

			assert.NoError(t, request.ValidateAndSanitize())
		})
	})
}
