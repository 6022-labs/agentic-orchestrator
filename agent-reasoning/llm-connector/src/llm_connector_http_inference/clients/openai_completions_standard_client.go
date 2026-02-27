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
	OPEN_AI_CHAT_COMPLETION_ENDPOINT = "/v1/chat/completions"
)

type OpenAICompletionsStandardClientInterface interface {
	ChatCompletion(baseUrl, apiKey string, request *models.ExternalOpenAICompletionsStandardRequest) (*models.ExternalOpenAICompletionsStandardResponse, error)
}

type OpenAICompletionsStandardClient struct {
	logger     *zap.Logger
	httpClient *http.Client
}

type newOpenAICompletionsStandardClientParams struct {
	dig.In

	Logger     *zap.Logger
	HttpClient *http.Client `name:"llm-inference-http-client"`
}

func NewOpenAICompletionsStandardClient(params newOpenAICompletionsStandardClientParams) *OpenAICompletionsStandardClient {
	return &OpenAICompletionsStandardClient{
		logger:     params.Logger,
		httpClient: params.HttpClient,
	}
}

func (c *OpenAICompletionsStandardClient) ChatCompletion(baseUrl, apiKey string, request *models.ExternalOpenAICompletionsStandardRequest) (*models.ExternalOpenAICompletionsStandardResponse, error) {
	url := fmt.Sprintf("%s%s", baseUrl, OPEN_AI_CHAT_COMPLETION_ENDPOINT)

	// Convert the request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	c.logger.Debug("OpenAI ChatCompletion request",
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
		c.logger.Warn("OpenAI ChatCompletion failed response",
			zap.String("url", url),
			zap.Int("status_code", resp.StatusCode),
			zap.ByteString("response_body", respBody),
		)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response models.ExternalOpenAICompletionsStandardResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
