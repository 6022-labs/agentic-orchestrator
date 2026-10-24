package requests_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/stretchr/testify/assert"
)

func WhenValidatingAndSanitizingVerifyWalletForBridgeRequestBeforeEach() *requests.VerifyWalletForBridgeRequest {
	return &requests.VerifyWalletForBridgeRequest{
		Platform:  enums.BridgePlatformSlack,
		Challenge: "challenge",
		Signature: "0xsignature",
	}
}

func TestWhenValidatingAndSanitizingVerifyWalletForBridgeRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given the platform is empty", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error on platform", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingVerifyWalletForBridgeRequestBeforeEach()
			request.Platform = ""

			err := request.ValidateAndSanitize()

			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "platform", validationErr.Field)
		})
	})

	t.Run("Given the platform is unknown", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error listing the allowed platforms", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingVerifyWalletForBridgeRequestBeforeEach()
			request.Platform = "discord"

			err := request.ValidateAndSanitize()

			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "platform", validationErr.Field)
			assert.Contains(t, validationErr.Message, enums.GetAllBridgePlatforms().String())
		})
	})

	t.Run("Given the platform is invalid and the other fields are empty", func(t *testing.T) {
		t.Parallel()

		t.Run("Should report the platform first", func(t *testing.T) {
			t.Parallel()
			request := &requests.VerifyWalletForBridgeRequest{Platform: "discord"}

			err := request.ValidateAndSanitize()

			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "platform", validationErr.Field)
		})
	})

	t.Run("Given the challenge is empty", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error on challenge", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingVerifyWalletForBridgeRequestBeforeEach()
			request.Challenge = ""

			err := request.ValidateAndSanitize()

			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "challenge", validationErr.Field)
		})
	})

	t.Run("Given the challenge is only whitespace", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error on challenge", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingVerifyWalletForBridgeRequestBeforeEach()
			request.Challenge = " \t\n"

			err := request.ValidateAndSanitize()

			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "challenge", validationErr.Field)
		})
	})

	t.Run("Given the signature is empty", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error on signature", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingVerifyWalletForBridgeRequestBeforeEach()
			request.Signature = ""

			err := request.ValidateAndSanitize()

			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "signature", validationErr.Field)
		})
	})

	t.Run("Given the signature is only whitespace", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error on signature", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingVerifyWalletForBridgeRequestBeforeEach()
			request.Signature = "   "

			err := request.ValidateAndSanitize()

			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "signature", validationErr.Field)
		})
	})

	t.Run("Given a valid request", func(t *testing.T) {
		t.Parallel()

		t.Run("Should succeed and leave the fields untouched", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingVerifyWalletForBridgeRequestBeforeEach()

			err := request.ValidateAndSanitize()

			assert.NoError(t, err)
			assert.Equal(t, WhenValidatingAndSanitizingVerifyWalletForBridgeRequestBeforeEach(), request)
		})
	})
}
