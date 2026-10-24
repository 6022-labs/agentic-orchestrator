package requests_test

import (
	"testing"

	common_errors "github.com/6022-labs/agentic-orchestrator/agent-identity/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/requests"
	"github.com/stretchr/testify/assert"
)

type WhenValidatingAndSanitizingAgentByCollectionTokenIdRequestTestingSuite struct {
	sut *requests.AgentByCollectionTokenIdRequest
}

func WhenValidatingAndSanitizingAgentByCollectionTokenIdRequestBeforeEach(
	agentCollectionAddress string,
) *WhenValidatingAndSanitizingAgentByCollectionTokenIdRequestTestingSuite {
	return &WhenValidatingAndSanitizingAgentByCollectionTokenIdRequestTestingSuite{
		sut: &requests.AgentByCollectionTokenIdRequest{
			AgentCollectionAddress: agentCollectionAddress,
			AgentCollectionTokenId: 1,
		},
	}
}

func TestWhenValidatingAndSanitizingAgentByCollectionTokenIdRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given a collection address without the 0x prefix", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error and leave the address untouched", func(t *testing.T) {
			t.Parallel()
			address := "AbCdEf0000000000000000000000000000000001"
			suite := WhenValidatingAndSanitizingAgentByCollectionTokenIdRequestBeforeEach(address)

			err := suite.sut.ValidateAndSanitize()

			var validationError *common_errors.ValidationError
			assert.ErrorAs(t, err, &validationError)
			assert.Equal(t, "agentCollectionAddress", validationError.Field)
			assert.Equal(t, address, suite.sut.AgentCollectionAddress)
		})
	})

	t.Run("Given an empty collection address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error and leave the address untouched", func(t *testing.T) {
			t.Parallel()
			address := ""
			suite := WhenValidatingAndSanitizingAgentByCollectionTokenIdRequestBeforeEach(address)

			err := suite.sut.ValidateAndSanitize()

			var validationError *common_errors.ValidationError
			assert.ErrorAs(t, err, &validationError)
			assert.Equal(t, "agentCollectionAddress", validationError.Field)
			assert.Equal(t, address, suite.sut.AgentCollectionAddress)
		})
	})

	t.Run("Given a collection address that is not hex", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error and leave the address untouched", func(t *testing.T) {
			t.Parallel()
			address := "0xabcdef000000000000000000000000000000000z"
			suite := WhenValidatingAndSanitizingAgentByCollectionTokenIdRequestBeforeEach(address)

			err := suite.sut.ValidateAndSanitize()

			var validationError *common_errors.ValidationError
			assert.ErrorAs(t, err, &validationError)
			assert.Equal(t, "agentCollectionAddress", validationError.Field)
			assert.Equal(t, address, suite.sut.AgentCollectionAddress)
		})
	})

	t.Run("Given a zero collection token ID", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenValidatingAndSanitizingAgentByCollectionTokenIdRequestBeforeEach(
				"0xabcdef0000000000000000000000000000000001",
			)
			suite.sut.AgentCollectionTokenId = 0

			err := suite.sut.ValidateAndSanitize()

			var validationError *common_errors.ValidationError
			assert.ErrorAs(t, err, &validationError)
			assert.Equal(t, "agentCollectionTokenId", validationError.Field)
		})
	})

	t.Run("Given a checksummed collection address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should lowercase it in place", func(t *testing.T) {
			t.Parallel()
			suite := WhenValidatingAndSanitizingAgentByCollectionTokenIdRequestBeforeEach(
				"0xAbCdEf0000000000000000000000000000000001",
			)

			err := suite.sut.ValidateAndSanitize()

			assert.NoError(t, err)
			assert.Equal(t, "0xabcdef0000000000000000000000000000000001", suite.sut.AgentCollectionAddress)
		})
	})
}
