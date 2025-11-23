package models

type ExternalGetLLMConnectionSetupStatusResponse struct {
	IsCompleted bool    `json:"isCompleted"`
	SetupStatus *string `json:"setupStatus"`
}
