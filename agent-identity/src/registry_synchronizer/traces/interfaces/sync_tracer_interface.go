package interfaces

import "context"

// SyncTracerInterface opens the spans of one block's indexing pass.
type SyncTracerInterface interface {
	StartBlockProcessing(ctx context.Context, blockNumber uint64) (context.Context, Span)
	// StartEventHandling opens a child span per contract event batch; the block number sits on the parent.
	StartEventHandling(ctx context.Context, eventName string) (context.Context, Span)
}
