package models_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services"
	"github.com/stretchr/testify/assert"
)

// The cross-repo vector: agentic-agent-node signs this exact card with anvil account 0.
const (
	vectorSigner    = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	vectorProtected = "eyJhbGciOiJFUzI1NksiLCJqd2siOnsiY3J2Ijoic2VjcDI1NmsxIiwia3R5IjoiRUMiLCJ4IjoiZ3hoVFcxUVFYVXA2cm1EQ" +
		"Wo4UmZsb2NZRzBfZnhpVzlHblVfcHpsXzdYVSIsInkiOiJOVWZ4SEtocFprYnk4Nnl3ampFQmF2ckNQbU1NWFJIMW4ySC05WHNOS3FVIn0sImtpZCI" +
		"6IjB4ZjM5RmQ2ZTUxYWFkODhGNkY0Y2U2YUI4ODI3Mjc5Y2ZmRmI5MjI2NiIsInR5cCI6IkpPU0UifQ"
	vectorSignature = "HoLnS3IMLIFR8cgB35NLrQbvgxvUq_NFj6-cbEVbUm4CX_RV-lSHMNiy2kaKPbyvp8htlPvGlT0dTNJkHPjhvA"

	vectorCardBody = `"protocol":"6022","protocolVersion":"0.1","name":"hermes","ensDomain":"hermes.agents6022.eth","end` +
		`points":{"responses":"https://agent.example.com/api/responses"},"paymentMethods":["l402","x402"]`
)

// TestWhenVerifyingAMappedCard closes the seam between decoding and verifying: the case-collision
// bug lived here, invisible to either side tested alone.
func TestWhenVerifyingAMappedCard(t *testing.T) {
	t.Parallel()

	verifier := services.NewAgentCardSignatureVerifier()
	servedBody := func(body string) string {
		return `{` + body + `,"signatures":[{"protected":"` + vectorProtected + `","signature":"` + vectorSignature + `"}]}`
	}

	t.Run("Given a genuinely signed card as the node served it", func(t *testing.T) {
		t.Parallel()

		t.Run("Should decode and then verify against the registry signer", func(t *testing.T) {
			t.Parallel()

			body := servedBody(vectorCardBody)
			card, err := decodeExternal(t, body).ToModel([]byte(body))

			assert.NoError(t, err)
			assert.Len(t, card.Signatures, 1)
			assert.NoError(t, verifier.Verify([]byte(body), card.Signatures, vectorSigner))
		})
	})

	t.Run("Given that card with its endpoint rewritten", func(t *testing.T) {
		t.Parallel()

		t.Run("Should decode but fail verification", func(t *testing.T) {
			t.Parallel()

			body := servedBody(
				`"protocol":"6022","protocolVersion":"0.1","name":"hermes","ensDomain":"hermes.agents6022.eth","endpoints":{"resp` +
					`onses":"https://attacker.example/api/responses"},"paymentMethods":["l402","x402"]`,
			)
			card, err := decodeExternal(t, body).ToModel([]byte(body))

			assert.NoError(t, err)
			assert.Error(t, verifier.Verify([]byte(body), card.Signatures, vectorSigner))
		})
	})

	t.Run("Given a signed card shadowed by a differently-cased signatures member", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not be waved through as unsigned", func(t *testing.T) {
			t.Parallel()

			body := `{` + vectorCardBody + `,"signatures":[{"protected":"` + vectorProtected +
				`","signature":"` + vectorSignature + `"}],"Signatures":[]}`
			card, err := decodeExternal(t, body).ToModel([]byte(body))

			assert.NoError(t, err)
			// The exact member is still seen, so the card is checked — and the extra member changes the digest.
			assert.Len(t, card.Signatures, 1)
			assert.Error(t, verifier.Verify([]byte(body), card.Signatures, vectorSigner))
		})
	})
}
