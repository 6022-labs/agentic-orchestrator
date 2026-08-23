package repositories

import (
	"context"

	"crypto/rand"
	"math/big"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	db_metrics_mocks "github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_db_mocks/metrics_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenGettingAgentByCollectionAddressAndCollectionTokenIdTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentRepository
}

func WhenGettingAgentByCollectionAddressAndCollectionTokenIdBeforeEach(
	t *testing.T,
) *WhenGettingAgentByCollectionAddressAndCollectionTokenIdTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbQueryMetrics := db_metrics_mocks.NewMockDbQueryMetricsInterface(gomock.NewController(t))
	dbQueryMetrics.EXPECT().
		RecordQuery(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes()

	dbClient, _ := registry_common_db.NewDatabaseClient(
		connectionString,
		zap.NewNop(),
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
		dbQueryMetrics,
	)
	if err = dbClient.Migrate(); err != nil {
		panic(err)
	}

	sut := repositories.NewAgentRepository(dbClient)

	return &WhenGettingAgentByCollectionAddressAndCollectionTokenIdTestingSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenGettingAgentByCollectionAddressAndCollectionTokenIdAfterEach(
	suite *WhenGettingAgentByCollectionAddressAndCollectionTokenIdTestingSuite,
) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingAgentByCollectionAddressAndCollectionTokenId(t *testing.T) {
	t.Parallel()

	t.Run("Given entity exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return entity", func(t *testing.T) {
			suite := WhenGettingAgentByCollectionAddressAndCollectionTokenIdBeforeEach(t)
			defer WhenGettingAgentByCollectionAddressAndCollectionTokenIdAfterEach(suite)

			existingEntity := model_builders.NewAgentBuilder().
				WithAttributes(map[string]string{
					"mission": "to the moon",
				}).
				Build()
			createdEntity, err := suite.sut.Save(context.Background(), existingEntity)
			if err != nil {
				t.Errorf("Error creating entity: %v", err)
			}

			fetchedEntity, err := suite.sut.GetByCollectionAddressAndCollectionTokenId(context.Background(),
				createdEntity.CollectionAddress,
				*createdEntity.CollectionTokenId,
			)

			assert.Nil(t, err)
			assert.Equal(t, createdEntity.ID, fetchedEntity.ID)
			assert.Equal(t, createdEntity.Name, fetchedEntity.Name)
			assert.Equal(t, createdEntity.Images, fetchedEntity.Images)
			assert.Equal(t, createdEntity.Creator, fetchedEntity.Creator)
			assert.Equal(t, createdEntity.EnsDomain, fetchedEntity.EnsDomain)
			assert.Equal(t, createdEntity.Attributes, fetchedEntity.Attributes)
		})
	})

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil without error", func(t *testing.T) {
			suite := WhenGettingAgentByCollectionAddressAndCollectionTokenIdBeforeEach(t)
			defer WhenGettingAgentByCollectionAddressAndCollectionTokenIdAfterEach(suite)

			max := new(big.Int).Lsh(big.NewInt(1), 64)
			nonExistantCollectionTokenIdBigInt, err := rand.Int(rand.Reader, max)
			if err != nil {
				t.Errorf("Error generating random uint64: %v", err)
			}

			nonExistantCollectionTokenId := nonExistantCollectionTokenIdBigInt.Uint64()
			nonExistantCollectionAddress := common_tests.RandomEthAddress()

			fetchedEntity, err := suite.sut.GetByCollectionAddressAndCollectionTokenId(context.Background(),
				nonExistantCollectionAddress.Hex(),
				nonExistantCollectionTokenId,
			)

			assert.Nil(t, err)
			assert.Nil(t, fetchedEntity)
		})
	})
}
