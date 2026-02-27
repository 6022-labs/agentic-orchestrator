package repositories_test

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_db"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_db/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdTestingSuite struct {
	dbName string

	sut *repositories.AgentRepository
}

func WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdBeforeEach() *WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbClient, _ := orchestrator_db.NewDatabaseClient(
		connectionString,
		zap.NewNop(),
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
	)
	if err = dbClient.Migrate(); err != nil {
		panic(err)
	}

	sut := repositories.NewAgentRepository(dbClient)

	return &WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdTestingSuite{
		sut:    sut,
		dbName: dbName,
	}
}

func WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdAfterEach(suite *WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingNotYetMintedAgentByCollectionAddressAndProposalId(t *testing.T) {
	t.Parallel()

	t.Run("Given agent exists and is not yet minted", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdBeforeEach()
			defer WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdAfterEach(suite)

			agent := model_builders.NewAgentBuilder().WithCollectionTokenId(nil).Build()
			createdAgent, err := suite.sut.Save(agent)
			if err != nil {
				t.Errorf("Error creating agent: %v", err)
			}

			fetchedAgent, err := suite.sut.GetNotYetMintedByCollectionAddressAndProposalId(
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
			suite := WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdBeforeEach()
			defer WhenGettingNotYetMintedAgentByCollectionAddressAndProposalIdAfterEach(suite)

			max := new(big.Int).Lsh(big.NewInt(1), 64)
			nonExistingProposalIdBigInt, err := rand.Int(rand.Reader, max)
			if err != nil {
				t.Errorf("Error generating random uint64: %v", err)
			}

			nonExistingProposalId := nonExistingProposalIdBigInt.Uint64()
			nonExistingCollectionAddress := common_tests.RandomEthAddress()

			fetchedAgent, err := suite.sut.GetNotYetMintedByCollectionAddressAndProposalId(
				nonExistingCollectionAddress.Hex(),
				nonExistingProposalId,
			)

			assert.NotNil(t, err)
			assert.Nil(t, fetchedAgent)
		})
	})
}
