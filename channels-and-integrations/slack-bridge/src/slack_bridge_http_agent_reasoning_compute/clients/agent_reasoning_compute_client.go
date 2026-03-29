package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_reasoning_compute/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_reasoning_compute/settings"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const (
	AGENT_COMPUTE_ENDPOINT = "/agent/%s/%d/compute"
)

type AgentReasoningComputeClientInterface interface {
	Completions(collectionAddress string, collectionTokenId uint64, body *models.ExternalAgentComputeRequest) (*models.ExternalAgentComputeResponses, error)
}

type AgentReasoningComputeClient struct {
	logger                        *zap.Logger
	httpClient                    *http.Client
	agentReasoningComputeSettings *settings.AgentReasoningComputeSettings
}

type newAgentReasoningComputeClientParams struct {
	dig.In

	Logger                        *zap.Logger
	HttpClient                    *http.Client `name:"agent-reasoning-compute-http-client"`
	AgentReasoningComputeSettings *settings.AgentReasoningComputeSettings
}

func NewAgentReasoningComputeClient(params newAgentReasoningComputeClientParams) *AgentReasoningComputeClient {
	return &AgentReasoningComputeClient{
		logger:                        params.Logger,
		httpClient:                    params.HttpClient,
		agentReasoningComputeSettings: params.AgentReasoningComputeSettings,
	}
}

func (client *AgentReasoningComputeClient) Completions(collectionAddress string, collectionTokenId uint64, body *models.ExternalAgentComputeRequest) (*models.ExternalAgentComputeResponses, error) {
	endpoint := fmt.Sprintf(AGENT_COMPUTE_ENDPOINT, collectionAddress, collectionTokenId)
	url := fmt.Sprintf("%s%s", client.agentReasoningComputeSettings.BaseUrl, endpoint)

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	client.logger.Debug(
		"Sending completions request to agent-reasoning compute",
		zap.String("url", url),
		zap.String("body", string(bodyBytes)),
		zap.Any("header", httpRequest.Header),
	)

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		client.logger.Warn(
			"Failed to get completions from agent-reasoning compute",
			zap.String("url", url),
			zap.Int("status_code", httpResponse.StatusCode),
			zap.String("body", string(responseBody)),
		)
		return nil, fmt.Errorf("unexpected status code: %d", httpResponse.StatusCode)
	}

	var result models.ExternalAgentComputeResponses
	err = json.NewDecoder(httpResponse.Body).Decode(&result)
	if err != nil {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		client.logger.Warn(
			"Failed to decode completions response body",
			zap.String("url", url),
			zap.Int("status_code", httpResponse.StatusCode),
			zap.String("body", string(responseBody)),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &result, nil
}
