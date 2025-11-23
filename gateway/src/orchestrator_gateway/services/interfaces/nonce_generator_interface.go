package interfaces

type NonceGeneratorInterface interface {
	GenerateNonce() (*string, error)
}
