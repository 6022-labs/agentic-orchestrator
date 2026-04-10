package models

import "github.com/google/uuid"

type ExternalCreateMCPConnectionRequest struct {
	AgentId    uuid.UUID         `json:"agentId"`
	Name       string            `json:"name"`
	Url        string            `json:"url"`
	HeaderAuth map[string]string `json:"headerAuth,omitempty"`
}
