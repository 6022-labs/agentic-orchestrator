package responses

import "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"

type ListAgentsInChannelResponse struct {
	Role              enums.AgentRole `json:"role"`
	EnsDomain         string          `json:"ens_domain"`
	CollectionAddress string          `json:"collection_address"`
	CollectionTokenId uint64          `json:"collection_token_id"`
}
