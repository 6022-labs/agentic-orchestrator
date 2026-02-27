package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services/interfaces"
	gateway_interfaces "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/responses"
)

type FinishJwtAuthenticationInterface interface {
	Finish(request *requests.FinishJwtAuthenticationRequest) (*responses.FinishJwtAuthenticationResponse, error)
}

type FinishJwtAuthentication struct {
	cacheStorage           interfaces.CacheStorageInterface
	jwtManager             gateway_interfaces.JwtManagerInterface
	siweDetailsExtractor   gateway_interfaces.SIWEDetailsExtractorInterface
	signedMessageValidator gateway_interfaces.SignedMessageValidatorInterface
}

func NewFinishJwtAuthentication(
	cacheStorage interfaces.CacheStorageInterface,
	jwtManager gateway_interfaces.JwtManagerInterface,
	siweDetailsExtractor gateway_interfaces.SIWEDetailsExtractorInterface,
	signedMessageValidator gateway_interfaces.SignedMessageValidatorInterface,
) *FinishJwtAuthentication {
	return &FinishJwtAuthentication{
		jwtManager:             jwtManager,
		cacheStorage:           cacheStorage,
		siweDetailsExtractor:   siweDetailsExtractor,
		signedMessageValidator: signedMessageValidator,
	}
}

func (f *FinishJwtAuthentication) Finish(request *requests.FinishJwtAuthenticationRequest) (*responses.FinishJwtAuthenticationResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	address, chainId, nonce, err := f.siweDetailsExtractor.ExtractSIWEMessageDetails(request.Message)
	if err != nil {
		return nil, err
	}

	// Ensure nonce exists in Redis (prevents replay attacks)
	nonceRedisPath := fmt.Sprintf(NonceLoginRedisPathTemplate, *address)
	storedNonce, err := f.cacheStorage.Read(nonceRedisPath)
	if err != nil || storedNonce != *nonce {
		return nil, fmt.Errorf("invalid nonce or nonce expired")
	}

	// Verify the full SIWE message signature
	if !f.signedMessageValidator.Validate(request.Message, request.Signature, *address) {
		return nil, fmt.Errorf("invalid SIWE signature")
	}

	// Generate a JWT token
	jwt, err := f.jwtManager.CreateJwtToken(*address, *chainId)
	if err != nil {
		return nil, err
	}

	// Delete nonce after successful login (prevents reuse)
	_ = f.cacheStorage.Delete(nonceRedisPath)

	return &responses.FinishJwtAuthenticationResponse{Jwt: jwt}, nil
}
