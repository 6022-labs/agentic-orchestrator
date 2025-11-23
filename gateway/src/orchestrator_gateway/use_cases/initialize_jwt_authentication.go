package use_cases

import (
	"fmt"
	"time"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/services/interfaces"
	gateway_interfaces "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/responses"
)

const (
	NonceLoginRedisPathTemplate = "login/nonces/%s"
)

type InitializeJwtAuthenticationInterface interface {
	Init(request *requests.InitializeJwtAuthenticationRequest) (*responses.InitializeJwtAuthenticationResponse, error)
}

type InitializeJwtAuthentication struct {
	cacheStorage   interfaces.CacheStorageInterface
	nonceGenerator gateway_interfaces.NonceGeneratorInterface
}

func NewInitializeJwtAuthentication(
	cacheStorage interfaces.CacheStorageInterface,
	nonceGenerator gateway_interfaces.NonceGeneratorInterface,
) *InitializeJwtAuthentication {
	return &InitializeJwtAuthentication{
		cacheStorage:   cacheStorage,
		nonceGenerator: nonceGenerator,
	}
}

func (g *InitializeJwtAuthentication) Init(request *requests.InitializeJwtAuthenticationRequest) (*responses.InitializeJwtAuthenticationResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	nonce, err := g.nonceGenerator.GenerateNonce()
	if err != nil {
		return nil, err
	}

	nonceRedisPath := fmt.Sprintf(NonceLoginRedisPathTemplate, request.Address)

	duration := time.Minute * 5
	err = g.cacheStorage.Store(nonceRedisPath, *nonce, &duration)
	if err != nil {
		return nil, err
	}

	return &responses.InitializeJwtAuthenticationResponse{Nonce: *nonce}, nil
}
