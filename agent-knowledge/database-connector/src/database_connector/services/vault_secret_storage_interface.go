package services

type VaultSecretStorageInterface interface {
	Store(key string, value string) error
	Read(key string) (string, error)
	Delete(key string) error
}
