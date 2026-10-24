package requests

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/google/uuid"
)

type ConversationByIdRequest struct {
	ConversationId uuid.UUID `swaggerignore:"true" uri:"id"`
}

func (request *ConversationByIdRequest) Validate() error {
	if request.ConversationId == uuid.Nil {
		return errors.NewValidationError("id", "must be a valid conversation id")
	}
	return nil
}
