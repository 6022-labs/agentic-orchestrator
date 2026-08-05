package services_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services"
	interfaces_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const recoverSignerAddress = "0xd389d997Df5654b87914492E1e0C6861bfe835DE"

type WhenRecoveringSignerTestingSuite struct {
	sut *services.SignatureVerifier

	mockEcdsaSigner *interfaces_mocks.MockECDSASignerInterface
}

func WhenRecoveringSignerBeforeEach(t *testing.T) *WhenRecoveringSignerTestingSuite {
	mockController := gomock.NewController(t)

	mockEcdsaSigner := interfaces_mocks.NewMockECDSASignerInterface(mockController)

	return &WhenRecoveringSignerTestingSuite{
		sut:             services.NewSignatureVerifier(mockEcdsaSigner),
		mockEcdsaSigner: mockEcdsaSigner,
	}
}

// prefixedHash mirrors the verifier's EIP-191 personal_sign hash.
func prefixedHash(message string) []byte {
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	return crypto.Keccak256Hash([]byte(prefixedMessage)).Bytes()
}

func TestWhenRecoveringSigner(t *testing.T) {
	t.Parallel()

	message := "localhost:5173 wants you to sign in with your Ethereum account"
	sigBytes := make([]byte, 65)
	sigBytes[64] = 1
	signatureHex := "0x" + hex.EncodeToString(sigBytes)

	t.Run("Given a recoverable signature", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the signer address", func(t *testing.T) {
			t.Parallel()

			suite := WhenRecoveringSignerBeforeEach(t)
			suite.mockEcdsaSigner.EXPECT().
				RecoverAddress(prefixedHash(message), sigBytes).
				Return(recoverSignerAddress, nil)

			address, err := suite.sut.Recover(message, signatureHex)

			assert.NoError(t, err)
			assert.Equal(t, common.HexToAddress(recoverSignerAddress), address)
		})
	})

	t.Run("Given the signer fails to recover", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRecoveringSignerBeforeEach(t)
			suite.mockEcdsaSigner.EXPECT().
				RecoverAddress(prefixedHash(message), sigBytes).
				Return("", errors.New("invalid signature"))

			_, err := suite.sut.Recover(message, signatureHex)

			assert.ErrorContains(t, err, "invalid signature")
		})
	})

	t.Run("Given the signature is not valid hex", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fail without calling the signer", func(t *testing.T) {
			t.Parallel()

			suite := WhenRecoveringSignerBeforeEach(t)

			_, err := suite.sut.Recover(message, "0xnot-hex")

			assert.ErrorContains(t, err, "not valid hex")
		})
	})
}
