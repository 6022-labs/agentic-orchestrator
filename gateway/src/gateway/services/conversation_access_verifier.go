package services

import (
	"context"
	"errors"

	gateway_errors "github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/google/uuid"
)

const ConversationAccessForbiddenCode = "conversation_access_forbidden"

// ConversationAccessVerifier admits a caller who controls at least one of the conversation's participating agents.
type ConversationAccessVerifier struct {
	agentAccessVerifier         interfaces.AgentAccessVerifierInterface
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface
}

func NewConversationAccessVerifier(
	agentAccessVerifier interfaces.AgentAccessVerifierInterface,
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface,
) *ConversationAccessVerifier {
	return &ConversationAccessVerifier{
		agentAccessVerifier:         agentAccessVerifier,
		conversationBrokerRequester: conversationBrokerRequester,
	}
}

func (v *ConversationAccessVerifier) EnsureAccess(
	ctx context.Context,
	callerId string,
	conversationId uuid.UUID,
) error {
	participants, err := v.conversationBrokerRequester.ListConversationParticipants(ctx, conversationId)
	if err != nil {
		if v.isTypedError(err) {
			return err
		}
		return gateway_errors.NewUnavailableError(
			"conversation_broker_unavailable",
			"conversation broker is unavailable",
			err,
		)
	}

	for _, participant := range participants {
		if participant == nil || participant.CollectionAddress == nil || participant.CollectionTokenId == nil {
			continue
		}

		err := v.agentAccessVerifier.EnsureAccess(ctx,
			callerId,
			*participant.CollectionAddress,
			*participant.CollectionTokenId,
		)
		if err == nil {
			return nil
		}
		if !v.isAccessDeniedError(err) {
			return err
		}
	}

	return gateway_errors.NewForbiddenError(
		ConversationAccessForbiddenCode,
		"caller has no access to any of this conversation's participating agents",
	)
}

// isTypedError keeps the broker's domain errors (e.g. NotFound) instead of masking them as unavailability.
func (v *ConversationAccessVerifier) isTypedError(err error) bool {
	var validationErr *gateway_errors.ValidationError
	var unauthorizedErr *gateway_errors.UnauthorizedError
	var forbiddenErr *gateway_errors.ForbiddenError
	var notFoundErr *gateway_errors.NotFoundError
	var unavailableErr *gateway_errors.UnavailableError
	var internalErr *gateway_errors.InternalError
	return errors.As(err, &validationErr) || errors.As(err, &unauthorizedErr) || errors.As(err, &forbiddenErr) ||
		errors.As(err, &notFoundErr) || errors.As(err, &unavailableErr) || errors.As(err, &internalErr)
}

// isAccessDeniedError tells "this participant is not yours" apart from infrastructure failures.
func (v *ConversationAccessVerifier) isAccessDeniedError(err error) bool {
	var forbiddenErr *gateway_errors.ForbiddenError
	var notFoundErr *gateway_errors.NotFoundError
	return errors.As(err, &forbiddenErr) || errors.As(err, &notFoundErr)
}
