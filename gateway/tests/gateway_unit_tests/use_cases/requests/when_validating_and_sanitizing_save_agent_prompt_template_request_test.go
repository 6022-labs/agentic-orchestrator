package requests_test

import (
	"strings"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/stretchr/testify/assert"
)

func WhenValidatingAndSanitizingSaveAgentPromptTemplateRequestBeforeEach() *requests.SaveAgentPromptTemplateRequest {
	return &requests.SaveAgentPromptTemplateRequest{
		GetAgentPromptTemplateRequest: requests.GetAgentPromptTemplateRequest{
			AgentCollectionAddress: checksummedAddress,
			AgentCollectionTokenId: "1",
			Key:                    "clone-system",
		},
		Content: "content",
	}
}

func TestWhenValidatingAndSanitizingSaveAgentPromptTemplateRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given an address without the 0x prefix", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error on agentCollectionAddress", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingSaveAgentPromptTemplateRequestBeforeEach()
			request.AgentCollectionAddress = strings.TrimPrefix(checksummedAddress, "0x")

			err := request.ValidateAndSanitize()

			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "agentCollectionAddress", validationErr.Field)
		})
	})

	t.Run("Given an address that is not hex", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error on agentCollectionAddress", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingSaveAgentPromptTemplateRequestBeforeEach()
			request.AgentCollectionAddress = nonHexAddress

			err := request.ValidateAndSanitize()

			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "agentCollectionAddress", validationErr.Field)
		})
	})

	t.Run("Given a zero collection token ID", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error on agentCollectionTokenId", func(t *testing.T) {
			t.Parallel()

			request := WhenValidatingAndSanitizingSaveAgentPromptTemplateRequestBeforeEach()
			request.AgentCollectionTokenId = "0"

			err := request.ValidateAndSanitize()

			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "agentCollectionTokenId", validationErr.Field)
		})
	})

	t.Run("Given a checksummed address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should lowercase it in place", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingAndSanitizingSaveAgentPromptTemplateRequestBeforeEach()

			err := request.ValidateAndSanitize()

			assert.NoError(t, err)
			assert.Equal(t, lowercasedAddress, request.AgentCollectionAddress)
		})
	})
}
