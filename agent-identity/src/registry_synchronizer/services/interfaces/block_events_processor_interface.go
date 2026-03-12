package interfaces

import "context"

type BlockEventsProcessorInterface interface {
	ProcessBlock(ctx context.Context, blockNumber uint64, item *QueueItem) error
}
