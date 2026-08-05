package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/ethereum/go-ethereum/common"
)

type VerifyAuthInterface interface {
	Execute(request *requests.VerifyAuthRequest) (*responses.VerifyAuthResponse, error)
}

type VerifyAuth struct {
	nonceStore         interfaces.NonceStoreInterface
	siweMessageBuilder interfaces.SiweMessageBuilderInterface
	signatureVerifier  interfaces.SignatureVerifierInterface
	jwtManager         interfaces.JwtManagerInterface
}

func NewVerifyAuth(
	nonceStore interfaces.NonceStoreInterface,
	siweMessageBuilder interfaces.SiweMessageBuilderInterface,
	signatureVerifier interfaces.SignatureVerifierInterface,
	jwtManager interfaces.JwtManagerInterface,
) *VerifyAuth {
	return &VerifyAuth{
		nonceStore:         nonceStore,
		siweMessageBuilder: siweMessageBuilder,
		signatureVerifier:  signatureVerifier,
		jwtManager:         jwtManager,
	}
}

func (uc *VerifyAuth) Execute(request *requests.VerifyAuthRequest) (*responses.VerifyAuthResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	address := common.HexToAddress(request.Address)
	entry := uc.nonceStore.Consume(address)
	if entry == nil {
		return nil, errors.NewUnauthorizedError("challenge_not_found", "challenge not found or expired")
	}

	message := uc.siweMessageBuilder.Build(address, entry.ChainId, entry.Nonce, entry.IssuedAt, entry.ExpiresAt)
	recovered, err := uc.signatureVerifier.Recover(message, request.Signature)
	if err != nil || recovered != address {
		return nil, errors.NewUnauthorizedError("signature_mismatch", "signature does not match address")
	}

	token, expiresAt, err := uc.jwtManager.CreateJwtToken(request.Address)
	if err != nil {
		return nil, err
	}

	return &responses.VerifyAuthResponse{Token: token, ExpiresAt: expiresAt}, nil
}
