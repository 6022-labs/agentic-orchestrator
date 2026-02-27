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
	OLLAMA_CHAT_MESSAGE_ENDPOINT = "/api/chat"
)

type OllamaChatMessageStandardClientInterface interface {
	ChatMessage(baseUrl string, request *models.ExternalOllamaChatMessageStandardRequest) (*models.ExternalOllamaChatMessageStandardResponse, error)
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

func (c *OllamaChatMessageStandardClient) ChatMessage(baseUrl string, request *models.ExternalOllamaChatMessageStandardRequest) (*models.ExternalOllamaChatMessageStandardResponse, error) {
	url := fmt.Sprintf("%s%s", baseUrl, OLLAMA_CHAT_MESSAGE_ENDPOINT)

	// Convert the request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	c.logger.Debug("Ollama Chat Message request",
		zap.String("url", url),
		zap.String("request.model", request.Model),
		zap.Any("request.tools", request.Tools),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		respBody, _ := io.ReadAll(resp.Body)
		c.logger.Warn("Ollama Chat Message failed response",
			zap.String("url", url),
			zap.Int("status_code", resp.StatusCode),
			zap.ByteString("response_body", respBody),
		)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response models.ExternalOllamaChatMessageStandardResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
