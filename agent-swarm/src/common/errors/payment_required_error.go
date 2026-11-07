package errors

import "fmt"

// PaymentRequiredError is returned when an agent node answers HTTP 402.
// PaymentRequirements carries the raw PAYMENT-REQUIRED header payload so a
// future x402 payer can decode it, sign a payment and retry the turn.
type PaymentRequiredError struct {
	EnsDomain           string `json:"ensDomain"`
	PaymentRequirements string `json:"paymentRequirements,omitempty"`
}

func NewPaymentRequiredError(ensDomain, paymentRequirements string) *PaymentRequiredError {
	return &PaymentRequiredError{EnsDomain: ensDomain, PaymentRequirements: paymentRequirements}
}

func (e *PaymentRequiredError) Error() string {
	return fmt.Sprintf("agent node %s requires payment", e.EnsDomain)
}
