package interfaces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
)

type PlatformBridgeRequesterInterface interface {
	SetWalletOwnershipToUser(ctx context.Context, platform enums.BridgePlatform, challenge, walletAddress string) error
}
