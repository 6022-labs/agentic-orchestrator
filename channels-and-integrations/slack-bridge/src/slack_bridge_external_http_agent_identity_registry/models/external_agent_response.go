package models

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type ExternalAgentResponse struct {
	ID uuid.UUID `json:"id"`

	Role      string `json:"role"`
	Owner     string `json:"owner"`
	EnsDomain string `json:"ensDomain"`

	Images     map[string]string `json:"images"`
	Attributes map[string]string `json:"attributes"`

	CollectionTokenId uint64 `json:"collectionTokenId"`
	CollectionAddress string `json:"collectionAddress"`
}

func (e *ExternalAgentResponse) ToSlackAgent() (*models.SlackAgent, error) {
	iconImage, found := e.Images["icon"]

	var role enums.AgentRole
	if err := role.Set(e.Role); err != nil {
		return nil, fmt.Errorf("failed to set agent role: %w", err)
	}

	return &models.SlackAgent{
		Role:              role,
		OwnerWallet:       e.Owner,
		EnsDomain:         e.EnsDomain,
		CollectionTokenId: e.CollectionTokenId,
		CollectionAddress: e.CollectionAddress,
		IconImageCid:      lo.Ternary(found, &iconImage, nil),
	}, nil
}
