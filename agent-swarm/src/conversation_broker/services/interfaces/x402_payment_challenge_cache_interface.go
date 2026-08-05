package interfaces

import "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"

// X402PaymentChallengeCacheInterface remembers each node's last 402 challenge
// so later calls can pre-sign the free-call payload and skip the 402 roundtrip.
type X402PaymentChallengeCacheInterface interface {
	Read(ensDomain string) (*models.X402PaymentChallenge, bool)
	Store(ensDomain string, paymentChallenge *models.X402PaymentChallenge)
}
