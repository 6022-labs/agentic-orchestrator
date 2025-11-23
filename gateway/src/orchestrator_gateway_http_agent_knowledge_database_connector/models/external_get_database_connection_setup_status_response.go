package models

type ExternalGetDatabaseConnectionSetupStatusResponse struct {
	IsCompleted bool    `json:"isCompleted"`
	SetupStatus *string `json:"setupStatus"`
}
