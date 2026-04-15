package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_agent_reasoning_agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_agent_reasoning_agent_compute/settings"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const (
	AGENT_COMPUTE_ENDPOINT = "/agent/%s/%d/compute"
)

type AgentReasoningAgentComputeClientInterface interface {
	Completions(collectionAddress string, collectionTokenId uint64, body *models.ExternalAgentComputeRequest) (*models.ExternalAgentComputeResponses, error)
}

type AgentReasoningAgentComputeClient struct {
	logger                             *zap.Logger
	httpClient                         *http.Client
	agentReasoningAgentComputeSettings *settings.AgentReasoningAgentComputeSettings
}

type NewAgentReasoningAgentComputeClientParams struct {
	dig.In

	Logger                             *zap.Logger
	HttpClient                         *http.Client `name:"agent-reasoning-agent-compute-http-client"`
	AgentReasoningAgentComputeSettings *settings.AgentReasoningAgentComputeSettings
}

func NewAgentReasoningAgentComputeClient(params NewAgentReasoningAgentComputeClientParams) *AgentReasoningAgentComputeClient {
	return &AgentReasoningAgentComputeClient{
		logger:                             params.Logger,
		httpClient:                         params.HttpClient,
		agentReasoningAgentComputeSettings: params.AgentReasoningAgentComputeSettings,
	}
}

func (client *AgentReasoningAgentComputeClient) Completions(
	collectionAddress string,
	collectionTokenId uint64,
	body *models.ExternalAgentComputeRequest,
) (*models.ExternalAgentComputeResponses, error) {
	endpoint := fmt.Sprintf(AGENT_COMPUTE_ENDPOINT, collectionAddress, collectionTokenId)
	url := client.agentReasoningAgentComputeSettings.BaseUrl + endpoint

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
	if err := json.NewDecoder(httpResponse.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}
	return &result, nil
}
