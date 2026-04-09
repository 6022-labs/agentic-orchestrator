package services

import (
	registry_interfaces "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_ipfs/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_ipfs/models"
	"go.uber.org/zap"
)

type IPFSFileStorage struct {
	logger     *zap.Logger
	ipfsClient clients.IPFSClientInterface
}

func NewIPFSFileStorage(logger *zap.Logger, ipfsClient clients.IPFSClientInterface) *IPFSFileStorage {
	return &IPFSFileStorage{logger: logger, ipfsClient: ipfsClient}
}

func (s *IPFSFileStorage) UploadFile(fileBytes []byte, filename string) (*string, error) {
	pin := true
	response, err := s.ipfsClient.Add(fileBytes, filename, &models.ExternalAddQueryRequest{Pin: &pin})
	if err != nil {
		return nil, err
	}

	go func(hash string) {
		err := s.ipfsClient.RoutingProvideWithoutResponse(&models.ExternalRoutingProvideRequest{Arg: hash})
		if err != nil {
			s.logger.Error("Failed to provide routing", zap.String("hash", hash), zap.Error(err))
		}
	}(response.Hash)

	return &response.Hash, nil
}

var _ registry_interfaces.IPFSFileStorageInterface = (*IPFSFileStorage)(nil)
