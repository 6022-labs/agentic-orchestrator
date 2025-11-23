package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_ipfs/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_ipfs/settings"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type IPFSClientInterface interface {
	Add(fileBytes []byte, fileName string, params *models.ExternalAddQueryRequest) (*models.ExternalAddResponse, error)
	RoutingProvideWithoutResponse(params *models.ExternalRoutingProvideRequest) error
}

const (
	IpfsAddPath            = "/api/v0/add"
	IpfsRoutingProvidePath = "/api/v0/routing/provide"
)

type IPFSClient struct {
	logger       *zap.Logger
	httpClient   *http.Client
	ipfsSettings *settings.IPFSSettings
}

type newIPFSClientParams struct {
	dig.In

	Logger       *zap.Logger
	IPFSSettings *settings.IPFSSettings
	HTTPClient   *http.Client `name:"ipfs-http-client"`
}

func NewIPFSClient(
	params newIPFSClientParams,
) *IPFSClient {
	return &IPFSClient{
		logger:       params.Logger,
		httpClient:   params.HTTPClient,
		ipfsSettings: params.IPFSSettings,
	}
}

func (client *IPFSClient) Add(fileBytes []byte, fileName string, params *models.ExternalAddQueryRequest) (*models.ExternalAddResponse, error) {
	// Create a multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Write the byte array to the form
	_, err = part.Write(fileBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to write file content: %w", err)
	}

	client.logger.Debug("File content written to form", zap.String("fileName", fileName))

	// Close the writer
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// Construct the request URL
	queryString := ""
	if params != nil {
		queryString, err = params.ToQueryParam()
		if err != nil {
			return nil, fmt.Errorf("failed to build query params: %w", err)
		}
	}
	url := fmt.Sprintf("%s%s?%s", client.ipfsSettings.IpfsNodeHttpUrl, IpfsAddPath, queryString)

	client.logger.Debug("Constructed add request URL", zap.String("url", url))

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client.logger.Debug("HTTP request created", zap.Any("request", req))

	// Send the request
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	client.logger.Debug("HTTP request sent", zap.Int("statusCode", resp.StatusCode))

	// Parse the response
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)

		client.logger.Warn(
			"Received non-200 response",
			zap.Int("statusCode", resp.StatusCode),
			zap.String("responseBody", string(bodyBytes)),
		)

		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	var addResponse models.ExternalAddResponse
	err = json.NewDecoder(resp.Body).Decode(&addResponse)
	if err != nil {
		bodyBytes, _ := io.ReadAll(resp.Body)

		client.logger.Warn(
			"Failed to decode response",
			zap.String("responseBody", string(bodyBytes)),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	client.logger.Debug("Add called successfully", zap.Any("response", addResponse))

	return &addResponse, nil
}

func (client *IPFSClient) RoutingProvideWithoutResponse(params *models.ExternalRoutingProvideRequest) error {
	var err error

	queryString := ""
	if params != nil {
		queryString, err = params.ToQueryParam()
		if err != nil {
			return fmt.Errorf("failed to build query params: %w", err)
		}
	}
	url := fmt.Sprintf("%s%s?%s", client.ipfsSettings.IpfsNodeHttpUrl, IpfsRoutingProvidePath, queryString)

	client.logger.Debug("Constructed routing provide request URL", zap.String("url", url))

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	client.logger.Debug("HTTP request created", zap.Any("request", req))

	// Send the request
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Parse the response
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)

		client.logger.Warn(
			"Received non-200 response",
			zap.Int("statusCode", resp.StatusCode),
			zap.String("responseBody", string(bodyBytes)),
		)

		return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	client.logger.Debug("Routing provide called successfully")

	return nil
}
