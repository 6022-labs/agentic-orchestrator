package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_reasoning_agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_reasoning_agent_compute/settings"
	"github.com/6022-labs/zapext"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const (
	SAVE_LLM_CONNECTION_ENDPOINT             = "/llm/connection"
	DELETE_LLM_CONNECTION_ENDPOINT           = "/llm/connection"
	GET_LLM_CONNECTION_SETTINGS_ENDPOINT     = "/llm/connection/settings"
	GET_LLM_CONNECTION_SETUP_STATUS_ENDPOINT = "/llm/connection/setup-status"
	CHAT_COMPLETION_ENDPOINT                 = "/chat/completion"
	CREATE_MCP_CONNECTION_ENDPOINT           = "/mcp/connection"
	REMOVE_MCP_CONNECTION_ENDPOINT           = "/mcp/connection"
	LIST_MCP_CONNECTION_ENDPOINT             = "/mcp/connections"
)

type AgentReasoningComputeClientInterface interface {
	SaveLLMConnection(request *models.ExternalSaveLLMConnectionRequest) error
	RemoveLLMConnection(request *models.ExternalAgentIdRequest) error
	GetLLMConnectionSettings(request *models.ExternalAgentIdRequest) (*models.ExternalGetLLMConnectionSettingsResponse, error)
	GetLLMConnectionSetupStatus(request *models.ExternalAgentIdRequest) (*models.ExternalGetLLMConnectionSetupStatusResponse, error)
	ChatCompletion(request *models.ExternalChatCompletionRequest) (*string, error)
	CreateMCPConnection(request *models.ExternalCreateMCPConnectionRequest) error
	RemoveMCPConnection(request *models.ExternalIdRequest) error
	ListMCPConnections(request *models.ExternalListMCPConnectionsRequest) (*models.ExternalListMCPConnectionsResponse, error)
}

type AgentReasoningComputeClient struct {
	logger                        *zap.Logger
	httpClient                    *http.Client
	agentReasoningComputeSettings *settings.AgentReasoningAgentComputeSettings
}

type newAgentReasoningComputeClientParams struct {
	dig.In

	Logger                        *zap.Logger
	HttpClient                    *http.Client `name:"compute-http-client"`
	AgentReasoningComputeSettings *settings.AgentReasoningAgentComputeSettings
}

func NewAgentReasoningComputeClient(
	params newAgentReasoningComputeClientParams,
) *AgentReasoningComputeClient {
	return &AgentReasoningComputeClient{
		logger:                        params.Logger,
		httpClient:                    params.HttpClient,
		agentReasoningComputeSettings: params.AgentReasoningComputeSettings,
	}
}

