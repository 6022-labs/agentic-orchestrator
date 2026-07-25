package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_channels_and_integrations_bridge/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_channels_and_integrations_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_channels_and_integrations_bridge/settings"
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

func (r *PlatformBridgeRequester) SetWalletOwnershipToUser(platform enums.BridgePlatform, challenge, walletAddress string) error {
	platformBaseUrl, err := r.getPlatformBaseUrl(platform)
	if err != nil {
		return err
	}

	return r.platformBridgeClient.SetWalletOwnershipToUser(platformBaseUrl, models.ExternalSetWalletOwnershipToUserRequest{
		Challenge:     challenge,
		WalletAddress: walletAddress,
	})
}

func (r *PlatformBridgeRequester) getPlatformBaseUrl(platform enums.BridgePlatform) (string, error) {
	switch platform {
	case enums.BridgePlatformSlack:
		return r.platformBridgeSettings.SlackBridgeUrl, nil
	default:
		return "", fmt.Errorf("unsupported platform: %s", platform)
	}
}
