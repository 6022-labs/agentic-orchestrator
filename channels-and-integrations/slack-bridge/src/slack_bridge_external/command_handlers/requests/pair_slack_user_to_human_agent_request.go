package requests

import (
	"fmt"
	"strconv"
	"strings"

	common_errors "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/errors"
	"github.com/ethereum/go-ethereum/common"
)

type PairSlackUserToHumanAgentRequest struct {
	SlackUserId                 string
	RequesterSlackUserId        string
	HumanAgentCollectionAddress string
	HumanAgentCollectionTokenId uint64
}

// TryParse reads "[slack user ID] [collection address] [token ID]" the same way as AgentCommandArguments.TryParse.
func (r *PairSlackUserToHumanAgentRequest) TryParse(text string) {
	arguments := strings.Fields(text)
	if len(arguments) != 3 {
		return
	}

	r.SlackUserId = arguments[0]
	r.HumanAgentCollectionAddress = arguments[1]
	if humanAgentCollectionTokenId, err := strconv.ParseUint(arguments[2], 10, 64); err == nil {
		r.HumanAgentCollectionTokenId = humanAgentCollectionTokenId
	}
}

func (r *PairSlackUserToHumanAgentRequest) ValidateAndSanitize() error {
	if len(strings.TrimSpace(r.SlackUserId)) == 0 {
		return fmt.Errorf("%w: slack user id cannot be empty", common_errors.ErrValidation)
	}
	if !strings.HasPrefix(r.HumanAgentCollectionAddress, "0x") || !common.IsHexAddress(r.HumanAgentCollectionAddress) {
		return fmt.Errorf(
			"%w: %q is not a collection address",
			common_errors.ErrValidation,
			r.HumanAgentCollectionAddress,
		)
	}
	if r.HumanAgentCollectionTokenId == 0 {
		return fmt.Errorf("%w: token ID must be a positive integer", common_errors.ErrValidation)
	}

	r.HumanAgentCollectionAddress = strings.ToLower(r.HumanAgentCollectionAddress)
	return nil
}
