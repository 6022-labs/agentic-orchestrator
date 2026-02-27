package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services/auth/signature"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

type WhenRecoveringPublicKeyTestingSuite struct {
	sut *signature.ECDSASigner
}

func WhenRecoveringPublicKeyBeforeEach() *WhenRecoveringPublicKeyTestingSuite {
	return &WhenRecoveringPublicKeyTestingSuite{
		sut: signature.NewECDSASigner(),
	}
}

func TestWhenRecoveringPublicKey(t *testing.T) {
	t.Parallel()

	t.Run("Given an invalid signature length", func(t *testing.T) {
		t.Parallel()

		suite := WhenRecoveringPublicKeyBeforeEach()
		data := []byte("test-data")
		invalidSignature := []byte("invalid-signature")

		address, err := suite.sut.RecoverPublicKey(data, invalidSignature)

		assert.Nil(t, address)
		assert.Error(t, err)
		assert.Equal(t, "invalid signature length", err.Error())
	})

	t.Run("Given a valid signature", func(t *testing.T) {
		t.Parallel()

		suite := WhenRecoveringPublicKeyBeforeEach()
		privateKey, err := crypto.GenerateKey()
		assert.NoError(t, err)

		data := []byte("test-data")
		signature, err := suite.sut.Sign(data, privateKey)
		assert.NoError(t, err)

		address, err := suite.sut.RecoverPublicKey(data, signature)

		expectedAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
		assert.NoError(t, err)
		assert.NotNil(t, address)
		assert.Equal(t, expectedAddress, *address)
	})
}
