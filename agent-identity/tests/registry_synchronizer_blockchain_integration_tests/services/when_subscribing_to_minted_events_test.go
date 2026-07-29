package services_test

import (
	"testing"
	"time"

	registry_abi "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

const (
	// mintedSubscriptionTimeout bounds the wait for a live event so a broken
	// subscription fails the test instead of hanging it.
	mintedSubscriptionTimeout = 30 * time.Second
	// mintedSubscriptionSettleWindow is how long the test waits to be convinced no event is coming.
	mintedSubscriptionSettleWindow = 2 * time.Second
)

type WhenSubscribingToMintedEventsTestingSuite struct {
	fixture *chainFixture
	sut     *services.AgentCollectionEventSubscriptionProvider
}

func WhenSubscribingToMintedEventsBeforeEach(t *testing.T) *WhenSubscribingToMintedEventsTestingSuite {
	fixture := newChainFixture(t)
	return &WhenSubscribingToMintedEventsTestingSuite{
		fixture: fixture,
		sut:     fixture.agentCollectionEventSubscriptionProvider(),
	}
}

func TestWhenSubscribingToMintedEvents(t *testing.T) {
	t.Parallel()

	t.Run("Given a live subscription on the deployed collection", func(t *testing.T) {
		t.Parallel()

		t.Run("Should deliver mints that happen after subscribing", func(t *testing.T) {
			t.Parallel()
			suite := WhenSubscribingToMintedEventsBeforeEach(t)

			events, subscription, err := suite.sut.StartMintedSubscription(
				[]common.Address{suite.fixture.stack.Collection},
			)
			assert.NoError(t, err)
			if !assert.NotNil(t, subscription) {
				return
			}
			defer subscription.Unsubscribe()

			suite.fixture.mintAgent("subscribed-agent")

			var received *registry_abi.AgentCollectionV1Minted
			select {
			case received = <-events:
			case err := <-subscription.Err():
				t.Fatalf("subscription failed: %v", err)
			case <-time.After(mintedSubscriptionTimeout):
				t.Fatal("timed out waiting for the Minted event")
			}

			if assert.NotNil(t, received) {
				assert.Equal(t, suite.fixture.anvilAccountSettings.Address, received.To)
				assert.Equal(t, uint64(1), received.TokenId.Uint64())
				assert.Positive(t, received.Raw.BlockNumber)
			}
		})

		t.Run("Should stop delivering once unsubscribed", func(t *testing.T) {
			t.Parallel()
			suite := WhenSubscribingToMintedEventsBeforeEach(t)

			events, subscription, err := suite.sut.StartMintedSubscription(
				[]common.Address{suite.fixture.stack.Collection},
			)
			assert.NoError(t, err)
			if !assert.NotNil(t, subscription) {
				return
			}

			subscription.Unsubscribe()
			suite.fixture.mintAgent("ignored-agent")

			select {
			case event, open := <-events:
				assert.False(t, open && event != nil, "no event should arrive after unsubscribing")
			case <-time.After(mintedSubscriptionSettleWindow):
			}
		})
	})

	t.Run("Given a subscription filtered on an unrelated address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not deliver mints from the deployed collection", func(t *testing.T) {
			t.Parallel()
			suite := WhenSubscribingToMintedEventsBeforeEach(t)

			events, subscription, err := suite.sut.StartMintedSubscription(
				[]common.Address{common_tests.RandomEthAddress()},
			)
			assert.NoError(t, err)
			if !assert.NotNil(t, subscription) {
				return
			}
			defer subscription.Unsubscribe()

			suite.fixture.mintAgent("other-collection-agent")

			select {
			case event := <-events:
				assert.Nil(t, event, "the filter should have excluded this collection")
			case <-time.After(mintedSubscriptionSettleWindow):
			}
		})
	})
}
