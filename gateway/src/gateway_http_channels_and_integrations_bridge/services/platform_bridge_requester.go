package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_channels_and_integrations_bridge/clients"
)

const jsonContentType = "application/json"

type PlatformBridgeRequester struct {
	platformClients map[enums.BridgePlatform]clients.ClientWithResponsesInterface
}

func NewPlatformBridgeRequester(
	platformClients map[enums.BridgePlatform]clients.ClientWithResponsesInterface,
) *PlatformBridgeRequester {
	return &PlatformBridgeRequester{
		platformClients: platformClients,
	}
}

func (r *PlatformBridgeRequester) SetWalletOwnershipToUser(platform enums.BridgePlatform, challenge, walletAddress string) error {
	client, found := r.platformClients[platform]
	if !found {
		return fmt.Errorf("unsupported platform: %s", platform)
	}

	body, err := json.Marshal(clients.RequestsSetWalletOwnershipToSlackUserRequest{
		Challenge:     &challenge,
		WalletAddress: &walletAddress,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	rsp, err := client.PostWalletSetOwnershipToSlackUserWithBodyWithResponse(context.Background(), jsonContentType, bytes.NewReader(body))
	if err != nil {
		return err
	}
	if rsp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, response: %s", rsp.StatusCode(), string(rsp.Body))
	}

	return nil
}
