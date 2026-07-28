package services_test

import (
	"math/big"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/stretchr/testify/assert"
)

type WhenGettingAgentOwnerTestSuite struct {
	fixture *chainFixture
	sut     *services.AgentCollectionRequester
}

func WhenGettingAgentOwnerBeforeEach(t *testing.T) *WhenGettingAgentOwnerTestSuite {
	fixture := newChainFixture(t)
	return &WhenGettingAgentOwnerTestSuite{
		fixture: fixture,
		sut:     fixture.agentCollectionRequester(),
	}
}

func TestWhenGettingAgentOwner(t *testing.T) {
	t.Parallel()

	t.Run("Given a minted agent", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the wallet it was minted to", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingAgentOwnerBeforeEach(t)
			suite.fixture.mintAgent("owned-agent")

			owner, err := suite.sut.GetAgentOwner(suite.fixture.stack.Collection, *big.NewInt(1))

			assert.NoError(t, err)
			if assert.NotNil(t, owner) {
				assert.Equal(t, suite.fixture.anvilAccountSettings.Address, *owner)
			}
		})
	})

	t.Run("Given a token id that was never minted", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fail rather than report a zero owner", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingAgentOwnerBeforeEach(t)

			_, err := suite.sut.GetAgentOwner(suite.fixture.stack.Collection, *big.NewInt(999))

			assert.Error(t, err)
		})
	})

	t.Run("Given an address that holds no collection contract", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fail rather than decode empty return data", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingAgentOwnerBeforeEach(t)

			_, err := suite.sut.GetAgentOwner(common_tests.RandomEthAddress(), *big.NewInt(1))

			assert.Error(t, err)
		})
	})
}
