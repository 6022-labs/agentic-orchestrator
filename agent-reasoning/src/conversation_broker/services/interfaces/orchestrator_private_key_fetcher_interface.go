package interfaces

// OrchestratorPrivateKeyFetcherInterface reads the orchestrator wallet key from the credential store.
type OrchestratorPrivateKeyFetcherInterface interface {
	FetchPrivateKey() (*string, error)
}
