package interfaces

import "context"

type BlockEventsDequeueDaemonInterface interface {
	Start(ctx context.Context) error
	Stop() error
}
