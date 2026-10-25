package requests_test

import (
	"strings"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/stretchr/testify/assert"
)

func WhenValidatingAndSanitizingSignedRequestBeforeEach() *requests.AuthenticateSignedRequestRequest {
	return &requests.AuthenticateSignedRequestRequest{
		Method:    "GET",
		Path:      "/agents",
		Timestamp: "1792944000",
		Nonce:     "0123456789abcdef",
		Signature: "0x" + strings.Repeat("ab", 65),
	}
}

func TestWhenValidatingAndSanitizingSignedRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given a timestamp that is not a number", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error on timestamp", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingSignedRequestBeforeEach()
			request.Timestamp = "yesterday"

			err := request.ValidateAndSanitize()

			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "timestamp", validationErr.Field)
		})
	})

	t.Run("Given a missing timestamp", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error on timestamp", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingSignedRequestBeforeEach()
			request.Timestamp = ""

			err := request.ValidateAndSanitize()

			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "timestamp", validationErr.Field)
		})
	})

	t.Run("Given a non-positive timestamp", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error on timestamp", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingSignedRequestBeforeEach()
			request.Timestamp = "0"

			err := request.ValidateAndSanitize()

			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "timestamp", validationErr.Field)
		})
	})

	t.Run("Given a nonce shorter than 16 characters", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error on nonce", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingSignedRequestBeforeEach()
			request.Nonce = "short"

			err := request.ValidateAndSanitize()

			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "nonce", validationErr.Field)
		})
	})

	t.Run("Given a nonce with a separator character", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error on nonce", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingSignedRequestBeforeEach()
			request.Nonce = "0123456789abcdef/x"

			err := request.ValidateAndSanitize()

			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "nonce", validationErr.Field)
		})
	})

	t.Run("Given a signature that is not 65 bytes of hex", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error on signature", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingSignedRequestBeforeEach()
			request.Signature = "0xabcd"

			err := request.ValidateAndSanitize()

			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "signature", validationErr.Field)
		})
	})

	t.Run("Given fields with surrounding whitespace", func(t *testing.T) {
		t.Parallel()

		t.Run("Should trim them in place", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingSignedRequestBeforeEach()
			signature := request.Signature
			request.Timestamp = " 1792944000 "
			request.Nonce = " 0123456789abcdef "
			request.Signature = " " + signature + " "

			err := request.ValidateAndSanitize()

			assert.Nil(t, err)
			assert.Equal(t, "1792944000", request.Timestamp)
			assert.Equal(t, "0123456789abcdef", request.Nonce)
			assert.Equal(t, signature, request.Signature)
		})
	})
}
