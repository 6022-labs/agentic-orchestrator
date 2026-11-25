package interfaces

import "context"

type BlockEventsDequeueDaemonInterface interface {
	// Start processes queued block events in the background; only the first call has an effect.
	Start(ctx context.Context)
	// Stop cancels the block in progress and waits for it to return.
	Stop()
	// Failed yields the error that ended the daemon; a Stop yields nothing.
	Failed() <-chan error
}
