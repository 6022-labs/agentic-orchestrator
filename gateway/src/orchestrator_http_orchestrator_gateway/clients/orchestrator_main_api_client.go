package clients

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_http_orchestrator_gateway/models"
	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type OrchestratorMainApiClientInterface interface {
	CallExportAgent(domain string, collectionTokenId uint64) (*models.ExternalExportAgentResponse, error)
	CallIsInternalAgent(domain string, collectionTokenId uint64) (*models.ExternalIsInternalAgentResponse, error)
	CallAgentMetadata(domain string, collectionTokenId uint64) (*models.ExternalAgentMetadataResponse, error)
}

const (
	OrchestratorMainApiClientExportAgentPath        = "/agent/export"
	OrchestratorMainApiClientFetchAgentMetadataPath = "/agent/metadata"
	OrchestratorMainApiClientIsInternalAgentPath    = "/agent/is-internal"
)

type OrchestratorMainApiClient struct {
	logger                        *zap.Logger
	client                        *http.Client
	ecdsaSigner                   interfaces.ECDSASignerInterface
	orchestratorPrivateKeyFetcher interfaces.OrchestratorPrivateKeyFetcherInterface
}

type newOrchestratorMainApiClientParams struct {
	dig.In

	Logger                        *zap.Logger
	Client                        *http.Client `name:"orchestrator-main-api-http-client"`
	ECDSASigner                   interfaces.ECDSASignerInterface
	OrchestratorPrivateKeyFetcher interfaces.OrchestratorPrivateKeyFetcherInterface
}

func NewOrchestratorMainApiClient(
	params newOrchestratorMainApiClientParams,
) *OrchestratorMainApiClient {
	return &OrchestratorMainApiClient{
		logger:                        params.Logger,
		client:                        params.Client,
		ecdsaSigner:                   params.ECDSASigner,
		orchestratorPrivateKeyFetcher: params.OrchestratorPrivateKeyFetcher,
	}
}

func (c *OrchestratorMainApiClient) CallExportAgent(domain string, collectionTokenId uint64) (*models.ExternalExportAgentResponse, error) {
	baseUrl := c.buildBaseUrl(domain)
	baseUrlWithEndpoint := fmt.Sprintf("%s%s", baseUrl, OrchestratorMainApiClientExportAgentPath)

	body := models.ExternalExportAgentRequest{
		AgentCollectionTokenId: collectionTokenId,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, baseUrlWithEndpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	err = c.addRequestSignature(request)
	if err != nil {
		return nil, err
	}

	c.logger.Debug(
		"Sending export agent request to orchestrator AI API",
		zap.String("url", baseUrlWithEndpoint),
		zap.String("body", string(bodyBytes)),
		zap.Any("header", request.Header),
	)

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)

		c.logger.Warn(
			"Failed to export agent",
			zap.String("url", baseUrlWithEndpoint),
			zap.Int("status_code", response.StatusCode),
			zap.String("body", string(bodyBytes)),
		)

		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var result models.ExternalExportAgentResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		bodyBytes, _ := io.ReadAll(response.Body)

		c.logger.Warn(
			"Failed to decode export agent response",
			zap.String("url", baseUrlWithEndpoint),
			zap.Int("status_code", response.StatusCode),
			zap.String("body", string(bodyBytes)),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &result, nil
}

func (c *OrchestratorMainApiClient) CallIsInternalAgent(domain string, collectionTokenId uint64) (*models.ExternalIsInternalAgentResponse, error) {
	baseUrl := c.buildBaseUrl(domain)
	baseUrlWithEndpoint := fmt.Sprintf("%s%s", baseUrl, OrchestratorMainApiClientIsInternalAgentPath)

	fullUrl := fmt.Sprintf("%s?agentCollectionTokenId=%d", baseUrlWithEndpoint, collectionTokenId)
	request, err := http.NewRequest(http.MethodGet, fullUrl, nil)
	if err != nil {
		return nil, err
	}

	err = c.addRequestSignature(request)
	if err != nil {
		return nil, err
	}

	c.logger.Debug(
		"Sending is internal agent request to orchestrator AI API",
		zap.String("url", fullUrl),
		zap.Any("header", request.Header),
	)

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)

		c.logger.Warn(
			"Failed to check if agent is internal",
			zap.String("url", fullUrl),
			zap.Int("status_code", response.StatusCode),
			zap.String("body", string(bodyBytes)),
		)

		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var result models.ExternalIsInternalAgentResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		bodyBytes, _ := io.ReadAll(response.Body)

		c.logger.Warn(
			"Failed to decode is internal agent response",
			zap.String("url", fullUrl),
			zap.Int("status_code", response.StatusCode),
			zap.String("body", string(bodyBytes)),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &result, nil
}

func (c *OrchestratorMainApiClient) CallAgentMetadata(domain string, collectionTokenId uint64) (*models.ExternalAgentMetadataResponse, error) {
	baseUrl := c.buildBaseUrl(domain)
	baseUrlWithEndpoint := fmt.Sprintf("%s%s", baseUrl, OrchestratorMainApiClientFetchAgentMetadataPath)

	fullUrl := fmt.Sprintf("%s?agentCollectionTokenId=%d", baseUrlWithEndpoint, collectionTokenId)
	request, err := http.NewRequest(http.MethodGet, fullUrl, nil)
	if err != nil {
		return nil, err
	}

	err = c.addRequestSignature(request)
	if err != nil {
		return nil, err
	}

	c.logger.Debug(
		"Sending agent metadata request to orchestrator AI API",
		zap.String("url", fullUrl),
		zap.Any("header", request.Header),
	)

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)

		c.logger.Warn(
			"Failed to fetch agent metadata",
			zap.String("url", fullUrl),
			zap.Int("status_code", response.StatusCode),
			zap.String("body", string(bodyBytes)),
		)

		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var result models.ExternalAgentMetadataResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		bodyBytes, _ := io.ReadAll(response.Body)

		c.logger.Warn(
			"Failed to decode agent metadata response",
			zap.String("url", fullUrl),
			zap.Int("status_code", response.StatusCode),
			zap.String("body", string(bodyBytes)),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &result, nil
}

func (c *OrchestratorMainApiClient) buildBaseUrl(domain string) string {
	return fmt.Sprintf("https://%s", domain)
}

func (c *OrchestratorMainApiClient) addRequestSignature(req *http.Request) error {
	var requestQuery []byte
	var requestBody []byte
	var err error

	requestQuery = []byte(req.URL.Query().Encode())
	if req.Body != nil {
		requestBody, err = io.ReadAll(req.Body)
		if err != nil {
			return err
		}
	}

	payload := append([]byte(requestQuery), requestBody...)

	privateKeyStr, err := c.orchestratorPrivateKeyFetcher.FetchPrivateKey()
	if err != nil {
		return err
	}

	privateKey, err := crypto.HexToECDSA(*privateKeyStr)
	if err != nil {
		return err
	}

	signature, err := c.ecdsaSigner.Sign(payload, privateKey)
	if err != nil {
		return err
	}

	if req.Header == nil {
		req.Header = make(http.Header)
	}

	// Encode the signature in Base64
	encodedSignature := base64.StdEncoding.EncodeToString(signature)
	req.Header.Set("X-Agentic-Orchestrator-Signature", encodedSignature)

	return nil
}
