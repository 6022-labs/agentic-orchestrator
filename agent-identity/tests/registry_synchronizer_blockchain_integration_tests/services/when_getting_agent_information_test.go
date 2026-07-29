package services_test

import (
	"math/big"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/services"
	"github.com/stretchr/testify/assert"
)

type WhenGettingAgentInformationTestingSuite struct {
	fixture *chainFixture
	sut     *services.AgentCollectionRequester
}

func WhenGettingAgentInformationBeforeEach(t *testing.T) *WhenGettingAgentInformationTestingSuite {
	fixture := newChainFixture(t)
	return &WhenGettingAgentInformationTestingSuite{
		fixture: fixture,
		sut:     fixture.agentCollectionRequester(),
	}
}

func TestWhenGettingAgentInformation(t *testing.T) {
	t.Parallel()

	t.Run("Given a minted agent", func(t *testing.T) {
		t.Parallel()

		t.Run("Should decode its name, creator, images and attributes", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingAgentInformationBeforeEach(t)
			suite.fixture.mintAgent("described-agent")

			information, err := suite.sut.GetAgentInformation(suite.fixture.stack.Collection, *big.NewInt(1))

			assert.NoError(t, err)
			if assert.NotNil(t, information) {
				assert.Equal(t, "described-agent", information.Name)
				assert.Equal(t, suite.fixture.anvilAccountSettings.Address, information.Creator)
				assert.False(t, information.CloneOf.IsSet)
				if assert.Len(t, information.Images, 1) {
					assert.Equal(t, "default", information.Images[0].Key)
				}
				if assert.Len(t, information.Attributes, 2) {
					assert.Equal(t, "role", information.Attributes[0].Key)
					assert.Equal(t, "facilitator", information.Attributes[0].Value)
				}
				if assert.Len(t, information.Addresses, 1) {
					assert.Equal(t, "evm", information.Addresses[0].AddressType)
				}
			}
		})
	})

	t.Run("Given a token id that was never minted", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fail rather than return an empty struct", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingAgentInformationBeforeEach(t)

			_, err := suite.sut.GetAgentInformation(suite.fixture.stack.Collection, *big.NewInt(999))

			assert.Error(t, err)
		})
	})
}
