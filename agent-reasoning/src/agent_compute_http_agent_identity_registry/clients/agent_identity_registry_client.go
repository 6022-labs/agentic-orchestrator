package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_agent_identity_registry/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_agent_identity_registry/settings"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const (
	GET_AGENT_ENDPOINT               = "/agent/%s/%d"
	GET_AGENT_BY_ENS_DOMAIN_ENDPOINT = "/agent/ens/%s"
)

type AgentIdentityRegistryClientInterface interface {
	GetAgent(agentCollectionAddress string, agentCollectionTokenId uint64) (*models.ExternalAgentResponse, error)
	GetMintedAgentByEnsDomain(ensDomain string) (*models.ExternalAgentResponse, error)
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

func (client *AgentIdentityRegistryClient) GetAgent(agentCollectionAddress string, agentCollectionTokenId uint64) (*models.ExternalAgentResponse, error) {
	url := fmt.Sprintf(client.agentIdentityRegistrySettings.BaseUrl+GET_AGENT_ENDPOINT, agentCollectionAddress, agentCollectionTokenId)

	httpRequest, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode == http.StatusOK {
		var response models.ExternalAgentResponse
		if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return &response, nil
	}

	if httpResponse.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	responseBody, _ := io.ReadAll(httpResponse.Body)
	return nil, fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
}

func (client *AgentIdentityRegistryClient) GetMintedAgentByEnsDomain(ensDomain string) (*models.ExternalAgentResponse, error) {
	url := fmt.Sprintf(client.agentIdentityRegistrySettings.BaseUrl+GET_AGENT_BY_ENS_DOMAIN_ENDPOINT, ensDomain)

	httpRequest, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if httpResponse.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	var response models.ExternalAgentResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
