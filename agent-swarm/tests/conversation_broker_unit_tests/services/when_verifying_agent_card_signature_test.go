package services_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services"
	"github.com/stretchr/testify/assert"
)

// The cross-repo contract vector: agentic-agent-node's when_signing_agent_card_test.go pins
// these same literals. If either side drifts, signed cards stop verifying in production.
const (
	vectorSigner    = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	vectorProtected = "eyJhbGciOiJFUzI1NksiLCJqd2siOnsiY3J2Ijoic2VjcDI1NmsxIiwia3R5IjoiRUMiLCJ4IjoiZ3hoVFcxUVFYVXA2cm1EQ" +
		"Wo4UmZsb2NZRzBfZnhpVzlHblVfcHpsXzdYVSIsInkiOiJOVWZ4SEtocFprYnk4Nnl3ampFQmF2ckNQbU1NWFJIMW4ySC05WHNOS3FVIn0sImtpZCI" +
		"6IjB4ZjM5RmQ2ZTUxYWFkODhGNkY0Y2U2YUI4ODI3Mjc5Y2ZmRmI5MjI2NiIsInR5cCI6IkpPU0UifQ"
	vectorSignature = "HoLnS3IMLIFR8cgB35NLrQbvgxvUq_NFj6-cbEVbUm4CX_RV-lSHMNiy2kaKPbyvp8htlPvGlT0dTNJkHPjhvA"

	vectorCardBody = `"protocol":"6022","protocolVersion":"0.1","name":"hermes","ensDomain":"hermes.agents6022.eth","end` +
		`points":{"responses":"https://agent.example.com/api/responses"},"paymentMethods":["l402","x402"]`

	// Same key and card plus a member this build's structs do not model.
	vectorUnknownMemberBody = `"futureExtension":{"nested":true},` + vectorCardBody
	vectorUnknownSignature  = "ofunN24s3-ePHfSvbvbFj2nkgbHWKF8U7QAZyJDiwyRSxkKSJpgZMn7QdQu99Lk3xl1yTCibMJqlKUtCmgyn_g"
)

type WhenVerifyingAgentCardSignatureTestingSuite struct {
	sut *services.AgentCardSignatureVerifier
}

func WhenVerifyingAgentCardSignatureBeforeEach() *WhenVerifyingAgentCardSignatureTestingSuite {
	return &WhenVerifyingAgentCardSignatureTestingSuite{
		sut: services.NewAgentCardSignatureVerifier(),
	}
}

// signedCardBytes rebuilds the served card: the body plus the signatures member the signer excluded.
func signedCardBytes(body, protected, signature string) []byte {
	return []byte(fmt.Sprintf(`{%s,"signatures":[{"protected":%q,"signature":%q}]}`, body, protected, signature))
}

func vectorSignatures() []models.AgentCardSignature {
	return []models.AgentCardSignature{
		{
			Protected: vectorProtected,
			Signature: vectorSignature,
		},
	}
}

