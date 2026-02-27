package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector_http_inference/models"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const (
	ANTHROPIC_MESSAGES_ENDPOINT = "/v1/messages"
)

type AnthropicMessagesStandardClientInterface interface {
	Messages(baseUrl, apiKey string, request *models.ExternalAnthropicMessagesStandardRequest) (*models.ExternalAnthropicMessagesStandardResponse, error)
}

type AnthropicMessagesStandardClient struct {
	logger     *zap.Logger
	httpClient *http.Client
}

type newAnthropicMessagesStandardClientParams struct {
	dig.In

	Logger     *zap.Logger
	HttpClient *http.Client `name:"llm-inference-http-client"`
}

func NewAnthropicMessagesStandardClient(params newAnthropicMessagesStandardClientParams) *AnthropicMessagesStandardClient {
	return &AnthropicMessagesStandardClient{
		logger:     params.Logger,
		httpClient: params.HttpClient,
	}
}

func (c *AnthropicMessagesStandardClient) Messages(baseUrl, apiKey string, request *models.ExternalAnthropicMessagesStandardRequest) (*models.ExternalAnthropicMessagesStandardResponse, error) {
	url := fmt.Sprintf("%s%s", baseUrl, ANTHROPIC_MESSAGES_ENDPOINT)

	// Convert the request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("anthropic-beta", "mcp-client-2025-04-04")

	c.logger.Debug("Anthropic Messages request",
		zap.String("url", url),
		zap.String("request.model", request.Model),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		respBody, _ := io.ReadAll(resp.Body)
		c.logger.Warn("Anthropic Messages failed response",
			zap.String("url", url),
			zap.Int("status_code", resp.StatusCode),
			zap.ByteString("response_body", respBody),
		)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response models.ExternalAnthropicMessagesStandardResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
