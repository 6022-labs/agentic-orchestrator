package models

import "github.com/google/uuid"

type ExternalRegisterBridgeRequest struct {
	ID      *uuid.UUID `json:"id,omitempty"`
	Kind    string     `json:"kind"`
	BaseUrl string     `json:"baseUrl"`
}

type ExternalBridgeResponse struct {
	ID      uuid.UUID `json:"id"`
	Kind    string    `json:"kind"`
	BaseUrl string    `json:"baseUrl"`
}
