package interfaces

import "context"

type BlockchainRequesterInterface interface {
	GetCurrentBlockNumber(ctx context.Context) (*uint64, error)
}
