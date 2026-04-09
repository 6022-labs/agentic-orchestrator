package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

type WhenCreatingWalletTestingSuite struct {
	sut *services.WalletCreator
}

func WhenCreatingWalletBeforeEach(t *testing.T) *WhenCreatingWalletTestingSuite {
	sut := services.NewWalletCreator()

	return &WhenCreatingWalletTestingSuite{
		sut: sut,
	}
}

func TestWhenCreatingWallet(t *testing.T) {
	t.Parallel()

	t.Run("Given wallet creation is requested", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a valid Ethereum address and private key", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingWalletBeforeEach(t)

			address, privateKey, err := basics.sut.Create()

			assert.Nil(t, err)
			assert.NotNil(t, address)
			assert.NotNil(t, privateKey)
			assert.True(t, common.IsHexAddress(*address))
			assert.NotEmpty(t, *privateKey)
		})

		t.Run("Should return unique addresses on multiple calls", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingWalletBeforeEach(t)

			address1, _, err1 := basics.sut.Create()
			address2, _, err2 := basics.sut.Create()

			assert.Nil(t, err1)
			assert.Nil(t, err2)
			assert.NotEqual(t, *address1, *address2)
		})
	})
}
