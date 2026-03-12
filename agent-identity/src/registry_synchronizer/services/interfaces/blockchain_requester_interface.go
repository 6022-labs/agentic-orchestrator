package interfaces

type BlockchainRequesterInterface interface {
	GetCurrentBlockNumber() (*uint64, error)
}
