package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/db_clients"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_db"
	db_repositories "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_db/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_db/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingAgentByWalletAddressTestSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      repositories.AgentRepositoryInterface
}

func WhenGettingAgentByWalletAddressBeforeEach() *WhenGettingAgentByWalletAddressTestSuite {
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

	sut := db_repositories.NewAgentRepository(dbClient)

	return &WhenGettingAgentByWalletAddressTestSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenGettingAgentByWalletAddressAfterEach(suite *WhenGettingAgentByWalletAddressTestSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingAgentByWalletAddress(t *testing.T) {
	t.Parallel()

	t.Run("Given entity exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return entity", func(t *testing.T) {
			suite := WhenGettingAgentByWalletAddressBeforeEach()
			defer WhenGettingAgentByWalletAddressAfterEach(suite)

			existingEntity := model_builders.NewAgentBuilder().Build()
			createdEntity, err := suite.sut.Save(existingEntity)
			if err != nil {
				t.Errorf("Error creating entity: %v", err)
			}

			fetchedEntity, err := suite.sut.GetByWalletAddress(createdEntity.WalletAddress)

			assert.Nil(t, err)
			assert.Equal(t, createdEntity.ID, fetchedEntity.ID)
			assert.Equal(t, createdEntity.Name, fetchedEntity.Name)
			assert.Equal(t, createdEntity.Images, fetchedEntity.Images)
			assert.Equal(t, createdEntity.Creator, fetchedEntity.Creator)
			assert.Equal(t, createdEntity.EnsDomain, fetchedEntity.EnsDomain)
			assert.Equal(t, createdEntity.IsInternal, fetchedEntity.IsInternal)
			assert.Equal(t, createdEntity.Attributes, fetchedEntity.Attributes)
			assert.Equal(t, createdEntity.WalletAddress, fetchedEntity.WalletAddress)
		})
	})

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			suite := WhenGettingAgentByWalletAddressBeforeEach()
			defer WhenGettingAgentByWalletAddressAfterEach(suite)

			fetchedEntity, err := suite.sut.GetByWalletAddress("non-existing-wallet-address")

			assert.Nil(t, fetchedEntity)
			assert.NotNil(t, err)
		})
	})
}
