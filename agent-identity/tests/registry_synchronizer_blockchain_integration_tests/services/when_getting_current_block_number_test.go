package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/services"
	"github.com/stretchr/testify/assert"
)

type WhenGettingCurrentBlockNumberTestSuite struct {
	fixture *chainFixture
	sut     *services.BlockchainRequester
}

func WhenGettingCurrentBlockNumberBeforeEach(t *testing.T) *WhenGettingCurrentBlockNumberTestSuite {
	fixture := newChainFixture(t)
	return &WhenGettingCurrentBlockNumberTestSuite{
		fixture: fixture,
		sut:     fixture.blockchainRequester(),
	}
}

func TestWhenGettingCurrentBlockNumber(t *testing.T) {
	t.Parallel()

	t.Run("Given a chain that has mined the stack deployment", func(t *testing.T) {
		t.Parallel()

		t.Run("Should report the head block", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingCurrentBlockNumberBeforeEach(t)
			// Parallel tests share the chain, so the head only ever moves forward.
			minimumBlock := suite.fixture.currentBlock()

			blockNumber, err := suite.sut.GetCurrentBlockNumber()

			assert.NoError(t, err)
			if assert.NotNil(t, blockNumber) {
				assert.GreaterOrEqual(t, *blockNumber, minimumBlock)
			}
		})

		t.Run("Should advance after a new transaction is mined", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingCurrentBlockNumberBeforeEach(t)

			before, err := suite.sut.GetCurrentBlockNumber()
			assert.NoError(t, err)

			suite.fixture.mintAgent("block-number-agent")

			after, err := suite.sut.GetCurrentBlockNumber()
			assert.NoError(t, err)
			if assert.NotNil(t, before) && assert.NotNil(t, after) {
				assert.Greater(t, *after, *before)
			}
		})
	})
}
