package services_test

import (
	"fmt"
	"strings"
	"testing"

	"encoding/hex"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

const (
	chain_id    = 80002
	public_key  = "0xd389d997Df5654b87914492E1e0C6861bfe835DE"
	private_key = "b4c2dd8c816efd854aecde7167b0a5db630854f426c50eb2b8784b9a89c9f95a"
)

type WhenValidatingSignedMessageTestSuite struct {
	sut *services.SignedMessageValidator
}

func WhenValidatingSignedMessageBeforeEach() *WhenValidatingSignedMessageTestSuite {
	sut := services.NewSignedMessageValidator()

	return &WhenValidatingSignedMessageTestSuite{
		sut: sut,
	}
}

func generateMessage(nonce string) string {
	return strings.Join([]string{
		"localhost:5173 wants you to sign in with your Ethereum account:\n",
		public_key,
		"\n\nPlease sign with your account\n\nURI: http://localhost:5173\nVersion: 1\nChain ID: ",
		fmt.Sprintf("%d", chain_id),
		"\nNonce: ",
		nonce,
		"\nIssued At: 2025-03-02T18:24:14.110Z",
	}, "")
}

func signMessage(message string) (string, error) {
	privateKey, err := crypto.HexToECDSA(private_key)
	if err != nil {
		return "", err
	}

	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(prefixedMessage))
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(signature), nil
}

func TestWhenValidatingSignedMessage(t *testing.T) {
	t.Parallel()

	t.Run("Given nonce does not match the one received", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return false", func(t *testing.T) {
			t.Parallel()

			baseNonce := "nonce"
			signedNonce := "invalid_signed_nonce"

			suite := WhenValidatingSignedMessageBeforeEach()

			baseNonceMessage := generateMessage(baseNonce)
			signedNonceMessage := generateMessage(signedNonce)
			signedNonceSignature, _ := signMessage(signedNonceMessage)

			result := suite.sut.Validate(baseNonceMessage, signedNonceSignature, public_key)

			assert.False(t, result)
		})
	})

	t.Run("Given nonce match the one received", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return true", func(t *testing.T) {
			t.Parallel()

			nonce := "nonce"

			baseNonceMessage := generateMessage(nonce)
			signedNonceMessage := generateMessage(nonce)
			signedNonceSignature, _ := signMessage(signedNonceMessage)

			suite := WhenValidatingSignedMessageBeforeEach()

			result := suite.sut.Validate(baseNonceMessage, signedNonceSignature, public_key)

			assert.True(t, result)
		})
	})
}
