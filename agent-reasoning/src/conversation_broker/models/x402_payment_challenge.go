package models

const (
	X402Version                    = 2
	X402SchemeExact                = "exact"
	X402AssetTransferMethodEip3009 = "eip3009"
	X402AssetTransferMethodPermit2 = "permit2"
)

// X402PaymentOption is one accepted way to pay; Extra carries the asset's EIP-712 domain for eip3009 signing.
type X402PaymentOption struct {
	Scheme            string         `json:"scheme"`
	Network           string         `json:"network"`
	Amount            string         `json:"amount"`
	Asset             string         `json:"asset"`
	PayTo             string         `json:"payTo"`
	MaxTimeoutSeconds int            `json:"maxTimeoutSeconds"`
	Extra             map[string]any `json:"extra,omitempty"`
}

// X402PaymentChallenge is the decoded PAYMENT-REQUIRED header of an x402 v2 402 response.
type X402PaymentChallenge struct {
	X402Version int                 `json:"x402Version"`
	Accepts     []X402PaymentOption `json:"accepts"`
	Error       string              `json:"error,omitempty"`
}