func (client *AgentReasoningComputeClient) SaveLLMConnection(request *models.ExternalSaveLLMConnectionRequest) error {
	url := client.agentReasoningComputeSettings.BaseUrl + SAVE_LLM_CONNECTION_ENDPOINT

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	client.logger.Debug("Sending Save LLM Connection request",
		zap.String("method", http.MethodPost),
		zap.String("url", url),
		zap.String("agentId", request.AgentId.String()),
		zap.Any("body", request),
	)

	httpRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusCreated {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	return nil
}

func (client *AgentReasoningComputeClient) RemoveLLMConnection(request *models.ExternalAgentIdRequest) error {
	url := client.agentReasoningComputeSettings.BaseUrl + DELETE_LLM_CONNECTION_ENDPOINT

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	client.logger.Debug("Sending Remove LLM Connection request",
		zap.String("method", http.MethodDelete),
		zap.String("url", url),
		zap.String("agentId", request.AgentId.String()),
		zap.Any("body", request),
	)

	httpRequest, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusNoContent {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	return nil
}

func (client *AgentReasoningComputeClient) GetLLMConnectionSettings(request *models.ExternalAgentIdRequest) (*models.ExternalGetLLMConnectionSettingsResponse, error) {
	url := fmt.Sprintf("%s%s?agentId=%s",
		client.agentReasoningComputeSettings.BaseUrl,
		GET_LLM_CONNECTION_SETTINGS_ENDPOINT,
		request.AgentId.String(),
	)

	client.logger.Debug("Sending Get LLM Connection Settings request",
		zap.String("method", http.MethodGet),
		zap.String("url", url),
	)

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
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	var response models.ExternalGetLLMConnectionSettingsResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func (client *AgentReasoningComputeClient) GetLLMConnectionSetupStatus(request *models.ExternalAgentIdRequest) (*models.ExternalGetLLMConnectionSetupStatusResponse, error) {
	url := client.agentReasoningComputeSettings.BaseUrl + GET_LLM_CONNECTION_SETUP_STATUS_ENDPOINT
	url += fmt.Sprintf("?agentId=%s", request.AgentId)

	client.logger.Debug("Sending Get LLM Connection Setup Status request",
		zap.String("method", http.MethodGet),
		zap.String("url", url),
	)

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
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	var response models.ExternalGetLLMConnectionSetupStatusResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func (client *AgentReasoningComputeClient) ChatCompletion(request *models.ExternalChatCompletionRequest) (*string, error) {
	url := client.agentReasoningComputeSettings.BaseUrl + CHAT_COMPLETION_ENDPOINT

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	client.logger.Debug("Sending ChatCompletion request",
		zap.String("method", http.MethodPost),
		zap.String("url", url),
		zap.String("agentId", request.AgentId.String()),
		zapext.ArrayTruncateStart("messages", request.Messages, 5),
	)

	httpRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
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
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	var response string
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func (client *AgentReasoningComputeClient) CreateMCPConnection(request *models.ExternalCreateMCPConnectionRequest) error {
	url := client.agentReasoningComputeSettings.BaseUrl + CREATE_MCP_CONNECTION_ENDPOINT

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	client.logger.Debug("Sending Create MCP Connection request",
		zap.String("method", http.MethodPost),
		zap.String("url", url),
		zap.String("agentId", request.AgentId.String()),
		zap.Any("body", request),
	)

	httpRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusCreated {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	return nil
}

func (client *AgentReasoningComputeClient) RemoveMCPConnection(request *models.ExternalIdRequest) error {
	url := client.agentReasoningComputeSettings.BaseUrl + REMOVE_MCP_CONNECTION_ENDPOINT

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	client.logger.Debug("Sending Remove MCP Connection request",
		zap.String("method", http.MethodDelete),
		zap.String("url", url),
		zap.Any("body", request),
	)

	httpRequest, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusNoContent {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	return nil
}

func (client *AgentReasoningComputeClient) ListMCPConnections(request *models.ExternalListMCPConnectionsRequest) (*models.ExternalListMCPConnectionsResponse, error) {
	base, err := url.Parse(client.agentReasoningComputeSettings.BaseUrl + LIST_MCP_CONNECTION_ENDPOINT)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := base.Query()
	if request.AgentId != nil {
		q.Set("agentId", request.AgentId.String())
	}
	if request.SortBy != nil {
		q.Set("sortBy", *request.SortBy)
	}
	if request.SortDirection != nil {
		q.Set("sortDirection", *request.SortDirection)
	}
	if request.Page != nil {
		q.Set("page", strconv.FormatUint(uint64(*request.Page), 10))
	}
	if request.ElementsPerPage != nil {
		q.Set("elementsPerPage", strconv.FormatUint(uint64(*request.ElementsPerPage), 10))
	}
	if request.Search != nil {
		q.Set("search", *request.Search)
	}
	base.RawQuery = q.Encode()
	url := base.String()

	client.logger.Debug("Sending List MCP Connections request",
		zap.String("method", http.MethodGet),
		zap.String("url", url),
		zap.Any("request", request),
	)

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
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	var response models.ExternalListMCPConnectionsResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
