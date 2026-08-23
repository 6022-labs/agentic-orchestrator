package traces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/traces/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "agentic_orchestrator_conversation_broker"

type BrokerTracer struct {
	tracer trace.Tracer
}

func NewBrokerTracer() *BrokerTracer {
	return &BrokerTracer{tracer: otel.Tracer(tracerName)}
}

func (t *BrokerTracer) StartTurn(ctx context.Context, conversationId string) (context.Context, interfaces.Span) {
	// The daemon picks a turn up long after the message arrived, so its trace is its own; the link keeps the cause.
	ctx, span := t.tracer.Start(ctx, "broker.turn.process",
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithNewRoot(),
		trace.WithLinks(trace.LinkFromContext(ctx)),
	)
	span.SetAttributes(attribute.String("gen_ai.conversation.id", conversationId))
	return ctx, &otelSpan{span: span}
}

func (t *BrokerTracer) StartAgentResolution(
	ctx context.Context,
	resolver string,
	ensDomain string,
) (context.Context, interfaces.Span) {
	ctx, span := t.tracer.Start(ctx, "broker.agent.resolve "+resolver)
	span.SetAttributes(
		attribute.String("resolver", resolver),
		attribute.String("gen_ai.agent.name", ensDomain),
	)
	return ctx, &otelSpan{span: span}
}
