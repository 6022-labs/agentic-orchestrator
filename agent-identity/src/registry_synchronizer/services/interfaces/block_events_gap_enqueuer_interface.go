package interfaces

import "context"

type BlockEventsGapEnqueuerInterface interface {
	EnqueueGap(ctx context.Context) error
}
