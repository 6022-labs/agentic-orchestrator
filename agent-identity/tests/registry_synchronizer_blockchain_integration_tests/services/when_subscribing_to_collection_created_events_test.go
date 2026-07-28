package services_test

import (
	"context"
	"testing"
	"time"

	registry_abi "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_blockchain_integration_tests/fixtures"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

const (
	// collectionCreatedSubscriptionTimeout bounds the wait for a live event so a broken
	// subscription fails the test instead of hanging it.
	collectionCreatedSubscriptionTimeout = 30 * time.Second
	// collectionCreatedSubscriptionSettleWindow is how long the test waits to be convinced no event is coming.
	collectionCreatedSubscriptionSettleWindow = 2 * time.Second
)

type WhenSubscribingToCollectionCreatedEventsTestSuite struct {
	fixture *chainFixture
	sut     *services.AgentCollectionsManagerEventSubscriptionProvider
}

func WhenSubscribingToCollectionCreatedEventsBeforeEach(t *testing.T) *WhenSubscribingToCollectionCreatedEventsTestSuite {
	fixture := newChainFixture(t)
	return &WhenSubscribingToCollectionCreatedEventsTestSuite{
		fixture: fixture,
		sut:     fixture.agentCollectionsManagerEventSubscriptionProvider(),
	}
}

func (suite *WhenSubscribingToCollectionCreatedEventsTestSuite) createCollection(name string) {
	if _, err := fixtures.CreateCollection(
		context.Background(),
		suite.fixture.httpClient,
		suite.fixture.anvilAccountSettings,
		suite.fixture.stack.Manager,
		suite.fixture.stack.Creator,
		name,
	); err != nil {
		suite.fixture.t.Fatalf("failed to create the collection: %v", err)
	}
}

func TestWhenSubscribingToCollectionCreatedEvents(t *testing.T) {
	t.Parallel()

	t.Run("Given a live subscription on the manager", func(t *testing.T) {
		t.Parallel()

		t.Run("Should deliver collections created after subscribing", func(t *testing.T) {
			t.Parallel()
			suite := WhenSubscribingToCollectionCreatedEventsBeforeEach(t)

			events, subscription, err := suite.sut.StartCollectionCreatedSubscription(
				[]common.Address{suite.fixture.stack.Manager},
			)
			assert.NoError(t, err)
			if !assert.NotNil(t, subscription) {
				return
			}
			defer subscription.Unsubscribe()

			suite.createCollection("livecollection")

			var received *registry_abi.AgentCollectionsManagerCollectionCreated
			select {
			case received = <-events:
			case err := <-subscription.Err():
				t.Fatalf("subscription failed: %v", err)
			case <-time.After(collectionCreatedSubscriptionTimeout):
				t.Fatal("timed out waiting for the CollectionCreated event")
			}

			if assert.NotNil(t, received) {
				assert.Equal(t, "livecollection", received.Name)
				assert.Equal(t, uint64(2), received.CollectionId.Uint64())
				assert.NotEqual(t, common.Address{}, received.CollectionAddress)
			}
		})

		t.Run("Should not replay the collection created before subscribing", func(t *testing.T) {
			t.Parallel()
			suite := WhenSubscribingToCollectionCreatedEventsBeforeEach(t)

			events, subscription, err := suite.sut.StartCollectionCreatedSubscription(
				[]common.Address{suite.fixture.stack.Manager},
			)
			assert.NoError(t, err)
			if !assert.NotNil(t, subscription) {
				return
			}
			defer subscription.Unsubscribe()

			select {
			case event := <-events:
				assert.Nil(t, event, "subscriptions start at the head, they do not backfill")
			case <-time.After(collectionCreatedSubscriptionSettleWindow):
			}
		})
	})
}
