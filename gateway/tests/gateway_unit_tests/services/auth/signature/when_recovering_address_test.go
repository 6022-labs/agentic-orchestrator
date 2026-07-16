package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/auth/signature"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type WhenRecoveringAddressTestingSuite struct {
	sut *signature.ECDSASigner
}

func WhenRecoveringAddressBeforeEach() *WhenRecoveringAddressTestingSuite {
	return &WhenRecoveringAddressTestingSuite{sut: signature.NewECDSASigner()}
}

func TestWhenRecoveringAddress(t *testing.T) {
	t.Parallel()

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	expectedAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	hash := crypto.Keccak256([]byte("payload"))

	t.Run("Given a signature with a normalized V", func(t *testing.T) {
		t.Parallel()

		sig, err := crypto.Sign(hash, privateKey)
		require.NoError(t, err)

		t.Run("Should recover the signer address", func(t *testing.T) {
			t.Parallel()

			suite := WhenRecoveringAddressBeforeEach()

			address, err := suite.sut.RecoverAddress(hash, sig)

			assert.NoError(t, err)
			assert.Equal(t, expectedAddress, address)
		})
	})

	t.Run("Given a signature with a legacy V of 27", func(t *testing.T) {
		t.Parallel()

		sig, err := crypto.Sign(hash, privateKey)
		require.NoError(t, err)
		legacy := make([]byte, len(sig))
		copy(legacy, sig)
		legacy[64] += 27

		t.Run("Should normalize V and recover the signer address", func(t *testing.T) {
			t.Parallel()

			suite := WhenRecoveringAddressBeforeEach()

			address, err := suite.sut.RecoverAddress(hash, legacy)

			assert.NoError(t, err)
			assert.Equal(t, expectedAddress, address)
			// The caller's slice keeps its legacy V untouched.
			assert.Equal(t, sig[64]+27, legacy[64])
		})
	})

	t.Run("Given a signature with an invalid length", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRecoveringAddressBeforeEach()

			_, err := suite.sut.RecoverAddress(hash, []byte{0x01, 0x02})

			assert.ErrorContains(t, err, "invalid signature length")
		})
	})
}
