package requests

import (
	"mime/multipart"
	"regexp"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/ethereum/go-ethereum/common"
)

type CreateAgentFirstStepRequest struct {
	AuthenticatedWallet common.Address        `swaggerignore:"true"`
	CollectionAddress   string                `form:"collectionAddress"`
	Name                string                `form:"name"`
	Role                enums.AgentRole       `form:"role"`
	WalletAddress       *string               `form:"walletAddress"` // only for human agents
	Mission             string                `form:"mission"`
	CloneOfAgentId      *uint64               `form:"cloneOfAgentId"` // only for clone agents
	Image               *multipart.FileHeader `form:"image" swaggerignore:"true"`
}

var agentNameRegex = regexp.MustCompile(`^[a-z0-9-]+$`)

func (r *CreateAgentFirstStepRequest) Validate() error {
	if len(strings.TrimSpace(r.Name)) < 1 {
		return errors.NewValidationError("name", "must be at least 1 character long")
	}

	if !agentNameRegex.MatchString(r.Name) {
		return errors.NewValidationError("name", "must only contain lowercase letters, numbers, and dashes")
	}

	if !enums.GetAllAgentRoles().Contains(r.Role) {
		return errors.NewValidationError("role", "must be one of: "+enums.GetAllAgentRoles().String())
	}

	if r.Role == enums.AgentRoleHuman {
		if r.WalletAddress == nil {
			return errors.NewValidationError("walletAddress", "is required for human agents")
		}

		if !common.IsHexAddress(*r.WalletAddress) {
			return errors.NewValidationError("walletAddress", "must be a valid Ethereum address")
		}
	}

	if r.Role == enums.AgentRoleClone {
		if r.CloneOfAgentId == nil {
			return errors.NewValidationError("cloneOfAgentId", "is required for clone agents")
		}
	}

	if r.Role != enums.AgentRoleHuman && r.WalletAddress != nil {
		return errors.NewValidationError("walletAddress", "can only be set for human agents")
	}

	if len(strings.TrimSpace(r.Mission)) < 1 {
		return errors.NewValidationError("mission", "must be at least 1 character long")
	}

	if r.Image == nil {
		return errors.NewValidationError("image", "cannot be empty")
	}

	if r.Image.Size > 5*1024*1024 {
		return errors.NewValidationError("image", "must be less than 5MB")
	}

	if !common.IsHexAddress(r.CollectionAddress) {
		return errors.NewValidationError("collectionAddress", "must be a valid Ethereum address")
	}

	return nil
}
