package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/settings"
	"github.com/stretchr/testify/assert"
)

type WhenNormalizingEip712SignatureTestSuite struct {
	sut *services.EIP712SignatureHelper
}

func WhenNormalizingEip712SignatureBeforeEach() *WhenNormalizingEip712SignatureTestSuite {
	blockchainSettings := &settings.BlockchainSettings{}
	agentCollectionSettings := &settings.AgentCollectionSettings{}

	sut := services.NewEIP712SignatureHelper(
		blockchainSettings,
		agentCollectionSettings,
	)

	return &WhenNormalizingEip712SignatureTestSuite{
		sut: sut,
	}
}

func TestWhenNormalizingEip712Signature(t *testing.T) {
	t.Parallel()

	t.Run("Given invalid signature", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenNormalizingEip712SignatureBeforeEach()

			signature := "0x1b0b4b7f1c1f6a7f7e"

			normalizedSignature, err := suite.sut.NormalizeSignature([]byte(signature))

			assert.Error(t, err)
			assert.Nil(t, normalizedSignature)
		})
	})

	t.Run("Given valid signature", func(t *testing.T) {
		t.Parallel()

		t.Run("Should normalize signature", func(t *testing.T) {
			t.Parallel()

			suite := WhenNormalizingEip712SignatureBeforeEach()

			signature := "5389353371ac10565913ca82fd7b51b430994dd0235185084d9f469db4e98069"
			signatureBytes := []byte(signature)
			signatureBytes = append(signatureBytes, byte(0))

			normalizedSignature, err := suite.sut.NormalizeSignature(signatureBytes)

			assert.Nil(t, err)
			assert.NotNil(t, normalizedSignature)
			assert.Equal(t, byte(27), normalizedSignature[64])
		})
	})
}
