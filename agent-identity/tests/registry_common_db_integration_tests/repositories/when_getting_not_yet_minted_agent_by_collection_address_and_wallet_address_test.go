package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingNotYetMintedAgentByCollectionAddressAndWalletAddressTestingSuite struct {
	dbName string

	sut *repositories.AgentRepository
}

func WhenGettingNotYetMintedAgentByCollectionAddressAndWalletAddressBeforeEach() *WhenGettingNotYetMintedAgentByCollectionAddressAndWalletAddressTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbClient, _ := registry_common_db.NewDatabaseClient(
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

	return &WhenGettingNotYetMintedAgentByCollectionAddressAndWalletAddressTestingSuite{
		sut:    sut,
		dbName: dbName,
	}
}

func WhenGettingNotYetMintedAgentByCollectionAddressAndWalletAddressAfterEach(suite *WhenGettingNotYetMintedAgentByCollectionAddressAndWalletAddressTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingNotYetMintedAgentByCollectionAddressAndWalletAddress(t *testing.T) {
	t.Parallel()

	t.Run("Given not yet minted agent exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return agent", func(t *testing.T) {
			suite := WhenGettingNotYetMintedAgentByCollectionAddressAndWalletAddressBeforeEach()
			defer WhenGettingNotYetMintedAgentByCollectionAddressAndWalletAddressAfterEach(suite)

			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(nil).
				Build()
			createdAgent, err := suite.sut.Save(agent)
			if err != nil {
				t.Fatalf("Error saving agent: %v", err)
			}

			fetchedAgent, err := suite.sut.GetNotYetMintedByCollectionAddressAndWalletAddress(createdAgent.CollectionAddress, createdAgent.WalletAddress)

			assert.Nil(t, err)
			assert.NotNil(t, fetchedAgent)
			assert.Equal(t, createdAgent.ID, fetchedAgent.ID)
			assert.Equal(t, createdAgent.Name, fetchedAgent.Name)
			assert.Equal(t, createdAgent.Images, fetchedAgent.Images)
			assert.Equal(t, createdAgent.Creator, fetchedAgent.Creator)
			assert.Equal(t, createdAgent.EnsDomain, fetchedAgent.EnsDomain)
			assert.Equal(t, createdAgent.IsInternal, fetchedAgent.IsInternal)
			assert.Equal(t, createdAgent.Attributes, fetchedAgent.Attributes)
			assert.Equal(t, createdAgent.WalletAddress, fetchedAgent.WalletAddress)
		})
	})

	t.Run("Given not yet minted agent does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			suite := WhenGettingNotYetMintedAgentByCollectionAddressAndWalletAddressBeforeEach()
			defer WhenGettingNotYetMintedAgentByCollectionAddressAndWalletAddressAfterEach(suite)

			fetchedAgent, err := suite.sut.GetNotYetMintedByCollectionAddressAndWalletAddress("nonexistent_collection_address", "nonexistent_wallet_address")

			assert.Nil(t, fetchedAgent)
			assert.NotNil(t, err)
		})
	})
}
