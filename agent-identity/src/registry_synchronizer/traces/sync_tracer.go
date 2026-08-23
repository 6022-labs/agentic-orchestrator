package traces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/traces/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "agentic_orchestrator_registry_synchronizer"

type SyncTracer struct {
	tracer trace.Tracer
}

func NewSyncTracer() *SyncTracer {
	return &SyncTracer{tracer: otel.Tracer(tracerName)}
}

func (t *SyncTracer) StartBlockProcessing(
	ctx context.Context,
	blockNumber uint64,
) (context.Context, interfaces.Span) {
	// A block is dequeued long after it was received, so its trace is its own; the link keeps the cause.
	ctx, span := t.tracer.Start(ctx, "sync.block.process",
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithNewRoot(),
		trace.WithLinks(trace.LinkFromContext(ctx)),
	)
	span.SetAttributes(attribute.Int64("sync.block.number", int64(blockNumber)))
	return ctx, &otelSpan{span: span}
}

func (t *SyncTracer) StartEventHandling(ctx context.Context, eventName string) (context.Context, interfaces.Span) {
	ctx, span := t.tracer.Start(ctx, "sync.event.handle "+eventName)
	span.SetAttributes(attribute.String("chain.event.name", eventName))
	return ctx, &otelSpan{span: span}
}
