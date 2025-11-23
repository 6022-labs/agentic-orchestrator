package services

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_platform_bridge/clients"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_platform_bridge/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_platform_bridge/settings"
)

type PlatformBridgeRequester struct {
	platformBridgeSettings *settings.PlatformBridgesSettings
	platformBridgeClient   clients.PlatformBridgeClientInterface
}

func NewPlatformBridgeRequester(
	platformBridgeSettings *settings.PlatformBridgesSettings,
	platformBridgeClient clients.PlatformBridgeClientInterface,
) *PlatformBridgeRequester {
	return &PlatformBridgeRequester{
		platformBridgeClient:   platformBridgeClient,
		platformBridgeSettings: platformBridgeSettings,
	}
}

func (r *PlatformBridgeRequester) SetWalletOwnershipToUser(platform enums.AgentContextMessagePlatform, challenge, walletAddress string) error {
	platformBaseUrl, err := r.getPlatformBaseUrl(platform)
	if err != nil {
		return err
	}

	return r.platformBridgeClient.SetWalletOwnershipToUser(platformBaseUrl, models.ExternalSetWalletOwnershipToUserRequest{
		Challenge:     challenge,
		WalletAddress: walletAddress,
	})
}

func (r *PlatformBridgeRequester) getPlatformBaseUrl(platform enums.AgentContextMessagePlatform) (string, error) {
	switch platform {
	case enums.AgentContextMessagePlatformSlack:
		return r.platformBridgeSettings.SlackBridgeUrl, nil
	default:
		return "", fmt.Errorf("unsupported platform: %s", platform)
	}
}
