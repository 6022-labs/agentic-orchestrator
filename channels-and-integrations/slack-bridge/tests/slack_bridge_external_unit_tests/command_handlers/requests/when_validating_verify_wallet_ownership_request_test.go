package requests_test

import (
	"testing"

	common_errors "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/stretchr/testify/assert"
)

func TestWhenValidatingVerifyWalletOwnershipRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given an empty user ID", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.VerifyWalletOwnershipRequest{UserId: ""}

			assert.ErrorIs(t, request.Validate(), common_errors.ErrValidation)
		})
	})

	t.Run("Given a blank user ID", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.VerifyWalletOwnershipRequest{UserId: "   "}

			assert.ErrorIs(t, request.Validate(), common_errors.ErrValidation)
		})
	})

	t.Run("Given a user ID", func(t *testing.T) {
		t.Parallel()

		t.Run("Should accept it", func(t *testing.T) {
			t.Parallel()

			request := &requests.VerifyWalletOwnershipRequest{UserId: "U123456"}

			assert.NoError(t, request.Validate())
		})
	})
}
