package clients

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_orchestrator/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_orchestrator/settings"
	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type OrchestratorGatewayClientInterface interface {
	Completions(collectionAddress string, collectionTokenId uint64, body *models.ExternalAgentComputeRequest) (*models.ExternalAgentComputeResponses, error)
	GetAgent(request *models.ExternalAgentByCollectionTokenIdRequest) (*models.ExternalAgentResponse, error)
}

const (
	OrchestratorMainApiAgentsComputeEndpointTemplate = "/agent/%s/%d/compute"
)

type OrchestratorGatewayClient struct {
	logger                        *zap.Logger
	httpClient                    *http.Client
	orchestratorSettings          *settings.OrchestratorSettings
	ecdsaSignature                interfaces.ECDSASignatureInterface
	orchestratorPrivateKeyFetcher interfaces.OrchestratorPrivateKeyFetcherInterface
}

type newOrchestratorGatewayClientParams struct {
	dig.In

	Logger                        *zap.Logger
	Client                        *http.Client `name:"orchestrator-main-api-http-client"`
	OrchestratorSettings          *settings.OrchestratorSettings
	ECDSASignature                interfaces.ECDSASignatureInterface
	OrchestratorPrivateKeyFetcher interfaces.OrchestratorPrivateKeyFetcherInterface
}

func NewOrchestratorGatewayClient(
	params newOrchestratorGatewayClientParams,
) *OrchestratorGatewayClient {
	return &OrchestratorGatewayClient{
		logger:                        params.Logger,
		httpClient:                    params.Client,
		orchestratorSettings:          params.OrchestratorSettings,
		ecdsaSignature:                params.ECDSASignature,
		orchestratorPrivateKeyFetcher: params.OrchestratorPrivateKeyFetcher,
	}
}

func (c *OrchestratorGatewayClient) Completions(collectionAddress string, collectionTokenId uint64, body *models.ExternalAgentComputeRequest) (*models.ExternalAgentComputeResponses, error) {
	endpoint := fmt.Sprintf(OrchestratorMainApiAgentsComputeEndpointTemplate, collectionAddress, collectionTokenId)
	baseUrlWithEndpoint := fmt.Sprintf("%s%s", c.orchestratorSettings.ApiBaseUrl, endpoint)

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, baseUrlWithEndpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	// Prepare query and body for signature
	queryBytes := []byte(request.URL.Query().Encode())
	signature, err := c.addRequestSignature(queryBytes, bodyBytes)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Agentic-Orchestrator-Signature", signature)

	c.logger.Debug(
		"Sending completions request to orchestrator",
		zap.String("url", baseUrlWithEndpoint),
		zap.String("body", string(bodyBytes)),
		zap.Any("header", request.Header),
	)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)

		c.logger.Warn(
			"Failed to get completions from orchestrator",
			zap.String("url", baseUrlWithEndpoint),
			zap.Int("status_code", response.StatusCode),
			zap.String("body", string(bodyBytes)),
		)

		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var result models.ExternalAgentComputeResponses
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		bodyBytes, _ := io.ReadAll(response.Body)

		c.logger.Warn(
			"Failed to decode completions response body",
			zap.String("url", baseUrlWithEndpoint),
			zap.Int("status_code", response.StatusCode),
			zap.String("body", string(bodyBytes)),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &result, nil
}

func (c *OrchestratorGatewayClient) GetAgent(request *models.ExternalAgentByCollectionTokenIdRequest) (*models.ExternalAgentResponse, error) {
	endpoint := fmt.Sprintf("/agent/%s/%d", request.AgentCollectionAddress, request.AgentCollectionTokenId)
	baseUrlWithEndpoint := fmt.Sprintf("%s%s", c.orchestratorSettings.ApiBaseUrl, endpoint)

	req, err := http.NewRequest(http.MethodGet, baseUrlWithEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	c.logger.Debug(
		"Sending GetAgent request to orchestrator",
		zap.String("url", baseUrlWithEndpoint),
		zap.Any("header", req.Header),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.logger.Warn(
			"Failed to get agent from orchestrator",
			zap.String("url", baseUrlWithEndpoint),
			zap.Int("status_code", resp.StatusCode),
			zap.String("body", string(bodyBytes)),
		)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var agent models.ExternalAgentResponse
	err = json.NewDecoder(resp.Body).Decode(&agent)
	if err != nil {
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.logger.Warn(
			"Failed to decode GetAgent response body",
			zap.String("url", baseUrlWithEndpoint),
			zap.Int("status_code", resp.StatusCode),
			zap.String("body", string(bodyBytes)),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &agent, nil
}

func (c *OrchestratorGatewayClient) addRequestSignature(query, body []byte) (string, error) {
	payload := append([]byte(query), body...)

	privateKeyStr, err := c.orchestratorPrivateKeyFetcher.FetchPrivateKey()
	if err != nil {
		return "", err
	}

	privateKey, err := crypto.HexToECDSA(*privateKeyStr)
	if err != nil {
		return "", err
	}

	signature, err := c.ecdsaSignature.Sign(payload, privateKey)
	if err != nil {
		return "", err
	}

	encodedSignature := base64.StdEncoding.EncodeToString(signature)
	return encodedSignature, nil
}
