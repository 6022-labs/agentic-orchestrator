package repositories_test

import (
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
	"gorm.io/gorm"
)

type WhenGettingAgentsByEnsDomainTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentRepository
}

func WhenGettingAgentsByEnsDomainBeforeEach() *WhenGettingAgentsByEnsDomainTestingSuite {
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

	return &WhenGettingAgentsByEnsDomainTestingSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenGettingAgentsByEnsDomainAfterEach(suite *WhenGettingAgentsByEnsDomainTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingAgentsByEnsDomain(t *testing.T) {
	t.Parallel()

	t.Run("Given entity exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentsByEnsDomainBeforeEach()
			defer WhenGettingAgentsByEnsDomainAfterEach(suite)

			existingEntity := model_builders.NewAgentBuilder().Build()
			createdEntity, err := suite.sut.Save(existingEntity)
			if err != nil {
				t.Errorf("Error creating entity: %v", err)
			}

			fetchedEntity, err := suite.sut.GetByEnsDomain(createdEntity.EnsDomain)

			assert.Nil(t, err)
			assert.NotNil(t, fetchedEntity)
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

	t.Run("Given no entity exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentsByEnsDomainBeforeEach()
			defer WhenGettingAgentsByEnsDomainAfterEach(suite)

			fetchedEntity, err := suite.sut.GetByEnsDomain("nonexistent.ens")

			assert.Equal(t, err, gorm.ErrRecordNotFound)
			assert.Nil(t, fetchedEntity)
		})
	})
}
