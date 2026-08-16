package mappings

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_http_agent_identity_registry/clients"
	"github.com/samber/lo"
)

// ExternalAgentResponse is the registry's generated agent wire shape, named locally to carry mapping.
type ExternalAgentResponse clients.ResponsesAgentResponse

// ToSlackAgent reads the role from attributes — the registry publishes it there, not as a top-level field.
func (response *ExternalAgentResponse) ToSlackAgent() (*models.SlackAgent, error) {
	if response == nil {
		return nil, nil
	}

	images := lo.FromPtr(response.Images)
	iconImage, found := images["icon"]

	var role enums.AgentRole
	if err := role.Set(lo.FromPtr(response.Attributes)["role"]); err != nil {
		return nil, fmt.Errorf("failed to set agent role: %w", err)
	}

	return &models.SlackAgent{
		Role:              role,
		OwnerWallet:       lo.FromPtr(response.Owner),
		EnsDomain:         lo.FromPtr(response.EnsDomain),
		CollectionTokenId: uint64(lo.FromPtr(response.CollectionTokenId)),
		CollectionAddress: lo.FromPtr(response.CollectionAddress),
		IconImageCid:      lo.Ternary(found, &iconImage, nil),
	}, nil
}
