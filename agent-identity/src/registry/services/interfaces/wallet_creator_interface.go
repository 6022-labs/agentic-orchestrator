package interfaces

type WalletCreatorInterface interface {
	Create() (*string, *string, error)
}
