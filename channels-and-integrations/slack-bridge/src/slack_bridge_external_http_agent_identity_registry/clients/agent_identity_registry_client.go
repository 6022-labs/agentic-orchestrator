package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_http_agent_identity_registry/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_http_agent_identity_registry/settings"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const (
	GET_AGENT_ENDPOINT = "/agent/%s/%d"
)

type AgentIdentityRegistryClientInterface interface {
	GetAgent(request *models.ExternalAgentByCollectionTokenIdRequest) (*models.ExternalAgentResponse, error)
}

type AgentIdentityRegistryClient struct {
	logger                        *zap.Logger
	httpClient                    *http.Client
	agentIdentityRegistrySettings *settings.AgentIdentityRegistrySettings
}

type newAgentIdentityRegistryClientParams struct {
	dig.In

	Logger                        *zap.Logger
	HttpClient                    *http.Client `name:"agent-identity-registry-http-client"`
	AgentIdentityRegistrySettings *settings.AgentIdentityRegistrySettings
}

func NewAgentIdentityRegistryClient(params newAgentIdentityRegistryClientParams) *AgentIdentityRegistryClient {
	return &AgentIdentityRegistryClient{
		logger:                        params.Logger,
		httpClient:                    params.HttpClient,
		agentIdentityRegistrySettings: params.AgentIdentityRegistrySettings,
	}
}

func (client *AgentIdentityRegistryClient) GetAgent(request *models.ExternalAgentByCollectionTokenIdRequest) (*models.ExternalAgentResponse, error) {
	url := fmt.Sprintf(client.agentIdentityRegistrySettings.BaseUrl+GET_AGENT_ENDPOINT, request.AgentCollectionAddress, request.AgentCollectionTokenId)

	httpRequest, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	client.logger.Debug(
		"Sending GetAgent request to agent identity registry",
		zap.String("url", url),
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
			"Failed to get agent from agent identity registry",
			zap.String("url", url),
			zap.Int("status_code", httpResponse.StatusCode),
			zap.String("body", string(responseBody)),
		)
		return nil, fmt.Errorf("unexpected status code: %d", httpResponse.StatusCode)
	}

	var agent models.ExternalAgentResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&agent)
	if err != nil {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		client.logger.Warn(
			"Failed to decode GetAgent response body",
			zap.String("url", url),
			zap.Int("status_code", httpResponse.StatusCode),
			zap.String("body", string(responseBody)),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &agent, nil
}
