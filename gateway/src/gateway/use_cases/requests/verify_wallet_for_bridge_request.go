package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/ethereum/go-ethereum/common"
)

type VerifyWalletForBridgeRequest struct {
	AuthenticatedWallet common.Address                    `swaggerignore:"true"`
	Platform            enums.BridgePlatform `json:"platform"`
	Challenge           string                            `json:"challenge"`
	Signature           string                            `json:"signature"`
}

func (r *VerifyWalletForBridgeRequest) Validate() error {
	if !enums.GetAllBridgePlatforms().Contains(r.Platform) {
		return errors.NewValidationError("platform", "invalid platform, must be one of: "+enums.GetAllBridgePlatforms().String())
	}

	if len(strings.TrimSpace(r.Challenge)) == 0 {
		return errors.NewValidationError("challenge", "challenge cannot be empty")
	}

	if len(strings.TrimSpace(r.Signature)) == 0 {
		return errors.NewValidationError("signature", "signature cannot be empty")
	}

	return nil
}
