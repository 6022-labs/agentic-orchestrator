package repositories_test

import (
	"context"

	"crypto/rand"
	"math/big"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
)

type WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdTestingSuite struct {
	dbName string

	sut *repositories.AgentRepository
}

func WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdBeforeEach(
	t *testing.T,
) *WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)

	sut := repositories.NewAgentRepository(dbClient)

	return &WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdTestingSuite{
		sut:    sut,
		dbName: dbName,
	}
}

func WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdAfterEach(
	suite *WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdTestingSuite,
) {
	container_builders.GetPostgresContainerBuilder().DropDatabase(suite.dbName)
}

func TestWhenGettingNotYetMintedAgentByCollectionAddressAndProposalId(t *testing.T) {
	t.Parallel()

	t.Run("Given agent exists and is not yet minted", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdBeforeEach(t)
			defer WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdAfterEach(suite)

			agent := model_builders.NewAgentBuilder().WithCollectionTokenId(nil).Build()
			createdAgent, err := suite.sut.Save(context.Background(), agent)
			if err != nil {
				t.Errorf("Error creating agent: %v", err)
			}

			fetchedAgent, err := suite.sut.GetNotYetMintedByCollectionAddressAndProposalId(context.Background(),
				createdAgent.CollectionAddress,
				*createdAgent.ProposalId,
			)

			assert.Nil(t, err)
			assert.NotNil(t, fetchedAgent)
			assert.Equal(t, createdAgent.ID, fetchedAgent.ID)
		})
	})

	t.Run("Given agent does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			suite := WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdBeforeEach(t)
			defer WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdAfterEach(suite)

			max := new(big.Int).Lsh(big.NewInt(1), 64)
			nonExistingProposalIdBigInt, err := rand.Int(rand.Reader, max)
			if err != nil {
				t.Errorf("Error generating random uint64: %v", err)
			}

			nonExistingProposalId := nonExistingProposalIdBigInt.Uint64()
			nonExistingCollectionAddress := common_tests.RandomEthAddress()

			fetchedAgent, err := suite.sut.GetNotYetMintedByCollectionAddressAndProposalId(context.Background(),
				nonExistingCollectionAddress.Hex(),
				nonExistingProposalId,
			)

			assert.Nil(t, err)
			assert.Nil(t, fetchedAgent)
		})
	})
}
