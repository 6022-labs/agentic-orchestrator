package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_http_inference/models"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const (
	ANTHROPIC_MESSAGES_ENDPOINT = "/v1/messages"
)

type AnthropicMessagesStandardClientInterface interface {
	Messages(baseUrl, apiKey, agentId string, request *models.ExternalAnthropicMessagesStandardRequest) (*models.ExternalAnthropicMessagesStandardResponse, error)
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

func (c *AnthropicMessagesStandardClient) Messages(baseUrl, apiKey, agentId string, request *models.ExternalAnthropicMessagesStandardRequest) (*models.ExternalAnthropicMessagesStandardResponse, error) {
	url := fmt.Sprintf("%s%s", baseUrl, ANTHROPIC_MESSAGES_ENDPOINT)

	// Convert the request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		c.logger.Error("Failed to marshal Anthropic Messages request",
			zap.String("agentId", agentId),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBody))
	if err != nil {
		c.logger.Error("Failed to create Anthropic Messages HTTP request",
			zap.String("agentId", agentId),
			zap.String("url", url),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("anthropic-beta", "mcp-client-2025-04-04")

	c.logger.Info("Sending Anthropic Messages request",
		zap.String("agentId", agentId),
		zap.String("url", url),
		zap.String("model", request.Model),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to send Anthropic Messages request",
			zap.String("agentId", agentId),
			zap.String("url", url),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		respBody, _ := io.ReadAll(resp.Body)
		c.logger.Error("Anthropic Messages request failed with error status",
			zap.String("agentId", agentId),
			zap.String("url", url),
			zap.Int("status_code", resp.StatusCode),
			zap.ByteString("response_body", respBody),
		)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	c.logger.Info("Anthropic Messages request completed successfully",
		zap.String("agentId", agentId),
		zap.String("url", url),
		zap.Int("status_code", resp.StatusCode),
	)

	var response models.ExternalAnthropicMessagesStandardResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.logger.Error("Failed to decode Anthropic Messages response",
			zap.String("agentId", agentId),
			zap.String("url", url),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
