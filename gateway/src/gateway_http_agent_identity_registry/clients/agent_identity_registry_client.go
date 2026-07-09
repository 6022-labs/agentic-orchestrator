package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	orchestrator_errors "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/settings"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const (
	LIST_AGENTS_ENDPOINT            = "/agents"
	GET_AGENT_ENDPOINT              = "/agent/%s/%d"
	LIST_AGENT_COLLECTIONS_ENDPOINT = "/agent-collections"
)

type AgentIdentityRegistryClientInterface interface {
	ListAgents(request *models.ExternalListAgentsRequest) (*models.ExternalListAgentsResponse, error)
	GetAgent(agentCollectionAddress string, agentCollectionTokenId uint64) (*models.ExternalAgentResponse, error)
	ListAgentCollections() ([]*models.ExternalAgentCollectionResponse, error)
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

func (client *AgentIdentityRegistryClient) ListAgents(request *models.ExternalListAgentsRequest) (*models.ExternalListAgentsResponse, error) {
	base, err := url.Parse(client.agentIdentityRegistrySettings.BaseUrl + LIST_AGENTS_ENDPOINT)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := base.Query()
	if request.Page != nil {
		q.Set("page", strconv.FormatUint(uint64(*request.Page), 10))
	}
	if request.ElementsPerPage != nil {
		q.Set("elementsPerPage", strconv.FormatUint(uint64(*request.ElementsPerPage), 10))
	}
	if request.SortBy != nil {
		q.Set("sortBy", *request.SortBy)
	}
	if request.SortDirection != nil {
		q.Set("sortDirection", *request.SortDirection)
	}
	base.RawQuery = q.Encode()

	httpRequest, err := http.NewRequest(http.MethodGet, base.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		if httpResponse.StatusCode == http.StatusBadRequest {
			var validationErr orchestrator_errors.ValidationError
			if err := json.NewDecoder(httpResponse.Body).Decode(&validationErr); err == nil {
				return nil, &validationErr
			}
		}

		responseBody, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	var response models.ExternalListAgentsResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
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

func (client *AgentIdentityRegistryClient) ListAgentCollections() ([]*models.ExternalAgentCollectionResponse, error) {
	url := client.agentIdentityRegistrySettings.BaseUrl + LIST_AGENT_COLLECTIONS_ENDPOINT

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

	if httpResponse.StatusCode != http.StatusOK {
		if httpResponse.StatusCode == http.StatusBadRequest {
			var validationErr orchestrator_errors.ValidationError
			if err := json.NewDecoder(httpResponse.Body).Decode(&validationErr); err == nil {
				return nil, &validationErr
			}
		}
		if httpResponse.StatusCode == http.StatusNotFound {
			var notFoundErr orchestrator_errors.NotFoundError
			if err := json.NewDecoder(httpResponse.Body).Decode(&notFoundErr); err == nil {
				return nil, &notFoundErr
			}
			return nil, orchestrator_errors.NewNotFoundError("Agent not found")
		}

		responseBody, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	var response []*models.ExternalAgentCollectionResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, nil
}
