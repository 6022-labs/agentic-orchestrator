package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector_http_inference/models"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const (
	OPEN_AI_RESPONSES_ENDPOINT = "/v1/responses"
)

type OpenAIResponsesStandardClientInterface interface {
	Responses(baseUrl, apiKey string, request *models.ExternalOpenAIResponsesStandardRequest) (*models.ExternalOpenAIResponsesStandardResponse, error)
}

type OpenAIResponsesStandardClient struct {
	logger     *zap.Logger
	httpClient *http.Client
}

type newOpenAIResponsesStandardClientParams struct {
	dig.In

	Logger     *zap.Logger
	HttpClient *http.Client `name:"llm-inference-http-client"`
}

func NewOpenAIResponsesStandardClient(params newOpenAIResponsesStandardClientParams) *OpenAIResponsesStandardClient {
	return &OpenAIResponsesStandardClient{
		logger:     params.Logger,
		httpClient: params.HttpClient,
	}
}

func (c *OpenAIResponsesStandardClient) Responses(baseUrl, apiKey string, request *models.ExternalOpenAIResponsesStandardRequest) (*models.ExternalOpenAIResponsesStandardResponse, error) {
	url := fmt.Sprintf("%s%s", baseUrl, OPEN_AI_RESPONSES_ENDPOINT)

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

	c.logger.Debug("OpenAI Responses request",
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
		c.logger.Warn("OpenAI Responses failed response",
			zap.String("url", url),
			zap.Int("status_code", resp.StatusCode),
			zap.ByteString("response_body", respBody),
		)

		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response models.ExternalOpenAIResponsesStandardResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
