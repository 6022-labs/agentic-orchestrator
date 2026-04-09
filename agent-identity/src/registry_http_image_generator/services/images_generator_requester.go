package services

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_image_generator/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_image_generator/models"
)

type ImagesGeneratorRequester struct {
	client clients.ImagesGeneratorClientInterface
}

func NewImagesGeneratorRequester(client clients.ImagesGeneratorClientInterface) *ImagesGeneratorRequester {
	return &ImagesGeneratorRequester{client: client}
}

func (r *ImagesGeneratorRequester) GenerateImages(request *interfaces.GenerateImagesRequests) (map[string][]byte, error) {
	response, err := r.client.GenerateImages(&models.ExternalGenerateImagesRequest{
		Name:              request.Name,
		Role:              request.Role,
		Creator:           request.Creator,
		EnsDomain:         request.EnsDomain,
		OriginEnsDomain:   request.OriginEnsDomain,
		OriginGeneration:  request.OriginGeneration,
		BaseImageBytes:    request.BaseImageBytes,
		BaseImageMimeType: request.BaseImageMimeType,
		Wallets:           request.Wallets,
		Stars:             request.Stars,
	})
	if err != nil {
		return nil, err
	}
	return response.Images, nil
}
