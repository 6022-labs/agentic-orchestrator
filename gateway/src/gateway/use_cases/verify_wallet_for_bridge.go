package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
)

type VerifyWalletForBridgeInterface interface {
	Verify(request *requests.VerifyWalletForBridgeRequest) error
}

type VerifyWalletForBridge struct {
	signatureVerifier       interfaces.SignatureVerifierInterface
	platformBridgeRequester interfaces.PlatformBridgeRequesterInterface
}

func NewVerifyWalletForBridge(
	signatureVerifier interfaces.SignatureVerifierInterface,
	platformBridgeRequester interfaces.PlatformBridgeRequesterInterface,
) *VerifyWalletForBridge {
	return &VerifyWalletForBridge{
		signatureVerifier:       signatureVerifier,
		platformBridgeRequester: platformBridgeRequester,
	}
}

func (uc *VerifyWalletForBridge) Verify(request *requests.VerifyWalletForBridgeRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	recovered, err := uc.signatureVerifier.Recover(request.Challenge, request.Signature)
	if err != nil || recovered != request.AuthenticatedWallet {
		return errors.NewValidationError("signature", "invalid signature")
	}

	err = uc.platformBridgeRequester.SetWalletOwnershipToUser(
		request.Platform,
		request.Challenge,
		request.AuthenticatedWallet.Hex(),
	)
	if err != nil {
		return fmt.Errorf("failed to set wallet ownership to user: %w", err)
	}

	return nil
}
