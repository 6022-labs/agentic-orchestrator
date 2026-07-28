package services_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_blockchain_integration_tests/fixtures"
	"github.com/stretchr/testify/assert"
)

type WhenGettingMintProposalTestSuite struct {
	fixture *chainFixture
	sut     *services.AgentCollectionRequester
}

func WhenGettingMintProposalBeforeEach(t *testing.T) *WhenGettingMintProposalTestSuite {
	fixture := newChainFixture(t)
	return &WhenGettingMintProposalTestSuite{
		fixture: fixture,
		sut:     fixture.agentCollectionRequester(),
	}
}

func (suite *WhenGettingMintProposalTestSuite) createProposal(agentName string) {
	if err := fixtures.CreateMintProposal(
		context.Background(),
		suite.fixture.httpClient,
		suite.fixture.anvilAccountSettings,
		suite.fixture.stack.Collection,
		agentName,
	); err != nil {
		suite.fixture.t.Fatalf("failed to create the mint proposal: %v", err)
	}
}

func TestWhenGettingMintProposal(t *testing.T) {
	t.Parallel()

	t.Run("Given proposals created on a moderated collection", func(t *testing.T) {
		t.Parallel()

		t.Run("Should batch-read them and return the one matching the id", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingMintProposalBeforeEach(t)
			suite.createProposal("first-proposed-agent")
			suite.createProposal("second-proposed-agent")

			// Proposal ids are 1-based even though the batch reads array slots from 0.
			proposal, err := suite.sut.GetMintProposal(
				suite.fixture.stack.Collection,
				big.NewInt(2),
				suite.fixture.currentBlock(),
			)

			assert.NoError(t, err)
			if assert.NotNil(t, proposal) {
				assert.Equal(t, uint64(2), proposal.Id.Uint64())
				assert.Equal(t, "second-proposed-agent", proposal.Name)
				assert.Equal(t, suite.fixture.anvilAccountSettings.Address, proposal.Owner)
			}
		})

		t.Run("Should read the state as of the requested block", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingMintProposalBeforeEach(t)
			suite.createProposal("historical-agent")
			blockWithOneProposal := suite.fixture.currentBlock()
			suite.createProposal("later-agent")

			atLatest, err := suite.sut.GetMintProposal(
				suite.fixture.stack.Collection,
				big.NewInt(2),
				suite.fixture.currentBlock(),
			)
			assert.NoError(t, err)
			assert.NotNil(t, atLatest)

			_, err = suite.sut.GetMintProposal(suite.fixture.stack.Collection, big.NewInt(2), blockWithOneProposal)
			assert.Error(t, err)
		})
	})

	t.Run("Given a proposal id that does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fail rather than return an empty proposal", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingMintProposalBeforeEach(t)
			suite.createProposal("only-proposed-agent")

			_, err := suite.sut.GetMintProposal(
				suite.fixture.stack.Collection,
				big.NewInt(42),
				suite.fixture.currentBlock(),
			)

			assert.Error(t, err)
		})
	})
}
