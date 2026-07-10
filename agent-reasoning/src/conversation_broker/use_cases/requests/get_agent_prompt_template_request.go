package requests

import (
	"strconv"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
)

// GetAgentPromptTemplateRequest identifies one agent-scoped template row.
type GetAgentPromptTemplateRequest struct {
	CollectionAddress string `params:"collectionAddress"`
	CollectionTokenId string `params:"collectionTokenId"`
	Key               string `params:"key"`
}

func (r *GetAgentPromptTemplateRequest) Validate() error {
	if strings.TrimSpace(r.CollectionAddress) == "" {
		return errors.NewValidationError("collectionAddress", "must not be empty")
	}
	if _, err := strconv.ParseUint(r.CollectionTokenId, 10, 64); err != nil {
		return errors.NewValidationError("collectionTokenId", "must be a valid unsigned integer")
	}
	var key enums.AgentPromptTemplateKey
	if err := key.Set(r.Key); err != nil {
		return errors.NewValidationError("key", "must be one of: "+enums.GetAllAgentPromptTemplateKeys().String())
	}
	return nil
}

func (r *GetAgentPromptTemplateRequest) ParsedCollectionTokenId() uint64 {
	v, _ := strconv.ParseUint(r.CollectionTokenId, 10, 64)
	return v
}

func (r *GetAgentPromptTemplateRequest) ParsedKey() enums.AgentPromptTemplateKey {
	var key enums.AgentPromptTemplateKey
	_ = key.Set(r.Key)
	return key
}
