package compute

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models/enums"
)

type ConversationContextAgent struct {
	CollectionAddress string
	CollectionTokenId uint64
}

type ConversationContext struct {
	Id     string
	Agents []*ConversationContextAgent
}

type ComputeContext struct {
	Platform            enums.AgentContextMessagePlatform
	ForceConclusion     bool
	ConversationContext *ConversationContext
	PromptParameters    map[string]any
}
