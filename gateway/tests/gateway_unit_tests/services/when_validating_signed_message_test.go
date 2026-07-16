package services_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services"
	interfaces_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const signerAddress = "0xd389d997Df5654b87914492E1e0C6861bfe835DE"

type WhenValidatingSignedMessageTestSuite struct {
	sut *services.SignedMessageValidator

	mockEcdsaSigner *interfaces_mocks.MockECDSASignerInterface
}

func WhenValidatingSignedMessageBeforeEach(t *testing.T) *WhenValidatingSignedMessageTestSuite {
	mockController := gomock.NewController(t)

	mockEcdsaSigner := interfaces_mocks.NewMockECDSASignerInterface(mockController)

	return &WhenValidatingSignedMessageTestSuite{
		sut:             services.NewSignedMessageValidator(mockEcdsaSigner),
		mockEcdsaSigner: mockEcdsaSigner,
	}
}

// prefixedHash mirrors the validator's EIP-191 personal_sign hash.
func prefixedHash(message string) []byte {
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	return crypto.Keccak256Hash([]byte(prefixedMessage)).Bytes()
}

func TestWhenValidatingSignedMessage(t *testing.T) {
	t.Parallel()

	message := "localhost:5173 wants you to sign in with your Ethereum account"
	sigBytes := make([]byte, 65)
	sigBytes[64] = 1
	signatureHex := "0x" + hex.EncodeToString(sigBytes)

	t.Run("Given the recovered address matches", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return true", func(t *testing.T) {
			t.Parallel()

			suite := WhenValidatingSignedMessageBeforeEach(t)
			suite.mockEcdsaSigner.EXPECT().
				RecoverAddress(prefixedHash(message), sigBytes).
				Return(signerAddress, nil)

			assert.True(t, suite.sut.Validate(message, signatureHex, signerAddress))
		})
	})

	t.Run("Given the recovered address matches with a different casing", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return true", func(t *testing.T) {
			t.Parallel()

			suite := WhenValidatingSignedMessageBeforeEach(t)
			suite.mockEcdsaSigner.EXPECT().
				RecoverAddress(prefixedHash(message), sigBytes).
				Return(strings.ToLower(signerAddress), nil)

			assert.True(t, suite.sut.Validate(message, signatureHex, signerAddress))
		})
	})

	t.Run("Given the recovered address differs", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return false", func(t *testing.T) {
			t.Parallel()

			suite := WhenValidatingSignedMessageBeforeEach(t)
			suite.mockEcdsaSigner.EXPECT().
				RecoverAddress(prefixedHash(message), sigBytes).
				Return("0x0000000000000000000000000000000000000001", nil)

			assert.False(t, suite.sut.Validate(message, signatureHex, signerAddress))
		})
	})

	t.Run("Given the signer fails to recover", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return false", func(t *testing.T) {
			t.Parallel()

			suite := WhenValidatingSignedMessageBeforeEach(t)
			suite.mockEcdsaSigner.EXPECT().
				RecoverAddress(prefixedHash(message), sigBytes).
				Return("", errors.New("invalid signature"))

			assert.False(t, suite.sut.Validate(message, signatureHex, signerAddress))
		})
	})

	t.Run("Given the signature is not valid hex", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return false without calling the signer", func(t *testing.T) {
			t.Parallel()

			suite := WhenValidatingSignedMessageBeforeEach(t)

			assert.False(t, suite.sut.Validate(message, "0xnot-hex", signerAddress))
		})
	})
}
