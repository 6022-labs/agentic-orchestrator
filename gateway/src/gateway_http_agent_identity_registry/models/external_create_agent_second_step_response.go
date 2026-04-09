package models

type ExternalCreateAgentSecondStepResponse struct {
	Nonce     string `json:"nonce"`
	Signer    string `json:"signer"`
	Signature string `json:"signature"`
}
