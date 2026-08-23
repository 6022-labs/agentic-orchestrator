package interfaces

import (
	"context"
	"time"
)

const (
	PaymentOutcomeSettled  = "settled"
	PaymentOutcomeRejected = "rejected"
	PaymentOutcomeSkipped  = "skipped"
)

// RemoteAgentMetricsInterface records agent invocations under the OTel GenAI conventions, plus x402 settlements.
type RemoteAgentMetricsInterface interface {
	// RecordInvocation emits gen_ai.client.operation.duration for one invoke_agent call; endpointUrl is the
	// full remote endpoint, split into semconv server.address/server.port.
	RecordInvocation(ctx context.Context, agentName string, endpointUrl string, duration time.Duration, err error)
	RecordPaymentSettlement(ctx context.Context, outcome string, network string)
}
