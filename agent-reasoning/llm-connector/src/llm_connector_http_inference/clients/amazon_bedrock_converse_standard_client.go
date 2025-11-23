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
	AMAZON_BEDROCK_CONVERSE_ENDPOINT_TEMPLATE = "/model/%s/converse"
)

type AmazonBedrockConverseStandardClientInterface interface {
	Messages(baseUrl, apiKey string, request *models.ExternalAmazonBedrockConverseStandardRequest) (*models.ExternalAmazonBedrockConverseStandardResponse, error)
}

type AmazonBedrockConverseStandardClient struct {
	logger     *zap.Logger
	httpClient *http.Client
}

type newAmazonBedrockConverseStandardClientParams struct {
	dig.In

	Logger     *zap.Logger
	HttpClient *http.Client `name:"llm-inference-http-client"`
}

func NewAmazonBedrockConverseStandardClient(params newAmazonBedrockConverseStandardClientParams) *AmazonBedrockConverseStandardClient {
	return &AmazonBedrockConverseStandardClient{
		logger:     params.Logger,
		httpClient: params.HttpClient,
	}
}

func (c *AmazonBedrockConverseStandardClient) Messages(baseUrl, apiKey string, request *models.ExternalAmazonBedrockConverseStandardRequest) (*models.ExternalAmazonBedrockConverseStandardResponse, error) {
	endpoint := fmt.Sprintf(AMAZON_BEDROCK_CONVERSE_ENDPOINT_TEMPLATE, request.ModelId)
	url := fmt.Sprintf("%s%s", baseUrl, endpoint)

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
	req.Header.Set("Authorization", "Bearer "+apiKey)

	c.logger.Debug("Amazon Bedrock Converse Messages request",
		zap.String("url", url),
		zap.String("request.modelId", request.ModelId),
		zap.Any("request.toolConfig", request.ToolConfig),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		respBody, _ := io.ReadAll(resp.Body)
		c.logger.Warn("Amazon Bedrock Converse Messages failed response",
			zap.String("url", url),
			zap.Int("status_code", resp.StatusCode),
			zap.ByteString("response_body", respBody),
		)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response models.ExternalAmazonBedrockConverseStandardResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
