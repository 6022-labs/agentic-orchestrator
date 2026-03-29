package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_knowledge_database_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_knowledge_database_connector/settings"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const (
	ADD_MESSAGES_TO_AGENT_CONTEXT_ENDPOINT            = "/agent/context/messages"
	GET_AGENT_CONTEXT_MESSAGES_ENDPOINT               = "/agent/context/messages"
	CREATE_AGENT_CONTEXT_ENDPOINT                     = "/agent/context"
	GET_AGENT_CONTEXTS_ENDPOINT                       = "/agent/%s/contexts"
	GET_AGENT_CONTEXT_BY_KEY_ENDPOINT                 = "/agent/%s/context/%s"
	LINK_CONVERSATION_TO_AGENT_CONTEXT_ENDPOINT       = "/agent/context/link-conversation"
	GET_LINKED_AGENT_CONTEXT_OF_CONVERSATION_ENDPOINT = "/agent/context/linked-conversation"
	GET_AGENT_PROMPT_TEMPLATE_ENDPOINT                = "/agent/prompt/template"
)

type AgentDatabaseConnectorClientInterface interface {
	AddMessagesToAgentContext(request *models.ExternalAddMessagesToAgentContextRequest) error
	GetAgentContextMessages(request *models.ExternalGetAgentContextMessagesRequest) (*models.ExternalGetAgentContextMessageResponses, error)
	CreateAgentContext(request *models.ExternalCreateAgentContextRequest) error
	GetAgentContexts(request *models.ExternalAgentIdRequest) (*models.ExternalAgentContextResponses, error)
	GetAgentContextByKey(request *models.ExternalGetAgentContextByKeyRequest) (*models.ExternalAgentContextResponse, error)
	LinkConversationToAgentContext(request *models.ExternalLinkConversationToAgentContextRequest) error
	GetLinkedAgentContextOfConversation(request *models.ExternalGetLinkedAgentContextOfConversationRequest) (*models.ExternalGetLinkedAgentContextOfConversationResponse, error)
	GetAgentPromptTemplate(request *models.ExternalGetAgentPromptTemplateRequest) (*models.ExternalGetAgentPromptTemplateResponse, error)
}

type AgentDatabaseConnectorClient struct {
	logger                                  *zap.Logger
	httpClient                              *http.Client
	agentKnowledgeDatabaseConnectorSettings *settings.AgentKnowledgeDatabaseConnectorSettings
}

type newAgentDatabaseConnectorClientParams struct {
	dig.In

	Logger                                  *zap.Logger
	HttpClient                              *http.Client `name:"database-connector-http-client"`
	AgentKnowledgeDatabaseConnectorSettings *settings.AgentKnowledgeDatabaseConnectorSettings
}

func NewAgentDatabaseConnectorClient(
	params newAgentDatabaseConnectorClientParams,
) *AgentDatabaseConnectorClient {
	return &AgentDatabaseConnectorClient{
		logger:                                  params.Logger,
		httpClient:                              params.HttpClient,
		agentKnowledgeDatabaseConnectorSettings: params.AgentKnowledgeDatabaseConnectorSettings,
	}
}

