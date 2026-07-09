package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type VerifyAuthInterface interface {
	Verify(request *requests.VerifyAuthRequest) (*responses.VerifyAuthResponse, error)
}

type VerifyAuth struct {
	nonceStore             interfaces.NonceStoreInterface
	signedMessageValidator interfaces.SignedMessageValidatorInterface
	jwtManager             interfaces.JwtManagerInterface
}

func NewVerifyAuth(
	nonceStore interfaces.NonceStoreInterface,
	signedMessageValidator interfaces.SignedMessageValidatorInterface,
	jwtManager interfaces.JwtManagerInterface,
) *VerifyAuth {
	return &VerifyAuth{
		nonceStore:             nonceStore,
		signedMessageValidator: signedMessageValidator,
		jwtManager:             jwtManager,
	}
}

func (uc *VerifyAuth) Verify(request *requests.VerifyAuthRequest) (*responses.VerifyAuthResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	if !uc.nonceStore.Consume(request.Nonce) {
		return nil, errors.NewValidationError("nonce", "nonce not found or expired")
	}

	if !uc.signedMessageValidator.Validate(request.Nonce, request.Signature, request.Address) {
		return nil, errors.NewValidationError("signature", "signature does not match address")
	}

	token, expiresAt, err := uc.jwtManager.CreateJwtToken(request.Address)
	if err != nil {
		return nil, err
	}

	return &responses.VerifyAuthResponse{Token: token, ExpiresAt: expiresAt}, nil
}
