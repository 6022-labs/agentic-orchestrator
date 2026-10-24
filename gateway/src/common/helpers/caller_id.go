package helpers

import (
	"context"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
)

// CallerIdMaxLength is kept identical to agent-node's caller id, in characters.
const CallerIdMaxLength = 128

var callerIdRule = "must be a non-empty caller id of at most " + strconv.Itoa(CallerIdMaxLength) + " characters"

type callerIdContextKey struct{}

// ContextWithCallerId attaches the request's caller id, normalized.
func ContextWithCallerId(ctx context.Context, callerId string) context.Context {
	return context.WithValue(ctx, callerIdContextKey{}, NormalizeCallerId(callerId))
}

// CallerIdFromContext reports false when no caller id was attached.
func CallerIdFromContext(ctx context.Context) (string, bool) {
	callerId, ok := ctx.Value(callerIdContextKey{}).(string)
	return callerId, ok && callerId != ""
}

// RequireCallerId fails Unauthorized without a caller id and Validation when it exceeds CallerIdMaxLength.
func RequireCallerId(ctx context.Context) (string, error) {
	callerId, ok := CallerIdFromContext(ctx)
	if !ok {
		return "", errors.NewUnauthorizedError("missing_caller_id", "an authenticated caller id is required")
	}
	if !IsCallerId(callerId) {
		return "", errors.NewValidationError("callerId", callerIdRule)
	}
	return callerId, nil
}

// NormalizeCallerId lowercases: callers are lowercase wallet addresses, compared case-insensitively.
func NormalizeCallerId(callerId string) string {
	return strings.ToLower(strings.TrimSpace(callerId))
}

// IsCallerId accepts any non-blank opaque id of at most CallerIdMaxLength characters.
func IsCallerId(callerId string) bool {
	return strings.TrimSpace(callerId) != "" && utf8.RuneCountInString(callerId) <= CallerIdMaxLength
}
