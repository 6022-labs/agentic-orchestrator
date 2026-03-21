package responses

type CreateAgentSecondStepResponse struct {
	Nonce     string `json:"nonce"`
	Signer    string `json:"signer"`
	Signature string `json:"signature"`
}
