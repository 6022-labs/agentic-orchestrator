package requests_test

import (
	"strings"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/stretchr/testify/assert"
)

const registerAgentAddress = "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed"

func TestWhenValidatingAndSanitizingRegisterAgentRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given a collection address without the 0x prefix", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.RegisterAgentRequest{
				CollectionAddress: strings.TrimPrefix(registerAgentAddress, "0x"),
				CollectionTokenId: 1,
			}

			err := request.ValidateAndSanitize()

			assert.Error(t, err)
			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given a collection address that is not hex", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.RegisterAgentRequest{
				CollectionAddress: "0xZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ",
				CollectionTokenId: 1,
			}

			err := request.ValidateAndSanitize()

			assert.Error(t, err)
			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given a too short collection address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.RegisterAgentRequest{CollectionAddress: "0x1234", CollectionTokenId: 1}

			err := request.ValidateAndSanitize()

			assert.Error(t, err)
			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given an empty collection address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.RegisterAgentRequest{CollectionAddress: "", CollectionTokenId: 1}

			err := request.ValidateAndSanitize()

			assert.Error(t, err)
			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given a zero collection token ID", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.RegisterAgentRequest{CollectionAddress: registerAgentAddress}

			err := request.ValidateAndSanitize()

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "collectionTokenId")
		})
	})

	t.Run("Given a checksummed collection address surrounded by whitespace", func(t *testing.T) {
		t.Parallel()

		t.Run("Should trim and lowercase it in place", func(t *testing.T) {
			t.Parallel()

			request := &requests.RegisterAgentRequest{
				CollectionAddress: "  " + registerAgentAddress + "\t",
				CollectionTokenId: 1,
			}

			assert.NoError(t, request.ValidateAndSanitize())
			assert.Equal(t, strings.ToLower(registerAgentAddress), request.CollectionAddress)
		})
	})

	t.Run("Given a valid request", func(t *testing.T) {
		t.Parallel()

		t.Run("Should accept it", func(t *testing.T) {
			t.Parallel()

			request := &requests.RegisterAgentRequest{
				CollectionAddress: strings.ToLower(registerAgentAddress),
				CollectionTokenId: 1,
			}

			assert.NoError(t, request.ValidateAndSanitize())
		})
	})
}
