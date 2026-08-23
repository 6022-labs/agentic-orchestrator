package use_cases

import (
	"context"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
)

type VerifyWalletForBridgeInterface interface {
	Verify(ctx context.Context, request *requests.VerifyWalletForBridgeRequest) error
}

type VerifyWalletForBridge struct {
	gatewayMetrics          metrics_interfaces.GatewayMetricsInterface
	signatureVerifier       interfaces.SignatureVerifierInterface
	platformBridgeRequester interfaces.PlatformBridgeRequesterInterface
}

func NewVerifyWalletForBridge(
	gatewayMetrics metrics_interfaces.GatewayMetricsInterface,
	signatureVerifier interfaces.SignatureVerifierInterface,
	platformBridgeRequester interfaces.PlatformBridgeRequesterInterface,
) *VerifyWalletForBridge {
	return &VerifyWalletForBridge{
		gatewayMetrics:          gatewayMetrics,
		signatureVerifier:       signatureVerifier,
		platformBridgeRequester: platformBridgeRequester,
	}
}

func (uc *VerifyWalletForBridge) Verify(
	ctx context.Context,
	request *requests.VerifyWalletForBridgeRequest,
) (err error) {
	defer func() { uc.gatewayMetrics.RecordWalletVerification(ctx, err) }()

	if err := request.Validate(); err != nil {
		return err
	}

	recovered, err := uc.signatureVerifier.Recover(request.Challenge, request.Signature)
	if err != nil || recovered != request.AuthenticatedWallet {
		return errors.NewValidationError("signature", "invalid signature")
	}

	err = uc.platformBridgeRequester.SetWalletOwnershipToUser(ctx,
		request.Platform,
		request.Challenge,
		request.AuthenticatedWallet.Hex(),
	)
	if err != nil {
		return fmt.Errorf("failed to set wallet ownership to user: %w", err)
	}

	return nil
}
