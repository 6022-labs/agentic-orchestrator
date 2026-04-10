package interfaces

type VaultSecretStorageInterface interface {
	StoreKey(path, key, value string) error
	Store(path string, value map[string]string) error
	ReadKey(path, key string) (string, error)
	Read(path string) (map[string]string, error)
	DeleteKey(path, key string) error
	Delete(path string) error
}
