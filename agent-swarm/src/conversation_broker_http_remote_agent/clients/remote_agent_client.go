package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_remote_agent/models"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const (
	WELL_KNOWN_6022_PATH  = "/.well-known/6022"
	PaymentRequiredHeader = "PAYMENT-REQUIRED"
	wellKnown6022Protocol = "6022"
)

// RemoteAgentClientInterface is pure transport
type RemoteAgentClientInterface interface {
	GetWellKnownCard(agentUrl string) (*models.ExternalWellKnownAgent6022, error)
	CreateResponse(endpointUrl, ensDomain string, request *models.ExternalCreateResponseRequest, headers *models.ExternalRequestHeaders) (*models.ExternalResponse, error)
	CreateChatCompletion(endpointUrl, ensDomain string, request *models.ExternalChatCompletionRequest, headers *models.ExternalRequestHeaders) (*models.ExternalChatCompletionResponse, error)
}

type RemoteAgentClient struct {
	logger              *zap.Logger
	httpClient          *http.Client
	wellKnownHttpClient *http.Client
}

type NewRemoteAgentClientParams struct {
	dig.In

	Logger              *zap.Logger
	HttpClient          *http.Client `name:"agent-node-http-client"`
	WellKnownHttpClient *http.Client `name:"agent-node-well-known-http-client"`
}

func NewRemoteAgentClient(params NewRemoteAgentClientParams) *RemoteAgentClient {
	return &RemoteAgentClient{
		logger:              params.Logger,
		httpClient:          params.HttpClient,
		wellKnownHttpClient: params.WellKnownHttpClient,
	}
}

func (client *RemoteAgentClient) GetWellKnownCard(agentUrl string) (*models.ExternalWellKnownAgent6022, error) {
	url := strings.TrimRight(agentUrl, "/") + WELL_KNOWN_6022_PATH

	httpResponse, err := client.wellKnownHttpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(body))
	}

	var card models.ExternalWellKnownAgent6022
	if err := json.NewDecoder(httpResponse.Body).Decode(&card); err != nil {
		return nil, fmt.Errorf("failed to decode 6022 card: %w", err)
	}

	if card.Protocol != wellKnown6022Protocol {
		return nil, fmt.Errorf("unexpected protocol %q in 6022 card", card.Protocol)
	}

	return &card, nil
}

func (client *RemoteAgentClient) CreateResponse(endpointUrl, ensDomain string, request *models.ExternalCreateResponseRequest, headers *models.ExternalRequestHeaders) (*models.ExternalResponse, error) {
	var response models.ExternalResponse
	if err := client.post(endpointUrl, ensDomain, headers, request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (client *RemoteAgentClient) CreateChatCompletion(endpointUrl, ensDomain string, request *models.ExternalChatCompletionRequest, headers *models.ExternalRequestHeaders) (*models.ExternalChatCompletionResponse, error) {
	var response models.ExternalChatCompletionResponse
	if err := client.post(endpointUrl, ensDomain, headers, request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// setRequestHeaders maps every non-empty `header`-tagged field to its HTTP header.
func (client *RemoteAgentClient) setRequestHeaders(httpRequest *http.Request, headers *models.ExternalRequestHeaders) {
	if headers == nil {
		return
	}

	value := reflect.ValueOf(*headers)
	for i := range value.NumField() {
		headerName := value.Type().Field(i).Tag.Get("header")
		if headerName == "" {
			continue
		}
		if headerValue, ok := value.Field(i).Interface().(string); ok && headerValue != "" {
			httpRequest.Header.Set(headerName, headerValue)
		}
	}
}

func (client *RemoteAgentClient) post(url, ensDomain string, headers *models.ExternalRequestHeaders, requestBody, responseBody any) error {
	payload, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	client.setRequestHeaders(httpRequest, headers)

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode == http.StatusPaymentRequired {
		return cb_errors.NewPaymentRequiredError(ensDomain, httpResponse.Header.Get(PaymentRequiredHeader))
	}

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(body))
	}

	if err := json.NewDecoder(httpResponse.Body).Decode(responseBody); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
