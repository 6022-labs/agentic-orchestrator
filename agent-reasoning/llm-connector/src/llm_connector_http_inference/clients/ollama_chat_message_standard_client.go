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
	OLLAMA_CHAT_MESSAGE_ENDPOINT = "/api/chat"
)

type OllamaChatMessageStandardClientInterface interface {
	ChatMessage(baseUrl, agentId string, request *models.ExternalOllamaChatMessageStandardRequest) (*models.ExternalOllamaChatMessageStandardResponse, error)
}

type OllamaChatMessageStandardClient struct {
	logger     *zap.Logger
	httpClient *http.Client
}

type newOllamaChatMessageStandardClientParams struct {
	dig.In

	Logger     *zap.Logger
	HttpClient *http.Client `name:"llm-inference-http-client"`
}

func NewOllamaChatMessageStandardClient(params newOllamaChatMessageStandardClientParams) *OllamaChatMessageStandardClient {
	return &OllamaChatMessageStandardClient{
		logger:     params.Logger,
		httpClient: params.HttpClient,
	}
}

func (c *OllamaChatMessageStandardClient) ChatMessage(baseUrl, agentId string, request *models.ExternalOllamaChatMessageStandardRequest) (*models.ExternalOllamaChatMessageStandardResponse, error) {
	url := fmt.Sprintf("%s%s", baseUrl, OLLAMA_CHAT_MESSAGE_ENDPOINT)

	// Convert the request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		c.logger.Error("Failed to marshal Ollama Chat Message request",
			zap.String("agentId", agentId),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBody))
	if err != nil {
		c.logger.Error("Failed to create Ollama Chat Message HTTP request",
			zap.String("agentId", agentId),
			zap.String("url", url),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	c.logger.Info("Sending Ollama Chat Message request",
		zap.String("agentId", agentId),
		zap.String("url", url),
		zap.String("model", request.Model),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to send Ollama Chat Message request",
			zap.String("agentId", agentId),
			zap.String("url", url),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		respBody, _ := io.ReadAll(resp.Body)
		c.logger.Error("Ollama Chat Message request failed with error status",
			zap.String("agentId", agentId),
			zap.String("url", url),
			zap.Int("status_code", resp.StatusCode),
			zap.ByteString("response_body", respBody),
		)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	c.logger.Info("Ollama Chat Message request completed successfully",
		zap.String("agentId", agentId),
		zap.String("url", url),
		zap.Int("status_code", resp.StatusCode),
	)

	var response models.ExternalOllamaChatMessageStandardResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.logger.Error("Failed to decode Ollama Chat Message response",
			zap.String("agentId", agentId),
			zap.String("url", url),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
