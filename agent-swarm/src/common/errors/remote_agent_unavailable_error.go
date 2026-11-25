package errors

import (
	"fmt"
	"net/http"
	"time"
)

// RemoteAgentUnavailableError is returned when an agent node answers HTTP 429, 502, 503 or 504.
type RemoteAgentUnavailableError struct {
	EnsDomain  string
	StatusCode int
	RetryAfter time.Duration
	Body       string
}

func NewRemoteAgentUnavailableError(
	ensDomain string,
	statusCode int,
	retryAfter time.Duration,
	body string,
) *RemoteAgentUnavailableError {
	return &RemoteAgentUnavailableError{
		EnsDomain:  ensDomain,
		StatusCode: statusCode,
		RetryAfter: retryAfter,
		Body:       body,
	}
}

func (e *RemoteAgentUnavailableError) Error() string {
	return fmt.Sprintf("agent node %s is unavailable: status %d, response: %s", e.EnsDomain, e.StatusCode, e.Body)
}

// IsRateLimited reports a 429: the agent refused the call before working on it.
func (e *RemoteAgentUnavailableError) IsRateLimited() bool {
	return e.StatusCode == http.StatusTooManyRequests
}
