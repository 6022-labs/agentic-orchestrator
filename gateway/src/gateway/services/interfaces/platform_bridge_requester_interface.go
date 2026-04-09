package interfaces

import "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"

type PlatformBridgeRequesterInterface interface {
	SetWalletOwnershipToUser(platform enums.AgentContextMessagePlatform, challenge, walletAddress string) error
}
