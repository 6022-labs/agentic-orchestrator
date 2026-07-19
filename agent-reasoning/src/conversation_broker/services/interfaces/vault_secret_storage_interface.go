package interfaces

// VaultSecretStorageInterface is the credential-store port backed by Vault KV v1.
type VaultSecretStorageInterface interface {
	Store(key string, value string) error
	Read(key string) (string, error)
	Delete(key string) error
}
