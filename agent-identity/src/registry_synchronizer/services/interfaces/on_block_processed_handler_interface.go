package interfaces

type OnBlockProcessedHandlerInterface interface {
	OnBlockProcessed(blockNumber uint64)
}
