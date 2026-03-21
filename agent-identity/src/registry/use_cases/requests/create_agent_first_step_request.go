package requests

import (
	"mime/multipart"
	"regexp"
	"strings"

	registry_errors "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/ethereum/go-ethereum/common"
)

type CreateAgentFirstStepRequest struct {
	AuthenticatedWalletRaw string                `form:"authenticatedWallet" json:"authenticatedWallet"`
	AuthenticatedWallet    common.Address        `json:"-"`
	CollectionAddress      string                `form:"collectionAddress" json:"collectionAddress"`
	Name                   string                `form:"name" json:"name"`
	Role                   enums.AgentRole       `form:"role" json:"role"`
	WalletAddress          *string               `form:"walletAddress" json:"walletAddress"`
	Mission                string                `form:"mission" json:"mission"`
	CloneOfAgentId         *uint64               `form:"cloneOfAgentId" json:"cloneOfAgentId"`
	Image                  *multipart.FileHeader `form:"image" json:"-"`
}

var agentNameRegex = regexp.MustCompile(`^[a-z0-9-]+$`)

func (r *CreateAgentFirstStepRequest) Validate() error {
	if len(strings.TrimSpace(r.Name)) < 1 {
		return registry_errors.NewValidationError("name", "must be at least 1 character long")
	}
	if !agentNameRegex.MatchString(r.Name) {
		return registry_errors.NewValidationError("name", "must only contain lowercase letters, numbers, and dashes")
	}
	if !enums.GetAllAgentRoles().Contains(r.Role) {
		return registry_errors.NewValidationError("role", "must be one of: "+enums.GetAllAgentRoles().String())
	}
	if r.Role == enums.AgentRoleHuman {
		if r.WalletAddress == nil {
			return registry_errors.NewValidationError("walletAddress", "is required for human agents")
		}
		if !common.IsHexAddress(*r.WalletAddress) {
			return registry_errors.NewValidationError("walletAddress", "must be a valid Ethereum address")
		}
	}
	if r.Role == enums.AgentRoleClone && r.CloneOfAgentId == nil {
		return registry_errors.NewValidationError("cloneOfAgentId", "is required for clone agents")
	}
	if r.Role != enums.AgentRoleHuman && r.WalletAddress != nil {
		return registry_errors.NewValidationError("walletAddress", "can only be set for human agents")
	}
	if len(strings.TrimSpace(r.Mission)) < 1 {
		return registry_errors.NewValidationError("mission", "must be at least 1 character long")
	}
	if r.Image == nil {
		return registry_errors.NewValidationError("image", "cannot be empty")
	}
	if r.Image.Size > 5*1024*1024 {
		return registry_errors.NewValidationError("image", "must be less than 5MB")
	}
	if !common.IsHexAddress(r.CollectionAddress) {
		return registry_errors.NewValidationError("collectionAddress", "must be a valid Ethereum address")
	}
	if r.AuthenticatedWallet == (common.Address{}) {
		return registry_errors.NewValidationError("authenticatedWallet", "is required")
	}
	return nil
}
