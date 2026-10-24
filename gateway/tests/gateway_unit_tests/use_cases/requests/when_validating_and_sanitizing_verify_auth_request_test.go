package requests_test

import (
	"strings"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/stretchr/testify/assert"
)

func WhenValidatingAndSanitizingVerifyAuthRequestBeforeEach() *requests.VerifyAuthRequest {
	return &requests.VerifyAuthRequest{
		Address:   checksummedAddress,
		Signature: "0x" + strings.Repeat("ab", 65),
	}
}

func TestWhenValidatingAndSanitizingVerifyAuthRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given an address without the 0x prefix", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error on address", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingVerifyAuthRequestBeforeEach()
			request.Address = strings.TrimPrefix(checksummedAddress, "0x")

			err := request.ValidateAndSanitize()

			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "address", validationErr.Field)
		})
	})

	t.Run("Given an address that is not hex", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error on address", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingVerifyAuthRequestBeforeEach()
			request.Address = nonHexAddress

			err := request.ValidateAndSanitize()

			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "address", validationErr.Field)
		})
	})

	t.Run("Given a checksummed address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should lowercase it in place", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingVerifyAuthRequestBeforeEach()

			err := request.ValidateAndSanitize()

			assert.NoError(t, err)
			assert.Equal(t, lowercasedAddress, request.Address)
		})
	})
}
