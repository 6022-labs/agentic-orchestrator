package repositories

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/db_clients"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_db"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_db/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_db/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingAgentByCollectionAddressAndCollectionTokenIdTestSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentRepository
}

func WhenGettingAgentByCollectionAddressAndCollectionTokenIdBeforeEach() *WhenGettingAgentByCollectionAddressAndCollectionTokenIdTestSuite {
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

	return &WhenGettingAgentByCollectionAddressAndCollectionTokenIdTestSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenGettingAgentByCollectionAddressAndCollectionTokenIdAfterEach(suite *WhenGettingAgentByCollectionAddressAndCollectionTokenIdTestSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingAgentByCollectionAddressAndCollectionTokenId(t *testing.T) {
	t.Parallel()

	t.Run("Given entity exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return entity", func(t *testing.T) {
			suite := WhenGettingAgentByCollectionAddressAndCollectionTokenIdBeforeEach()
			defer WhenGettingAgentByCollectionAddressAndCollectionTokenIdAfterEach(suite)

			existingEntity := model_builders.NewAgentBuilder().
				WithAttributes(map[string]string{
					"mission": "to the moon",
				}).
				Build()
			createdEntity, err := suite.sut.Save(existingEntity)
			if err != nil {
				t.Errorf("Error creating entity: %v", err)
			}

			fetchedEntity, err := suite.sut.GetByCollectionAddressAndCollectionTokenId(
				createdEntity.CollectionAddress,
				*createdEntity.CollectionTokenId,
			)

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

		t.Run("Should return error", func(t *testing.T) {
			suite := WhenGettingAgentByCollectionAddressAndCollectionTokenIdBeforeEach()
			defer WhenGettingAgentByCollectionAddressAndCollectionTokenIdAfterEach(suite)

			max := new(big.Int).Lsh(big.NewInt(1), 64)
			nonExistantCollectionTokenIdBigInt, err := rand.Int(rand.Reader, max)
			if err != nil {
				t.Errorf("Error generating random uint64: %v", err)
			}

			nonExistantCollectionTokenId := nonExistantCollectionTokenIdBigInt.Uint64()
			nonExistantCollectionAddress := common_tests.RandomEthAddress()

			fetchedEntity, err := suite.sut.GetByCollectionAddressAndCollectionTokenId(
				nonExistantCollectionAddress.Hex(),
				nonExistantCollectionTokenId,
			)

			assert.NotNil(t, err)
			assert.Nil(t, fetchedEntity)
		})
	})
}
