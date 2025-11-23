package services

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_ipfs/clients"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_ipfs/models"
	"go.uber.org/zap"
)

type IPFSFileStorage struct {
	logger     *zap.Logger
	ipfsClient clients.IPFSClientInterface
}

func NewIPFSFileStorage(
	logger *zap.Logger,
	ipfsClient clients.IPFSClientInterface,
) *IPFSFileStorage {
	return &IPFSFileStorage{
		logger:     logger,
		ipfsClient: ipfsClient,
	}
}

func (s *IPFSFileStorage) UploadFile(fileBytes []byte, filename string) (*string, error) {
	pin := true

	response, err := s.ipfsClient.Add(fileBytes, filename, &models.ExternalAddQueryRequest{
		Pin: &pin,
	})
	if err != nil {
		return nil, err
	}

	go func() {
		err = s.ipfsClient.RoutingProvideWithoutResponse(&models.ExternalRoutingProvideRequest{
			Arg: response.Hash,
		})
		if err != nil {
			s.logger.Error("Failed to provide routing", zap.String("hash", response.Hash), zap.Error(err))
		}
	}()

	return &response.Hash, nil
}