func (client *AgentDatabaseConnectorClient) AddMessagesToAgentContext(request *models.ExternalAddMessagesToAgentContextRequest) error {
	url := client.agentKnowledgeDatabaseConnectorSettings.BaseUrl + ADD_MESSAGES_TO_AGENT_CONTEXT_ENDPOINT

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	client.logger.Debug("Sending AddMessagesToAgentContext request",
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
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode > 299 {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	return nil
}

func (client *AgentDatabaseConnectorClient) GetAgentContextMessages(request *models.ExternalGetAgentContextMessagesRequest) (*models.ExternalGetAgentContextMessageResponses, error) {
	url := client.agentKnowledgeDatabaseConnectorSettings.BaseUrl + GET_AGENT_CONTEXT_MESSAGES_ENDPOINT + "?"
	url += "agentId=" + request.AgentId.String()
	url += "&key=" + request.Key

	client.logger.Debug("Sending GetAgentContextMessages request",
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

	if httpResponse.StatusCode > 299 {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	var response []*models.ExternalGetAgentContextMessageResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	agentContextMessages := models.ExternalGetAgentContextMessageResponses(response)

	return &agentContextMessages, nil
}

func (client *AgentDatabaseConnectorClient) CreateAgentContext(request *models.ExternalCreateAgentContextRequest) error {
	url := client.agentKnowledgeDatabaseConnectorSettings.BaseUrl + CREATE_AGENT_CONTEXT_ENDPOINT

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	client.logger.Debug("Sending CreateAgentContext request",
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
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode > 299 {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	return nil
}

func (client *AgentDatabaseConnectorClient) GetAgentContexts(request *models.ExternalAgentIdRequest) (*models.ExternalAgentContextResponses, error) {
	baseUrl := fmt.Sprintf(client.agentKnowledgeDatabaseConnectorSettings.BaseUrl+GET_AGENT_CONTEXTS_ENDPOINT, request.AgentId.String())
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := u.Query()
	if request.Search != nil {
		q.Set("search", *request.Search)
	}
	if request.SortBy != nil {
		q.Set("sortBy", *request.SortBy)
	}
	if request.SortDirection != nil {
		q.Set("sortDirection", *request.SortDirection)
	}
	u.RawQuery = q.Encode()
	urlStr := u.String()

	client.logger.Debug("Sending GetAgentContexts request",
		zap.String("method", http.MethodGet),
		zap.String("url", urlStr),
	)

	httpRequest, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode > 299 {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	var response []*models.ExternalAgentContextResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	agentContexts := models.ExternalAgentContextResponses(response)

	return &agentContexts, nil
}

func (client *AgentDatabaseConnectorClient) GetAgentContextByKey(request *models.ExternalGetAgentContextByKeyRequest) (*models.ExternalAgentContextResponse, error) {
	urlStr := fmt.Sprintf(
		client.agentKnowledgeDatabaseConnectorSettings.BaseUrl+GET_AGENT_CONTEXT_BY_KEY_ENDPOINT,
		request.AgentId.String(),
		url.PathEscape(request.Key),
	)

	client.logger.Debug("Sending GetAgentContextByKey request",
		zap.String("method", http.MethodGet),
		zap.String("url", urlStr),
	)

	httpRequest, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if httpResponse.StatusCode > 299 {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	var response models.ExternalAgentContextResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func (client *AgentDatabaseConnectorClient) LinkConversationToAgentContext(request *models.ExternalLinkConversationToAgentContextRequest) error {
	url := client.agentKnowledgeDatabaseConnectorSettings.BaseUrl + LINK_CONVERSATION_TO_AGENT_CONTEXT_ENDPOINT

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	client.logger.Debug("Sending LinkConversationToAgentContext request",
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
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode > 299 {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	return nil
}

func (client *AgentDatabaseConnectorClient) GetLinkedAgentContextOfConversation(request *models.ExternalGetLinkedAgentContextOfConversationRequest) (*models.ExternalGetLinkedAgentContextOfConversationResponse, error) {
	url := client.agentKnowledgeDatabaseConnectorSettings.BaseUrl + GET_LINKED_AGENT_CONTEXT_OF_CONVERSATION_ENDPOINT + "?"
	url += "agentId=" + request.AgentId.String()
	url += "&conversationId=" + request.ConversationId

	client.logger.Debug("Sending GetLinkedAgentContextOfConversation request",
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

	if httpResponse.StatusCode > 299 {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	var response models.ExternalGetLinkedAgentContextOfConversationResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func (client *AgentDatabaseConnectorClient) GetAgentPromptTemplate(request *models.ExternalGetAgentPromptTemplateRequest) (*models.ExternalGetAgentPromptTemplateResponse, error) {
	url := client.agentKnowledgeDatabaseConnectorSettings.BaseUrl + GET_AGENT_PROMPT_TEMPLATE_ENDPOINT
	url += "?agentId=" + request.AgentId.String()
	url += "&templateKey=" + request.TemplateKey.String()

	client.logger.Debug("Sending GetAgentPromptTemplate request",
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

	if httpResponse.StatusCode > 299 {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	var response models.ExternalGetAgentPromptTemplateResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
