package use_cases

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

const signedRequestNonceCachePathTemplate = "auth/signed-requests/%s/%s"

type AuthenticateSignedRequestInterface interface {
	Execute(
		ctx context.Context,
		request *requests.AuthenticateSignedRequestRequest,
	) (*responses.AuthenticateSignedRequestResponse, error)
}

type AuthenticateSignedRequest struct {
	siweSettings                *settings.SiweSettings
	signedRequestMessageBuilder interfaces.SignedRequestMessageBuilderInterface
	signatureVerifier           interfaces.SignatureVerifierInterface
	cacheStorage                interfaces.CacheStorageInterface
}

func NewAuthenticateSignedRequest(
	siweSettings *settings.SiweSettings,
	signedRequestMessageBuilder interfaces.SignedRequestMessageBuilderInterface,
	signatureVerifier interfaces.SignatureVerifierInterface,
	cacheStorage interfaces.CacheStorageInterface,
) *AuthenticateSignedRequest {
	return &AuthenticateSignedRequest{
		siweSettings:                siweSettings,
		signedRequestMessageBuilder: signedRequestMessageBuilder,
		signatureVerifier:           signatureVerifier,
		cacheStorage:                cacheStorage,
	}
}

func (uc *AuthenticateSignedRequest) Execute(
	ctx context.Context,
	request *requests.AuthenticateSignedRequestRequest,
) (*responses.AuthenticateSignedRequestResponse, error) {
	if err := request.ValidateAndSanitize(); err != nil {
		return nil, err
	}

	timestamp, _ := strconv.ParseInt(request.Timestamp, 10, 64)
	maxAge := time.Duration(uc.siweSettings.SignedRequestMaxAgeSeconds) * time.Second
	if age := time.Since(time.Unix(timestamp, 0)); age > maxAge || age < -maxAge {
		return nil, errors.NewUnauthorizedError(
			"signed_request_expired",
			"the request timestamp is outside the allowed window",
		)
	}

	message := uc.signedRequestMessageBuilder.Build(
		request.Method,
		request.Path,
		request.Query,
		timestamp,
		request.Nonce,
		request.Body,
	)
	wallet, err := uc.signatureVerifier.Recover(message, request.Signature)
	if err != nil {
		return nil, errors.NewUnauthorizedError("signature_invalid", "the signature can't be recovered")
	}

	walletAddress := strings.ToLower(wallet.Hex())
	// Kept twice the window so a nonce outlives every timestamp that could still pass.
	stored, err := uc.cacheStorage.StoreIfAbsent(
		ctx,
		fmt.Sprintf(signedRequestNonceCachePathTemplate, walletAddress, request.Nonce),
		request.Timestamp,
		2*maxAge,
	)
	if err != nil {
		return nil, err
	}
	if !stored {
		return nil, errors.NewUnauthorizedError("nonce_already_used", "this nonce was already used")
	}

	return &responses.AuthenticateSignedRequestResponse{Wallet: walletAddress}, nil
}
