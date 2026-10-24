package requests_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func WhenValidatingConversationByIdRequestBeforeEach() *requests.ConversationByIdRequest {
	return &requests.ConversationByIdRequest{
		ConversationId: uuid.New(),
	}
}

func TestWhenValidatingConversationByIdRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given an empty conversation id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error on id", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingConversationByIdRequestBeforeEach()
			request.ConversationId = uuid.Nil

			err := request.Validate()

			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "id", validationErr.Field)
		})
	})

	t.Run("Given a valid request", func(t *testing.T) {
		t.Parallel()

		t.Run("Should accept it", func(t *testing.T) {
			t.Parallel()
			request := WhenValidatingConversationByIdRequestBeforeEach()

			assert.NoError(t, request.Validate())
		})
	})
}
