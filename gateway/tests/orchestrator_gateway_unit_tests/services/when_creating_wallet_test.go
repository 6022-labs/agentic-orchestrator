package services_test

import (
	"crypto/ecdsa"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

type WhenCreatingWalletTestSuite struct {
	sut *services.WalletCreator
}

func WhenCreatingWalletBeforeEach() WhenCreatingWalletTestSuite {
	sut := services.NewWalletCreator()

	return WhenCreatingWalletTestSuite{
		sut: sut,
	}
}

func TestWhenCreatingWallet(t *testing.T) {
	t.Parallel()

	t.Run("Should return an address", func(t *testing.T) {
		t.Parallel()

		basics := WhenCreatingWalletBeforeEach()

		address, _, err := basics.sut.Create()

		assert.Nil(t, err)
		assert.NotNil(t, address)
	})

	t.Run("Should return a private key", func(t *testing.T) {
		t.Parallel()

		basics := WhenCreatingWalletBeforeEach()

		_, privateKey, err := basics.sut.Create()

		assert.Nil(t, err)
		assert.NotNil(t, privateKey)
	})

	t.Run("Private key should belong to public key", func(t *testing.T) {
		t.Parallel()

		basics := WhenCreatingWalletBeforeEach()

		address, privateKey, err := basics.sut.Create()

		assert.Nil(t, err)

		// Convert private key string back to ECDSA private key
		privKey, err := crypto.HexToECDSA(*privateKey)
		assert.Nil(t, err)

		// Extract public key from private key
		publicKey := privKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		assert.True(t, ok, "public key is not of type *ecdsa.PublicKey")

		// Derive address from public key
		derivedAddress := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

		// Assert that derived address matches the original address
		assert.Equal(t, *address, derivedAddress)
	})
}
