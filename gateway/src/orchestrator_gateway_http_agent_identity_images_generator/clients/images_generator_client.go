package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway_http_agent_identity_images_generator/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway_http_agent_identity_images_generator/settings"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const (
	GENERATE_IMAGES_ENDPOINT = "/generate-images"
)

type ImagesGeneratorClientInterface interface {
	GenerateImages(request *models.ExternalGenerateImagesRequest) (*models.ExternalGenerateImagesResponse, error)
}

type ImagesGeneratorClient struct {
	logger                  *zap.Logger
	httpClient              *http.Client
	imagesGeneratorSettings *settings.ImagesGeneratorSettings
}

type newImagesGeneratorClientParams struct {
	dig.In

	Logger                  *zap.Logger
	HttpClient              *http.Client `name:"images-generator-http-client"`
	ImagesGeneratorSettings *settings.ImagesGeneratorSettings
}

func NewImagesGeneratorClient(
	params newImagesGeneratorClientParams,
) *ImagesGeneratorClient {
	return &ImagesGeneratorClient{
		logger:                  params.Logger,
		httpClient:              params.HttpClient,
		imagesGeneratorSettings: params.ImagesGeneratorSettings,
	}
}

func (client *ImagesGeneratorClient) GenerateImages(request *models.ExternalGenerateImagesRequest) (*models.ExternalGenerateImagesResponse, error) {

	url := client.imagesGeneratorSettings.BaseUrl + GENERATE_IMAGES_ENDPOINT

	var b bytes.Buffer
	mp := multipart.NewWriter(&b)

	if err := mp.WriteField("name", request.Name); err != nil {
		return nil, fmt.Errorf("failed to write field 'name': %w", err)
	}
	if err := mp.WriteField("role", request.Role); err != nil {
		return nil, fmt.Errorf("failed to write field 'role': %w", err)
	}
	if err := mp.WriteField("creator", request.Creator); err != nil {
		return nil, fmt.Errorf("failed to write field 'creator': %w", err)
	}
	if err := mp.WriteField("ensDomain", request.EnsDomain); err != nil {
		return nil, fmt.Errorf("failed to write field 'ensDomain': %w", err)
	}
	if request.OriginEnsDomain != nil {
		if err := mp.WriteField("originEnsDomain", *request.OriginEnsDomain); err != nil {
			return nil, fmt.Errorf("failed to write field 'originEnsDomain': %w", err)
		}
	}
	if request.OriginGeneration != nil {
		if err := mp.WriteField("originGeneration", fmt.Sprintf("%d", *request.OriginGeneration)); err != nil {
			return nil, fmt.Errorf("failed to write field 'originGeneration': %w", err)
		}
	}
	if err := mp.WriteField("baseImageMimeType", request.BaseImageMimeType); err != nil {
		return nil, fmt.Errorf("failed to write field 'baseImageMimeType': %w", err)
	}
	for _, wallet := range request.Wallets {
		if err := mp.WriteField("wallets", wallet); err != nil {
			return nil, fmt.Errorf("failed to write field 'wallets': %w", err)
		}
	}
	if err := mp.WriteField("stars", fmt.Sprintf("%d", request.Stars)); err != nil {
		return nil, fmt.Errorf("failed to write field 'stars': %w", err)
	}

	if len(request.BaseImageBytes) > 0 {
		fw, err := mp.CreateFormFile("baseImageBytes", "base_image")
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}
		_, err = fw.Write(request.BaseImageBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to write base image bytes: %w", err)
		}
	}

	err := mp.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	client.logger.Debug("Sending Generate Images request (multipart)",
		zap.String("method", http.MethodPost),
		zap.String("url", url),
		zap.Int("bodyLen", b.Len()),
	)

	httpRequest, err := http.NewRequest(http.MethodPost, url, &b)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", mp.FormDataContentType())

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", httpResponse.StatusCode, string(responseBody))
	}

	var response models.ExternalGenerateImagesResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
