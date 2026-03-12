package interfaces

type BlockEventsDequeueDaemonInterface interface {
	Start() error
	Stop() error
}
