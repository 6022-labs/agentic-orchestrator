package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
)

type WhenUpdatingSlackAgentOwnerWalletTestingSuite struct {
	dbName   string
	dbClient *slack_bridge_db.DatabaseClient

	sut *repositories.SlackAgentRepository
}

func WhenUpdatingSlackAgentOwnerWalletBeforeEach(t *testing.T) *WhenUpdatingSlackAgentOwnerWalletTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)

	sut := repositories.NewSlackAgentRepository(dbClient)

	return &WhenUpdatingSlackAgentOwnerWalletTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,

		sut: sut,
	}
}

func WhenUpdatingSlackAgentOwnerWalletAfterEach(suite *WhenUpdatingSlackAgentOwnerWalletTestingSuite) {
	container_builders.GetPostgresContainerBuilder().DropDatabase(suite.dbName)
}

func TestWhenUpdatingSlackAgentOwnerWallet(t *testing.T) {
	t.Parallel()

	t.Run("Given slack agent does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenUpdatingSlackAgentOwnerWalletBeforeEach(t)
			defer WhenUpdatingSlackAgentOwnerWalletAfterEach(suite)

			err := suite.sut.UpdateOwnerWallet(context.Background(), "0x123", 1, "0x456")

			assert.NotNil(t, err)
		})
	})

	t.Run("Given slack agent exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should update owner wallet", func(t *testing.T) {
			t.Parallel()

			suite := WhenUpdatingSlackAgentOwnerWalletBeforeEach(t)
			defer WhenUpdatingSlackAgentOwnerWalletAfterEach(suite)

			slackAgent := model_builders.NewSlackAgentBuilder().Build()

			_, err := suite.sut.Save(context.Background(), slackAgent)
			if err != nil {
				t.Fatal(err)
			}

			err = suite.sut.UpdateOwnerWallet(
				context.Background(),
				slackAgent.CollectionAddress,
				slackAgent.CollectionTokenId,
				"0x456",
			)
			if err != nil {
				t.Fatal(err)
			}

			updatedAgent, err := suite.sut.GetByCollectionAddressAndCollectionTokenId(context.Background(),
				slackAgent.CollectionAddress,
				slackAgent.CollectionTokenId,
			)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, "0x456", updatedAgent.OwnerWallet)
		})
	})
}
