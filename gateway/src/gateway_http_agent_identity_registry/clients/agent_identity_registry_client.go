package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
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
	LIST_AGENTS_ENDPOINT                     = "/agents"
	GET_AGENT_ENDPOINT                       = "/agent/%s/%d"
	LIST_AGENT_COLLECTIONS_ENDPOINT          = "/agent-collections"
	IS_INTERNAL_AGENT_ENDPOINT               = "/agent/%s/%d/is-internal"
	CREATE_AGENT_FIRST_STEP                  = "/agent/create/first-step"
	CREATE_AGENT_SECOND_STEP                 = "/agent/create/second-step"
	GET_AGENT_BY_ENS_DOMAIN_ENDPOINT = "/agent/ens/%s"
)

type AgentIdentityRegistryClientInterface interface {
	ListAgents(request *models.ExternalListAgentsRequest) (*models.ExternalListAgentsResponse, error)
	GetAgent(agentCollectionAddress string, agentCollectionTokenId uint64) (*models.ExternalAgentResponse, error)
	ListAgentCollections() ([]*models.ExternalAgentCollectionResponse, error)
	IsInternalAgent(agentCollectionAddress string, agentCollectionTokenId uint64) (*models.ExternalIsInternalAgentResponse, error)
	CreateAgentFirstStep(request *models.ExternalCreateAgentFirstStepRequest) (*models.ExternalAgentResponse, error)
	CreateAgentSecondStep(request *models.ExternalCreateAgentSecondStepRequest) (*models.ExternalCreateAgentSecondStepResponse, error)
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
	if request.IsInternal != nil {
		q.Set("isInternal", strconv.FormatBool(*request.IsInternal))
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

func (client *AgentIdentityRegistryClient) IsInternalAgent(agentCollectionAddress string, agentCollectionTokenId uint64) (*models.ExternalIsInternalAgentResponse, error) {
	url := fmt.Sprintf(client.agentIdentityRegistrySettings.BaseUrl+IS_INTERNAL_AGENT_ENDPOINT, agentCollectionAddress, agentCollectionTokenId)

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

	var response models.ExternalIsInternalAgentResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func (client *AgentIdentityRegistryClient) CreateAgentFirstStep(request *models.ExternalCreateAgentFirstStepRequest) (*models.ExternalAgentResponse, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("authenticatedWallet", request.AuthenticatedWallet); err != nil {
		return nil, fmt.Errorf("failed to write field authenticatedWallet: %w", err)
	}
	if err := writer.WriteField("collectionAddress", request.CollectionAddress); err != nil {
		return nil, fmt.Errorf("failed to write field collectionAddress: %w", err)
	}
	if err := writer.WriteField("name", request.Name); err != nil {
		return nil, fmt.Errorf("failed to write field name: %w", err)
	}
	if err := writer.WriteField("role", string(request.Role)); err != nil {
		return nil, fmt.Errorf("failed to write field role: %w", err)
	}
	if request.WalletAddress != nil {
		if err := writer.WriteField("walletAddress", *request.WalletAddress); err != nil {
			return nil, fmt.Errorf("failed to write field walletAddress: %w", err)
		}
	}
	if err := writer.WriteField("mission", request.Mission); err != nil {
		return nil, fmt.Errorf("failed to write field mission: %w", err)
	}
	if request.CloneOfAgentId != nil {
		if err := writer.WriteField("cloneOfAgentId", strconv.FormatUint(*request.CloneOfAgentId, 10)); err != nil {
			return nil, fmt.Errorf("failed to write field cloneOfAgentId: %w", err)
		}
	}

	fileContent, err := request.Image.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer fileContent.Close()

	fileWriter, err := writer.CreateFormFile("image", request.Image.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create image multipart field: %w", err)
	}

	if _, err := io.Copy(fileWriter, fileContent); err != nil {
		return nil, fmt.Errorf("failed to copy image file content: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	httpRequest, err := http.NewRequest(http.MethodPost, client.agentIdentityRegistrySettings.BaseUrl+CREATE_AGENT_FIRST_STEP, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", writer.FormDataContentType())

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResponse.Body.Close()

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

func (client *AgentIdentityRegistryClient) CreateAgentSecondStep(request *models.ExternalCreateAgentSecondStepRequest) (*models.ExternalCreateAgentSecondStepResponse, error) {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpRequest, err := http.NewRequest(
		http.MethodPost,
		client.agentIdentityRegistrySettings.BaseUrl+CREATE_AGENT_SECOND_STEP,
		bytes.NewBuffer(requestBytes),
	)
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

	var response models.ExternalCreateAgentSecondStepResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
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
