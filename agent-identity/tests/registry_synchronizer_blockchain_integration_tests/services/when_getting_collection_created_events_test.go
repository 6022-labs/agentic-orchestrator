package services_test

import (
	"context"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_blockchain_integration_tests/fixtures"
	"github.com/stretchr/testify/assert"
)

type WhenGettingCollectionCreatedEventsTestingSuite struct {
	fixture *chainFixture
	sut     *services.AgentCollectionsManagerRequester
}

func WhenGettingCollectionCreatedEventsBeforeEach(t *testing.T) *WhenGettingCollectionCreatedEventsTestingSuite {
	fixture := newChainFixture(t)
	return &WhenGettingCollectionCreatedEventsTestingSuite{
		fixture: fixture,
		sut:     fixture.agentCollectionsManagerRequester(),
	}
}

func (suite *WhenGettingCollectionCreatedEventsTestingSuite) createCollection(name string) {
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

func TestWhenGettingCollectionCreatedEvents(t *testing.T) {
	t.Parallel()

	t.Run("Given the collection created during stack deployment", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return it when the range covers the deployment", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingCollectionCreatedEventsBeforeEach(t)

			events, err := suite.sut.GetAgentCollectionsManagerCollectionCreatedEvent(0, suite.fixture.currentBlock())

			assert.NoError(t, err)
			if assert.Len(t, events, 1) {
				assert.Equal(t, suite.fixture.stack.Collection, events[0].CollectionAddress)
				assert.Equal(t, fixtures.StackCollectionName, events[0].Name)
				assert.Equal(t, uint64(1), events[0].CollectionId.Uint64())
			}
		})
	})

	t.Run("Given a second collection created after a checkpoint", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return only the collection inside the range", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingCollectionCreatedEventsBeforeEach(t)
			checkpoint := suite.fixture.currentBlock()
			suite.createCollection("others")

			events, err := suite.sut.GetAgentCollectionsManagerCollectionCreatedEvent(
				checkpoint+1,
				suite.fixture.currentBlock(),
			)

			assert.NoError(t, err)
			if assert.Len(t, events, 1) {
				assert.Equal(t, "others", events[0].Name)
				assert.Equal(t, uint64(2), events[0].CollectionId.Uint64())
			}
		})
	})
}
