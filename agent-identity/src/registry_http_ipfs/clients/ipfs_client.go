package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_ipfs/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_ipfs/settings"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const (
	IpfsAddPath            = "/api/v0/add"
	IpfsRoutingProvidePath = "/api/v0/routing/provide"
)

type IPFSClientInterface interface {
	Add(fileBytes []byte, fileName string, params *models.ExternalAddQueryRequest) (*models.ExternalAddResponse, error)
	RoutingProvideWithoutResponse(params *models.ExternalRoutingProvideRequest) error
}

type IPFSClient struct {
	logger       *zap.Logger
	httpClient   *http.Client
	ipfsSettings *settings.IPFSSettings
}

type NewIPFSClientParams struct {
	dig.In

	Logger       *zap.Logger
	IPFSSettings *settings.IPFSSettings
	HTTPClient   *http.Client `name:"ipfs-http-client"`
}

func NewIPFSClient(params NewIPFSClientParams) *IPFSClient {
	return &IPFSClient{
		logger:       params.Logger,
		httpClient:   params.HTTPClient,
		ipfsSettings: params.IPFSSettings,
	}
}

// applyAuth attaches the bearer token when configured (Kubo API.Authorizations).
func (client *IPFSClient) applyAuth(req *http.Request) {
	if len(strings.TrimSpace(client.ipfsSettings.BearerToken)) != 0 {
		req.Header.Set("Authorization", "Bearer "+client.ipfsSettings.BearerToken)
	}
}

func (client *IPFSClient) Add(fileBytes []byte, fileName string, params *models.ExternalAddQueryRequest) (*models.ExternalAddResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err = part.Write(fileBytes); err != nil {
		return nil, fmt.Errorf("failed to write file content: %w", err)
	}

	if err = writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	queryString := ""
	if params != nil {
		queryString, err = params.ToQueryParam()
		if err != nil {
			return nil, fmt.Errorf("failed to build query params: %w", err)
		}
	}
	url := fmt.Sprintf("%s%s?%s", client.ipfsSettings.NodeHttpUrl, IpfsAddPath, queryString)

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client.applyAuth(req)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-200 response: %d (%s)", resp.StatusCode, string(bodyBytes))
	}

	var addResponse models.ExternalAddResponse
	if err = json.NewDecoder(resp.Body).Decode(&addResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

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
	url := fmt.Sprintf("%s%s?%s", client.ipfsSettings.NodeHttpUrl, IpfsRoutingProvidePath, queryString)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	client.applyAuth(req)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("received non-200 response: %d (%s)", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
