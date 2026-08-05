package models

// ExternalResponse models only the response fields the broker reads.
type ExternalResponse struct {
	Id     string                    `json:"id"`
	Status string                    `json:"status"`
	Output []*ExternalResponseOutput `json:"output"`
	Error  *ExternalResponseError    `json:"error"`
}

type ExternalResponseOutput struct {
	Type    string                           `json:"type"`
	Role    string                           `json:"role"`
	Content []*ExternalResponseOutputContent `json:"content"`
}

type ExternalResponseOutputContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ExternalResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
