package models

const (
	// Uniswap's canonical Permit2 — same address on every chain via CREATE2.
	X402Permit2ContractAddress = "0x000000000022D473030F116dDEE9F6B43aC78BA3"
	// Permit2's domain has NO version field — adding one breaks the digest.
	X402Permit2DomainName = "Permit2"
	// The x402 proxy the signature authorizes to pull the tokens.
	X402Permit2Spender = "0x402085c248EeA27D92E8b30b2C58ed07f9E20001"
)

// X402Eip3009Authorization mirrors EIP-3009 TransferWithAuthorization; Nonce is 0x-hex 32 bytes, the rest decimal.
type X402Eip3009Authorization struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	ValidAfter  string `json:"validAfter"`
	ValidBefore string `json:"validBefore"`
	Nonce       string `json:"nonce"`
}

type X402Permit2TokenPermissions struct {
	Token  string `json:"token"`
	Amount string `json:"amount"`
}

type X402Permit2Witness struct {
	To         string `json:"to"`
	ValidAfter string `json:"validAfter"`
}

// X402Permit2Authorization mirrors the signed PermitWitnessTransferFrom; Nonce is a full uint256.
type X402Permit2Authorization struct {
	Permitted X402Permit2TokenPermissions `json:"permitted"`
	From      string                      `json:"from"`
	Spender   string                      `json:"spender"`
	Nonce     string                      `json:"nonce"`
	Deadline  string                      `json:"deadline"`
	Witness   X402Permit2Witness          `json:"witness"`
}

// X402ExactEvmPayload carries exactly one signed authorization per transfer method.
type X402ExactEvmPayload struct {
	Signature            string                    `json:"signature"`
	Authorization        *X402Eip3009Authorization `json:"authorization,omitempty"`
	Permit2Authorization *X402Permit2Authorization `json:"permit2Authorization,omitempty"`
}

// X402PaymentPayload is the PAYMENT-SIGNATURE header content before base64 encoding.
type X402PaymentPayload struct {
	X402Version int                 `json:"x402Version"`
	Scheme      string              `json:"scheme"`
	Network     string              `json:"network"`
	Payload     X402ExactEvmPayload `json:"payload"`
}
