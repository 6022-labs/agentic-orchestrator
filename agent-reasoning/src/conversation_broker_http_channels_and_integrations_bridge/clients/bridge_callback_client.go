package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_channels_and_integrations_bridge/models"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

// BridgeCallbackEndpoint is the path the broker appends to a bridge's
// base URL when posting a completion callback.
const BridgeCallbackEndpoint = "/broker/callback"

type BridgeCallbackClientInterface interface {
	PostCompletion(bridgeBaseUrl string, body *models.ExternalCompletionCallbackRequest) error
}

type BridgeCallbackClient struct {
	logger     *zap.Logger
	httpClient *http.Client
}

type NewBridgeCallbackClientParams struct {
	dig.In

	Logger     *zap.Logger
	HttpClient *http.Client `name:"bridge-callback-http-client"`
}

func NewBridgeCallbackClient(params NewBridgeCallbackClientParams) *BridgeCallbackClient {
	return &BridgeCallbackClient{
		logger:     params.Logger,
		httpClient: params.HttpClient,
	}
}

func (client *BridgeCallbackClient) PostCompletion(bridgeBaseUrl string, body *models.ExternalCompletionCallbackRequest) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal callback body: %w", err)
	}

	url := strings.TrimRight(bridgeBaseUrl, "/") + BridgeCallbackEndpoint

	httpRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	client.logger.Debug(
		"Sending completion callback to bridge",
		zap.String("url", url),
		zap.String("body", string(bodyBytes)),
	)

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		client.logger.Warn(
			"Bridge callback returned non-2xx",
			zap.String("url", url),
			zap.Int("status_code", httpResponse.StatusCode),
			zap.String("body", string(responseBody)),
		)
		return fmt.Errorf("unexpected status code: %d", httpResponse.StatusCode)
	}

	return nil
}