func TestWhenVerifyingAgentCardSignature(t *testing.T) {
	t.Parallel()

	t.Run("Given the card carries no signature at all", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"Should accept it — integrity is the agent node's responsibility, not ours to enforce",
			func(t *testing.T) {
				t.Parallel()

				suite := WhenVerifyingAgentCardSignatureBeforeEach()

				err := suite.sut.Verify([]byte(`{`+vectorCardBody+`}`), nil, vectorSigner)

				assert.NoError(t, err)
			},
		)
	})

	t.Run("Given the card carries more than 16 signatures", func(t *testing.T) {
		t.Parallel()

		t.Run("Should refuse it before verifying any", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAgentCardSignatureBeforeEach()

			raw := signedCardBytes(vectorCardBody, vectorProtected, vectorSignature)
			signature := models.AgentCardSignature{Protected: vectorProtected, Signature: vectorSignature}
			signatures := slices.Repeat([]models.AgentCardSignature{signature}, 17)

			assert.ErrorContains(t, suite.sut.Verify(raw, signatures, vectorSigner), "at most 16 are checked")
		})
	})

	t.Run("Given the card is not a JSON object", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAgentCardSignatureBeforeEach()

			assert.Error(t, suite.sut.Verify([]byte(`["nope"]`), vectorSignatures(), vectorSigner))
		})
	})

	t.Run("Given the card carries a member this build does not model", func(t *testing.T) {
		t.Parallel()

		t.Run("Should still verify, because the digest comes from the received bytes", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAgentCardSignatureBeforeEach()

			raw := signedCardBytes(vectorUnknownMemberBody, vectorProtected, vectorUnknownSignature)
			signatures := []models.AgentCardSignature{{Protected: vectorProtected, Signature: vectorUnknownSignature}}

			assert.NoError(t, suite.sut.Verify(raw, signatures, vectorSigner))
		})
	})

	t.Run("Given an expected signer without the 0x prefix", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAgentCardSignatureBeforeEach()

			raw := signedCardBytes(vectorCardBody, vectorProtected, vectorSignature)
			err := suite.sut.Verify(raw, vectorSignatures(), "f39fd6e51aad88f6f4ce6ab8827279cfffb92266")

			assert.ErrorContains(t, err, "not a hex address")
		})
	})

	t.Run("Given no expected signer could be resolved", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject rather than trust the card's own jwk", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAgentCardSignatureBeforeEach()

			raw := signedCardBytes(vectorCardBody, vectorProtected, vectorSignature)

			assert.Error(t, suite.sut.Verify(raw, vectorSignatures(), ""))
		})
	})

	t.Run("Given an expected signer that is not hex", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAgentCardSignatureBeforeEach()

			raw := signedCardBytes(vectorCardBody, vectorProtected, vectorSignature)
			err := suite.sut.Verify(raw, vectorSignatures(), "0xf39fd6")

			assert.ErrorContains(t, err, "not a hex address")
		})
	})

	t.Run("Given the protected header declares another algorithm", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAgentCardSignatureBeforeEach()

			// {"alg":"none","jwk":{...},"kid":"...","typ":"JOSE"}
			noneAlg := "eyJhbGciOiJub25lIiwiandrIjp7ImNydiI6InNlY3AyNTZrMSIsImt0eSI6IkVDIiwieCI6Imd4aFRXMVFRWFVwNnJtREFqOFJmbG" +
				"9jWUcwX2Z4aVc5R25VX3B6bF83WFUiLCJ5IjoiTlVmeEhLaHBaa2J5ODZ5d2pqRUJhdnJDUG1NTVhSSDFuMkgtOVhzTktxVSJ9LCJraWQiOiIweG" +
				"YzOUZkNmU1MWFhZDg4RjZGNGNlNmFCODgyNzI3OWNmZkZiOTIyNjYiLCJ0eXAiOiJKT1NFIn0"
			signatures := []models.AgentCardSignature{{Protected: noneAlg, Signature: vectorSignature}}

			assert.ErrorContains(
				t,
				suite.sut.Verify(signedCardBytes(vectorCardBody, noneAlg, vectorSignature), signatures, vectorSigner),
				"algorithm",
			)
		})
	})

	t.Run("Given the card was signed by a different key", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it even though the signature itself is valid", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAgentCardSignatureBeforeEach()

			raw := signedCardBytes(vectorCardBody, vectorProtected, vectorSignature)
			err := suite.sut.Verify(raw, vectorSignatures(), "0x70997970C51812dc3A010C7d01b50e0d17dc79C8")

			assert.ErrorContains(t, err, "expected 0x70997970c51812dc3a010c7d01b50e0d17dc79c8")
		})
	})

	t.Run("Given the card body was tampered with after signing", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAgentCardSignatureBeforeEach()

			// The endpoint redirect is exactly the attack card signing exists to stop.
			tampered := `"protocol":"6022","protocolVersion":"0.1","name":"hermes","ensDomain":"hermes.agents6022.eth","endpoi` +
				`nts":{"responses":"https://attacker.example/api/responses"},"paymentMethods":["l402","x402"]`
			err := suite.sut.Verify(
				signedCardBytes(tampered, vectorProtected, vectorSignature),
				vectorSignatures(),
				vectorSigner,
			)

			assert.Error(t, err)
		})
	})

	t.Run("Given the signature was replaced", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAgentCardSignatureBeforeEach()

			forged := "AoLnS3IMLIFR8cgB35NLrQbvgxvUq_NFj6-cbEVbUm4CX_RV-lSHMNiy2kaKPbyvp8htlPvGlT0dTNJkHPjhvA"
			signatures := []models.AgentCardSignature{{Protected: vectorProtected, Signature: forged}}

			assert.Error(
				t,
				suite.sut.Verify(signedCardBytes(vectorCardBody, vectorProtected, forged), signatures, vectorSigner),
			)
		})
	})

	t.Run("Given several signatures where only the last is valid", func(t *testing.T) {
		t.Parallel()

		t.Run("Should accept it, since keys rotate", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAgentCardSignatureBeforeEach()

			raw := signedCardBytes(vectorCardBody, vectorProtected, vectorSignature)
			signatures := []models.AgentCardSignature{
				{Protected: "not-base64!!", Signature: vectorSignature},
				{Protected: vectorProtected, Signature: vectorSignature},
			}

			assert.NoError(t, suite.sut.Verify(raw, signatures, vectorSigner))
		})
	})

	t.Run("Given the card carries the agent node's signature", func(t *testing.T) {
		t.Parallel()

		t.Run("Should accept it", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAgentCardSignatureBeforeEach()

			raw := signedCardBytes(vectorCardBody, vectorProtected, vectorSignature)
			err := suite.sut.Verify(raw, vectorSignatures(), vectorSigner)

			assert.NoError(t, err)
		})

		t.Run("Should accept it whatever order the members arrive in", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAgentCardSignatureBeforeEach()

			reordered := `"paymentMethods":["l402","x402"],"name":"hermes","protocolVersion":"0.1","endpoints":{"responses":"h` +
				`ttps://agent.example.com/api/responses"},"protocol":"6022","ensDomain":"hermes.agents6022.eth"`
			err := suite.sut.Verify(
				signedCardBytes(reordered, vectorProtected, vectorSignature),
				vectorSignatures(),
				vectorSigner,
			)

			assert.NoError(t, err)
		})

		t.Run("Should match a lowercase expected signer, since addresses are not case-sensitive", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAgentCardSignatureBeforeEach()

			raw := signedCardBytes(vectorCardBody, vectorProtected, vectorSignature)
			err := suite.sut.Verify(raw, vectorSignatures(), "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")

			assert.NoError(t, err)
		})
	})
}
