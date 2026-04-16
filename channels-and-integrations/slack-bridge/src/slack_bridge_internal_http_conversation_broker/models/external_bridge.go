package models

import "github.com/google/uuid"

type ExternalRegisterBridgeRequest struct {
	ID          *uuid.UUID `json:"id,omitempty"`
	Kind        string     `json:"kind"`
	CallbackUrl string     `json:"callbackUrl"`
}

type ExternalBridgeResponse struct {
	ID          uuid.UUID `json:"id"`
	Kind        string    `json:"kind"`
	CallbackUrl string    `json:"callbackUrl"`
}
