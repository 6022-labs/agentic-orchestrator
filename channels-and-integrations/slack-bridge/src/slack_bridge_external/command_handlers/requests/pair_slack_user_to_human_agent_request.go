package requests

import (
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
		return common_errors.NewValidationError("slackUserId", "cannot be empty")
	}
	if !strings.HasPrefix(r.HumanAgentCollectionAddress, "0x") || !common.IsHexAddress(r.HumanAgentCollectionAddress) {
		return common_errors.NewValidationError("collectionAddress", "must be a 0x-prefixed hex address")
	}
	if r.HumanAgentCollectionTokenId == 0 {
		return common_errors.NewValidationError("collectionTokenId", "must be a positive integer")
	}

	r.HumanAgentCollectionAddress = strings.ToLower(r.HumanAgentCollectionAddress)
	return nil
}
