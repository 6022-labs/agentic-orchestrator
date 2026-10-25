package services_test

import (
	"context"
	"net/http"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/services"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
)

type WhenGettingCurrentBlockNumberTestingSuite struct {
	fixture *chainFixture
	sut     *services.BlockchainRequester
}

func WhenGettingCurrentBlockNumberBeforeEach(t *testing.T) *WhenGettingCurrentBlockNumberTestingSuite {
	fixture := newChainFixture(t)
	return &WhenGettingCurrentBlockNumberTestingSuite{
		fixture: fixture,
		sut:     fixture.blockchainRequester(),
	}
}

func TestWhenGettingCurrentBlockNumber(t *testing.T) {
	t.Parallel()

	t.Run("Given the node rejects every request", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fail after a single request", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingCurrentBlockNumberBeforeEach(t)
			client, requests := suite.fixture.flakyHttpClient(
				[]int32{1, 2, 3, 4, 5, 6, 7, 8},
				suite.fixture.replyStatus(http.StatusBadRequest),
			)

			_, err := suite.fixture.blockchainRequesterThrough(client).GetCurrentBlockNumber(context.Background())

			var httpErr rpc.HTTPError
			if assert.ErrorAs(t, err, &httpErr) {
				assert.Equal(t, http.StatusBadRequest, httpErr.StatusCode)
			}
			assert.Equal(t, int32(1), requests.Load())
		})
	})

	t.Run("Given the node rate-limits the first request", func(t *testing.T) {
		t.Parallel()

		t.Run("Should retry and return the head block", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingCurrentBlockNumberBeforeEach(t)
			client, requests := suite.fixture.flakyHttpClient(
				[]int32{1},
				suite.fixture.replyStatus(http.StatusTooManyRequests),
			)

			blockNumber, err := suite.fixture.blockchainRequesterThrough(client).
				GetCurrentBlockNumber(context.Background())

			assert.NoError(t, err)
			assert.NotNil(t, blockNumber)
			assert.Equal(t, int32(2), requests.Load())
		})
	})

	t.Run("Given the node answers the first request with a JSON-RPC limit exceeded error", func(t *testing.T) {
		t.Parallel()

		t.Run("Should retry and return the head block", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingCurrentBlockNumberBeforeEach(t)
			client, requests := suite.fixture.flakyHttpClient([]int32{1}, suite.fixture.replyJsonRpcError(-32005))

			blockNumber, err := suite.fixture.blockchainRequesterThrough(client).
				GetCurrentBlockNumber(context.Background())

			assert.NoError(t, err)
			assert.NotNil(t, blockNumber)
			assert.Equal(t, int32(2), requests.Load())
		})
	})

	t.Run("Given a chain that has mined the stack deployment", func(t *testing.T) {
		t.Parallel()

		t.Run("Should report the head block", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingCurrentBlockNumberBeforeEach(t)
			// Parallel tests share the chain, so the head only ever moves forward.
			minimumBlock := suite.fixture.currentBlock()

			blockNumber, err := suite.sut.GetCurrentBlockNumber(context.Background())

			assert.NoError(t, err)
			if assert.NotNil(t, blockNumber) {
				assert.GreaterOrEqual(t, *blockNumber, minimumBlock)
			}
		})

		t.Run("Should advance after a new transaction is mined", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingCurrentBlockNumberBeforeEach(t)

			before, err := suite.sut.GetCurrentBlockNumber(context.Background())
			assert.NoError(t, err)

			suite.fixture.mintAgent("block-number-agent")

			after, err := suite.sut.GetCurrentBlockNumber(context.Background())
			assert.NoError(t, err)
			if assert.NotNil(t, before) && assert.NotNil(t, after) {
				assert.Greater(t, *after, *before)
			}
		})
	})
}
