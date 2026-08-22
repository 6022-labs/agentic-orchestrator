package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

type WhenSavingEntityUsingTransactionTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentRepository
}

func WhenSavingEntityUsingTransactionBeforeEach() *WhenSavingEntityUsingTransactionTestingSuite {
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

	return &WhenSavingEntityUsingTransactionTestingSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenSavingEntityUsingTransactionAfterEach(suite *WhenSavingEntityUsingTransactionTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenSavingEntityUsingTransaction(t *testing.T) {
	t.Parallel()

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return created entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingEntityUsingTransactionBeforeEach()
			defer WhenSavingEntityUsingTransactionAfterEach(suite)

			databaseTransaction, err := suite.dbClient.BeginTransaction()
			if err != nil {
				t.Fatalf("Failed to begin transaction: %v", err)
			}

			entityToCreate := model_builders.NewAgentBuilder().Build()
			createdEntity, err := suite.sut.SaveUsingTransaction(
				context.Background(),
				entityToCreate,
				databaseTransaction,
			)

			assert.Nil(t, err)
			assert.Equal(t, entityToCreate, createdEntity)
		})

		t.Run("Should save entity in database", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingEntityUsingTransactionBeforeEach()
			defer WhenSavingEntityUsingTransactionAfterEach(suite)

			databaseTransaction, err := suite.dbClient.BeginTransaction()
			if err != nil {
				t.Fatalf("Failed to begin transaction: %v", err)
			}

			entityToCreate := model_builders.NewAgentBuilder().Build()
			entityToCreate, _ = suite.sut.SaveUsingTransaction(
				context.Background(),
				entityToCreate,
				databaseTransaction,
			)

			databaseTransaction.Commit()

			foundedEntity := new(models.Agent)
			suite.dbClient.GetDB().First(foundedEntity, entityToCreate.ID)

			assert.Equal(t, entityToCreate.ID, foundedEntity.ID)
		})
	})

	t.Run("Given entity exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return created entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingEntityUsingTransactionBeforeEach()
			defer WhenSavingEntityUsingTransactionAfterEach(suite)

			entityToCreate := model_builders.NewAgentBuilder().Build()
			createdEntity, _ := suite.sut.Save(context.Background(), entityToCreate)

			databaseTransaction, err := suite.dbClient.BeginTransaction()
			if err != nil {
				t.Fatalf("Failed to begin transaction: %v", err)
			}

			entityToUpdate := model_builders.NewAgentBuilder().WithID(createdEntity.ID).Build()
			updatedEntity, err := suite.sut.SaveUsingTransaction(
				context.Background(),
				entityToUpdate,
				databaseTransaction,
			)

			assert.Nil(t, err)
			assert.Equal(t, entityToUpdate, updatedEntity)
		})

		t.Run("Should update entity in database", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingEntityBeforeEach()
			defer WhenSavingEntityAfterEach(suite)

			entityToCreate := model_builders.NewAgentBuilder().Build()
			createdEntity, _ := suite.sut.Save(context.Background(), entityToCreate)

			databaseTransaction, err := suite.dbClient.BeginTransaction()
			if err != nil {
				t.Fatalf("Failed to begin transaction: %v", err)
			}

			entityToUpdate := model_builders.NewAgentBuilder().WithID(createdEntity.ID).Build()
			suite.sut.SaveUsingTransaction(context.Background(), entityToUpdate, databaseTransaction)

			databaseTransaction.Commit()

			entityToCreate, _ = suite.sut.GetById(context.Background(), createdEntity.ID)

			assert.Equal(t, entityToCreate.ID, entityToUpdate.ID)
			assert.Equal(t, entityToCreate.Name, entityToUpdate.Name)
			assert.Equal(t, entityToCreate.Owner, entityToUpdate.Owner)
			assert.Equal(t, entityToCreate.CollectionTokenId, entityToUpdate.CollectionTokenId)

			assert.NotNil(t, entityToCreate.CreatedAt)
			assert.NotNil(t, entityToCreate.UpdatedAt)
		})
	})

	t.Run("Should return created entity", func(t *testing.T) {
		t.Parallel()

		suite := WhenSavingEntityUsingTransactionBeforeEach()
		defer WhenSavingEntityUsingTransactionAfterEach(suite)

		databaseTransaction, err := suite.dbClient.BeginTransaction()
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		entityToCreate := model_builders.NewAgentBuilder().Build()
		createdEntity, err := suite.sut.SaveUsingTransaction(context.Background(), entityToCreate, databaseTransaction)

		assert.Nil(t, err)
		assert.Equal(t, entityToCreate, createdEntity)
	})

	t.Run("Given transaction is roolbacked", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not save entity in database", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingEntityUsingTransactionBeforeEach()
			defer WhenSavingEntityUsingTransactionAfterEach(suite)

			databaseTransaction, err := suite.dbClient.BeginTransaction()
			if err != nil {
				t.Fatalf("Failed to begin transaction: %v", err)
			}

			entityToCreate := model_builders.NewAgentBuilder().Build()
			createdEntity, _ := suite.sut.SaveUsingTransaction(
				context.Background(),
				entityToCreate,
				databaseTransaction,
			)

			databaseTransaction.Rollback()

			foundedEntity := new(models.Agent)
			suite.dbClient.GetDB().First(foundedEntity, createdEntity.ID)

			assert.NotEqual(t, entityToCreate.ID, foundedEntity.ID)
		})
	})

	t.Run("Given transaction is committed", func(t *testing.T) {
		t.Parallel()

		t.Run("Should save entity in database", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingEntityUsingTransactionBeforeEach()
			defer WhenSavingEntityUsingTransactionAfterEach(suite)

			databaseTransaction, err := suite.dbClient.BeginTransaction()
			if err != nil {
				t.Fatalf("Failed to begin transaction: %v", err)
			}

			entityToCreate := model_builders.NewAgentBuilder().Build()
			entityToCreate, _ = suite.sut.SaveUsingTransaction(
				context.Background(),
				entityToCreate,
				databaseTransaction,
			)

			databaseTransaction.Commit()

			foundedEntity := new(models.Agent)
			suite.dbClient.GetDB().First(foundedEntity, entityToCreate.ID)

			assert.Equal(t, entityToCreate.ID, foundedEntity.ID)
			assert.Equal(t, entityToCreate.Name, foundedEntity.Name)
		})
	})
}
