package requests

import (
	"fmt"
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
		return fmt.Errorf("%w: %q is not a collection address", common_errors.ErrValidation, a.CollectionAddress)
	}
	// Token IDs start at 1 on-chain, so 0 means TryParse could not read one.
	if a.CollectionTokenId == 0 {
		return fmt.Errorf("%w: token ID must be a positive integer", common_errors.ErrValidation)
	}

	a.CollectionAddress = strings.ToLower(a.CollectionAddress)
	return nil
}
