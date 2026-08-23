package interfaces

import "context"

// BrokerTracerInterface opens the spans of one agent turn.
type BrokerTracerInterface interface {
	StartTurn(ctx context.Context, conversationId string) (context.Context, Span)
	// StartAgentResolution wraps one resolver lookup; resolver is ens, registry or card.
	StartAgentResolution(ctx context.Context, resolver string, ensDomain string) (context.Context, Span)
}
