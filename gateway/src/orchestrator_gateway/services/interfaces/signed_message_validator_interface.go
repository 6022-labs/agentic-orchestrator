package interfaces

type SignedMessageValidatorInterface interface {
	Validate(message, signature, expectedAddress string) bool
}
