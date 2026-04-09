package interfaces

type OrchestratorPrivateKeyFetcherInterface interface {
	FetchPrivateKey() (*string, error)
}
