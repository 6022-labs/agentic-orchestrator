package traces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/traces/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// tracerName is shared by every channels-and-integrations bridge; channel.platform separates their spans.
const tracerName = "agentic_orchestrator_channels_and_integrations"

type ChannelTracer struct {
	tracer   trace.Tracer
	platform attribute.KeyValue
}

func NewChannelTracer(platform string) *ChannelTracer {
	return &ChannelTracer{
		tracer:   otel.Tracer(tracerName),
		platform: attribute.String("channel.platform", platform),
	}
}

func (t *ChannelTracer) StartEventHandling(ctx context.Context, eventType string) (context.Context, interfaces.Span) {
	ctx, span := t.tracer.Start(ctx, "channel.event.handle "+eventType)
	span.SetAttributes(t.platform, attribute.String("channel.event.type", eventType))
	return ctx, &otelSpan{span: span}
}

func (t *ChannelTracer) StartCommandHandling(ctx context.Context, command string) (context.Context, interfaces.Span) {
	ctx, span := t.tracer.Start(ctx, "channel.command.handle "+command)
	span.SetAttributes(t.platform, attribute.String("channel.command.name", command))
	return ctx, &otelSpan{span: span}
}
