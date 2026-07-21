package interfaces

import "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"

// X402PaymentPayloadSignerInterface signs x402 payloads with the orchestrator wallet;
// the node honors amount "0" only when a free rule shares network+payTo(+asset) with an accepts entry.
type X402PaymentPayloadSignerInterface interface {
	// Sign signs an exact payload of the given atomic-unit amount for the first supported accepts entry.
	Sign(paymentChallenge *models.X402PaymentChallenge, amount string) (*models.X402PaymentPayload, error)
}
