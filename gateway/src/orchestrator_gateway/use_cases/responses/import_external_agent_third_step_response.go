package responses

type ImportExternalAgentThirdStepResponse struct {
	Signer    string `json:"signer"`
	Signature string `json:"signature"`
}
