package services_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

type WhenSigningWithECDSATestingSuite struct {
	sut *services.ECDSASignature
}

func WhenSigningWithECDSABeforeEach(t *testing.T) *WhenSigningWithECDSATestingSuite {
	sut := services.NewECDSASignature()

	return &WhenSigningWithECDSATestingSuite{
		sut: sut,
	}
}

func generateTestPrivateKey(t *testing.T) *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	assert.NoError(t, err)
	return privateKey
}

func TestWhenSigningWithECDSA(t *testing.T) {
	t.Parallel()

	t.Run("Given a valid payload and private key", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a valid signature", func(t *testing.T) {
			t.Parallel()

			suite := WhenSigningWithECDSABeforeEach(t)
			privateKey := generateTestPrivateKey(t)
			data := []byte("some-payload-to-sign")

			signature, err := suite.sut.Sign(data, privateKey)

			assert.NoError(t, err)
			assert.NotNil(t, signature)
			assert.Equal(t, 65, len(signature))
		})
	})

	t.Run("Given a signature is recovered", func(t *testing.T) {
		t.Parallel()

		t.Run("Should recover the correct address", func(t *testing.T) {
			t.Parallel()

			suite := WhenSigningWithECDSABeforeEach(t)
			privateKey := generateTestPrivateKey(t)
			expectedAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
			data := []byte("some-payload-to-sign")

			signature, err := suite.sut.Sign(data, privateKey)
			assert.NoError(t, err)

			recoveredAddress, err := suite.sut.RecoverPublicKey(data, signature)

			assert.NoError(t, err)
			assert.NotNil(t, recoveredAddress)
			assert.Equal(t, expectedAddress, *recoveredAddress)
		})
	})

	t.Run("Given an invalid signature length", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSigningWithECDSABeforeEach(t)
			data := []byte("some-payload-to-sign")
			invalidSignature := []byte("short")

			recoveredAddress, err := suite.sut.RecoverPublicKey(data, invalidSignature)

			assert.Error(t, err)
			assert.Nil(t, recoveredAddress)
			assert.Contains(t, err.Error(), "invalid signature length")
		})
	})

	t.Run("Given different data is used for recovery", func(t *testing.T) {
		t.Parallel()

		t.Run("Should recover a different address", func(t *testing.T) {
			t.Parallel()

			suite := WhenSigningWithECDSABeforeEach(t)
			privateKey := generateTestPrivateKey(t)
			expectedAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
			data := []byte("original-payload")
			differentData := []byte("different-payload")

			signature, err := suite.sut.Sign(data, privateKey)
			assert.NoError(t, err)

			recoveredAddress, err := suite.sut.RecoverPublicKey(differentData, signature)

			assert.NoError(t, err)
			assert.NotNil(t, recoveredAddress)
			assert.NotEqual(t, expectedAddress, *recoveredAddress)
		})
	})
}
