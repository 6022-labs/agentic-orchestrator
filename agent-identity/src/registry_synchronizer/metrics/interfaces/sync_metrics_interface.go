package interfaces

import (
	"context"
	"time"
)

const (
	BlockOutcomeProcessed = "processed"
	BlockOutcomeFailed    = "failed"
)

// SyncMetricsInterface records the sync.* metrics of the block indexing pipeline.
type SyncMetricsInterface interface {
	// RecordQueueState publishes the depths after a mutation; a snapshot self-corrects, a delta cannot.
	RecordQueueState(ctx context.Context, readyBlocks int64, pendingBlocks int64)
	// RecordHighestBlockWithEvents is not the chain head: it stops advancing on an idle chain.
	RecordHighestBlockWithEvents(ctx context.Context, blockNumber uint64)
	RecordBlockProcessed(ctx context.Context, outcome string, duration time.Duration)
	// RecordCheckpoint publishes the last block durably indexed.
	RecordCheckpoint(ctx context.Context, blockNumber uint64)
	RecordBackfilledBlocks(ctx context.Context, blockCount int64)
}
