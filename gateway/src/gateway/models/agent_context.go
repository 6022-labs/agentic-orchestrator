package models

import "github.com/google/uuid"

type AgentContext struct {
	Id  uuid.UUID
	Key string
}
