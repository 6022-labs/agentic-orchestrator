package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
)

type VerifyWalletForBridgeInterface interface {
	Verify(request *requests.VerifyWalletForBridgeRequest) error
}

type VerifyWalletForBridge struct {
	signedMessageValidator  interfaces.SignedMessageValidatorInterface
	platformBridgeRequester interfaces.PlatformBridgeRequesterInterface
}

func NewVerifyWalletForBridge(
	signedMessageValidator interfaces.SignedMessageValidatorInterface,
	platformBridgeRequester interfaces.PlatformBridgeRequesterInterface,
) *VerifyWalletForBridge {
	return &VerifyWalletForBridge{
		signedMessageValidator:  signedMessageValidator,
		platformBridgeRequester: platformBridgeRequester,
	}
}

func (uc *VerifyWalletForBridge) Verify(request *requests.VerifyWalletForBridgeRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	if !uc.signedMessageValidator.Validate(request.Challenge, request.Signature, request.AuthenticatedWallet.Hex()) {
		return errors.NewValidationError("signature", "invalid signature")
	}

	err := uc.platformBridgeRequester.SetWalletOwnershipToUser(request.Platform, request.Challenge, request.AuthenticatedWallet.Hex())
	if err != nil {
		return fmt.Errorf("failed to set wallet ownership to user: %w", err)
	}

	return nil
}
