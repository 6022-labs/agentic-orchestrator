package interfaces

import "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"

// AgentCardSignatureVerifierInterface verifies an A2A §8.4 detached JWS over a fetched agent card.
type AgentCardSignatureVerifierInterface interface {
	// Verify fails unless at least one signature is valid and its key derives to expectedSigner.
	// An empty slice is not an error — integrity is the node's responsibility, not ours to enforce.
	Verify(raw []byte, signatures []models.AgentCardSignature, expectedSigner string) error
}
