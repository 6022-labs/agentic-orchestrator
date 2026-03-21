package models

import "github.com/google/uuid"

type ExternalIdRequest struct {
	Id uuid.UUID `json:"id"`
}
