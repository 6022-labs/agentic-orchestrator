package models

import "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"

// RemoteAgentInput is one fully-prepared inference turn for a remote agent
// node: resolved instructions plus the conversation window as input messages.
type RemoteAgentInput struct {
	Instructions string
	Input        []*RemoteAgentInputMessage
	Metadata     map[string]string
}

type RemoteAgentInputMessage struct {
	Role    enums.RemoteAgentInputRole
	Content string
}
