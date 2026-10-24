package interfaces

import "context"

// AgentAccessVerifierInterface guards agent-scoped operations to the NFT owner or the agent's own evm wallet.
// EnsureAccess fails with NotFoundError (unknown agent), ForbiddenError (neither) or UnavailableError (registry).
type AgentAccessVerifierInterface interface {
	EnsureAccess(
		ctx context.Context,
		callerId string,
		agentCollectionAddress string,
		agentCollectionTokenId uint64,
	) error
}
