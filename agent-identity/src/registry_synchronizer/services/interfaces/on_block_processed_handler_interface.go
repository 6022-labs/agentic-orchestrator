package interfaces

import "context"

type OnBlockProcessedHandlerInterface interface {
	OnBlockProcessed(ctx context.Context, blockNumber uint64)
}
