package interfaces

type BlockEventHandlerInterface interface {
	HandleBlockProcessed(blockNumber uint64) error
}
