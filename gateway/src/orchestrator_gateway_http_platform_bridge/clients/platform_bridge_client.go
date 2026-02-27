package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway_http_platform_bridge/models"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const (
	SET_WALLET_OWNERSHIP_ENDPOINT = "/wallet/set-ownership-to-slack-user"
)

type PlatformBridgeClientInterface interface {
	SetWalletOwnershipToUser(platformBaseUrl string, request models.ExternalSetWalletOwnershipToUserRequest) error
}

type PlatformBridgeClient struct {
	logger     *zap.Logger
	httpClient *http.Client
}

type newPlatformBridgeClientParams struct {
	dig.In

	Logger     *zap.Logger
	HttpClient *http.Client `name:"platform-bridge-http-client"`
}

func NewPlatformBridgeClient(
	params newPlatformBridgeClientParams,
) *PlatformBridgeClient {
	return &PlatformBridgeClient{
		logger:     params.Logger,
		httpClient: params.HttpClient,
	}
}

func (client *PlatformBridgeClient) SetWalletOwnershipToUser(platformBaseUrl string, request models.ExternalSetWalletOwnershipToUserRequest) error {
	url := platformBaseUrl + SET_WALLET_OWNERSHIP_ENDPOINT
	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	client.logger.Debug("Sending Set Wallet Ownership to User request",
		zap.String("method", http.MethodPost),
		zap.String("url", url),
		zap.Any("body", request),
	)

	httpRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status: %s", httpResponse.Status)
	}

	return nil
}
