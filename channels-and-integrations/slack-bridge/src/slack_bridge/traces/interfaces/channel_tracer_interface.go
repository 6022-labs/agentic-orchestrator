package interfaces

import "context"

// ChannelTracerInterface opens the spans of one inbound platform interaction; names stay platform-agnostic.
type ChannelTracerInterface interface {
	StartEventHandling(ctx context.Context, eventType string) (context.Context, Span)
	StartCommandHandling(ctx context.Context, command string) (context.Context, Span)
}
