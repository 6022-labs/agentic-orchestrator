package models

import "github.com/google/uuid"

type ExternalCreateAgentSecondStepRequest struct {
	AuthenticatedWallet string    `json:"authenticatedWallet"`
	AgentId             uuid.UUID `json:"agentId"`
}
