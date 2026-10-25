package services_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

type WhenGettingMintedEventsTestingSuite struct {
	fixture *chainFixture
	sut     *services.AgentCollectionRequester
}

func WhenGettingMintedEventsBeforeEach(t *testing.T) *WhenGettingMintedEventsTestingSuite {
	fixture := newChainFixture(t)
	return &WhenGettingMintedEventsTestingSuite{
		fixture: fixture,
		sut:     fixture.agentCollectionRequester(),
	}
}

func TestWhenGettingMintedEvents(t *testing.T) {
	t.Parallel()

	t.Run("Given no collection address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nothing", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingMintedEventsBeforeEach(t)
			suite.fixture.mintAgent("unscanned-agent")

			events, err := suite.sut.GetAgentCollectionV1MintedEvent(context.Background(),
				0,
				suite.fixture.currentBlock(),
				nil,
			)

			assert.NoError(t, err)
			assert.Empty(t, events)
		})
	})

	t.Run("Given agents minted inside the scanned block range", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return one event per mint", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingMintedEventsBeforeEach(t)
			from := suite.fixture.currentBlock() + 1
			suite.fixture.mintAgent("first-minted-agent")
			suite.fixture.mintAgent("second-minted-agent")

			events, err := suite.sut.GetAgentCollectionV1MintedEvent(context.Background(),
				from,
				suite.fixture.currentBlock(),
				[]common.Address{suite.fixture.stack.Collection},
			)

			assert.NoError(t, err)
			if assert.Len(t, events, 2) {
				assert.Equal(t, suite.fixture.anvilAccountSettings.Address, events[0].To)
				assert.Equal(t, uint64(1), events[0].TokenId.Uint64())
				assert.Equal(t, uint64(2), events[1].TokenId.Uint64())
				assert.Positive(t, events[0].Raw.BlockNumber)
			}
		})

		t.Run("Should return nothing for a range that predates the mints", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingMintedEventsBeforeEach(t)
			beforeMint := suite.fixture.currentBlock()
			suite.fixture.mintAgent("late-minted-agent")

			events, err := suite.sut.GetAgentCollectionV1MintedEvent(context.Background(),
				0,
				beforeMint,
				[]common.Address{suite.fixture.stack.Collection},
			)

			assert.NoError(t, err)
			assert.Empty(t, events)
		})
	})

	t.Run("Given agents minted in two collections", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the mints of both", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingMintedEventsBeforeEach(t)
			secondCollection := suite.fixture.createCollection("second-collection")
			from := suite.fixture.currentBlock() + 1
			suite.fixture.mintAgent("first-collection-agent")
			suite.fixture.mintAgentIn(secondCollection, "second-collection-agent")

			events, err := suite.sut.GetAgentCollectionV1MintedEvent(context.Background(),
				from,
				suite.fixture.currentBlock(),
				[]common.Address{suite.fixture.stack.Collection, secondCollection},
			)

			assert.NoError(t, err)
			if assert.Len(t, events, 2) {
				assert.Equal(t, suite.fixture.stack.Collection, events[0].Raw.Address)
				assert.Equal(t, secondCollection, events[1].Raw.Address)
			}
		})
	})

	t.Run("Given more collection addresses than one log query accepts", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the mints of the collection in the last batch", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingMintedEventsBeforeEach(t)
			from := suite.fixture.currentBlock() + 1
			suite.fixture.mintAgent("last-batch-agent")
			addresses := make([]common.Address, 0, 601)
			for range 600 {
				addresses = append(addresses, common_tests.RandomEthAddress())
			}
			addresses = append(addresses, suite.fixture.stack.Collection)

			events, err := suite.sut.GetAgentCollectionV1MintedEvent(context.Background(),
				from,
				suite.fixture.currentBlock(),
				addresses,
			)

			assert.NoError(t, err)
			if assert.Len(t, events, 1) {
				assert.Equal(t, suite.fixture.stack.Collection, events[0].Raw.Address)
			}
		})
	})

	t.Run("Given a collection address that holds no contract", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nothing rather than fail", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingMintedEventsBeforeEach(t)
			suite.fixture.mintAgent("unrelated-agent")

			events, err := suite.sut.GetAgentCollectionV1MintedEvent(context.Background(),
				0,
				suite.fixture.currentBlock(),
				[]common.Address{common_tests.RandomEthAddress()},
			)

			assert.NoError(t, err)
			assert.Empty(t, events)
		})
	})
}
