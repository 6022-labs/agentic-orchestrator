package models

// ExternalRequestHeaders carries optional protocol headers of a node call; nil means none, empty fields are skipped.
type ExternalRequestHeaders struct {
	PaymentSignature string `header:"PAYMENT-SIGNATURE"`
}
