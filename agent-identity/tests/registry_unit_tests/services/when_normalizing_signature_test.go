package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/stretchr/testify/assert"
)

type WhenNormalizingSignatureTestingSuite struct {
	sut *services.EIP712SignatureHelper
}

func WhenNormalizingSignatureBeforeEach(t *testing.T) *WhenNormalizingSignatureTestingSuite {
	blockchainSettings := &settings.BlockchainSettings{ChainId: 1}
	agentCollectionSettings := &settings.AgentCollectionSettings{SigningDomain: "AgentCollection"}

	sut := services.NewEIP712SignatureHelper(blockchainSettings, agentCollectionSettings)

	return &WhenNormalizingSignatureTestingSuite{
		sut: sut,
	}
}

func TestWhenNormalizingSignature(t *testing.T) {
	t.Parallel()

	t.Run("Given a valid 65-byte signature with v < 27", func(t *testing.T) {
		t.Parallel()

		t.Run("Should add 27 to the v value", func(t *testing.T) {
			t.Parallel()
			basics := WhenNormalizingSignatureBeforeEach(t)

			signature := make([]byte, 65)
			for i := 0; i < 64; i++ {
				signature[i] = byte(i)
			}
			signature[64] = 0

			result, err := basics.sut.NormalizeSignature(signature)

			assert.Nil(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, byte(27), result[64])
		})

		t.Run("Should add 27 to v=1", func(t *testing.T) {
			t.Parallel()
			basics := WhenNormalizingSignatureBeforeEach(t)

			signature := make([]byte, 65)
			signature[64] = 1

			result, err := basics.sut.NormalizeSignature(signature)

			assert.Nil(t, err)
			assert.Equal(t, byte(28), result[64])
		})
	})

	t.Run("Given a valid 65-byte signature with v >= 27", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not modify the v value", func(t *testing.T) {
			t.Parallel()
			basics := WhenNormalizingSignatureBeforeEach(t)

			signature := make([]byte, 65)
			signature[64] = 27

			result, err := basics.sut.NormalizeSignature(signature)

			assert.Nil(t, err)
			assert.Equal(t, byte(27), result[64])
		})

		t.Run("Should not modify v=28", func(t *testing.T) {
			t.Parallel()
			basics := WhenNormalizingSignatureBeforeEach(t)

			signature := make([]byte, 65)
			signature[64] = 28

			result, err := basics.sut.NormalizeSignature(signature)

			assert.Nil(t, err)
			assert.Equal(t, byte(28), result[64])
		})
	})

	t.Run("Given an invalid signature length", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error for too short signature", func(t *testing.T) {
			t.Parallel()
			basics := WhenNormalizingSignatureBeforeEach(t)

			signature := make([]byte, 64)

			result, err := basics.sut.NormalizeSignature(signature)

			assert.Nil(t, result)
			assert.NotNil(t, err)
		})

		t.Run("Should return an error for too long signature", func(t *testing.T) {
			t.Parallel()
			basics := WhenNormalizingSignatureBeforeEach(t)

			signature := make([]byte, 66)

			result, err := basics.sut.NormalizeSignature(signature)

			assert.Nil(t, result)
			assert.NotNil(t, err)
		})

		t.Run("Should return an error for empty signature", func(t *testing.T) {
			t.Parallel()
			basics := WhenNormalizingSignatureBeforeEach(t)

			signature := make([]byte, 0)

			result, err := basics.sut.NormalizeSignature(signature)

			assert.Nil(t, result)
			assert.NotNil(t, err)
		})
	})
}
