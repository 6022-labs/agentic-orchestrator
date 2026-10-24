package requests

import (
	"strconv"
	"strings"

	common_errors "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/errors"
	"github.com/ethereum/go-ethereum/common"
)

// AgentCommandArguments holds the "[collection address] [token ID]" arguments of an agent slash command.
type AgentCommandArguments struct {
	CollectionAddress string
	CollectionTokenId uint64
}

// TryParse leaves unreadable fields empty for ValidateAndSanitize to reject.
func (a *AgentCommandArguments) TryParse(text string) {
	arguments := strings.Fields(text)
	if len(arguments) != 2 {
		return
	}

	a.CollectionAddress = arguments[0]
	if collectionTokenId, err := strconv.ParseUint(arguments[1], 10, 64); err == nil {
		a.CollectionTokenId = collectionTokenId
	}
}

func (a *AgentCommandArguments) ValidateAndSanitize() error {
	// IsHexAddress alone also accepts an unprefixed address, which would never match the 0x-prefixed rows.
	if !strings.HasPrefix(a.CollectionAddress, "0x") || !common.IsHexAddress(a.CollectionAddress) {
		return common_errors.NewValidationError("collectionAddress", "must be a 0x-prefixed hex address")
	}
	// Token IDs start at 1 on-chain, so 0 means TryParse could not read one.
	if a.CollectionTokenId == 0 {
		return common_errors.NewValidationError("collectionTokenId", "must be a positive integer")
	}

	a.CollectionAddress = strings.ToLower(a.CollectionAddress)
	return nil
}
