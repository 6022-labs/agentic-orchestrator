package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
)

type WhenGettingSlackUserByPairedHumanAgentTestingSuite struct {
	dbName   string
	dbClient *slack_bridge_db.DatabaseClient

	sut *repositories.SlackUserRepository
}

func WhenGettingSlackUserByPairedHumanAgentBeforeEach(
	t *testing.T,
) *WhenGettingSlackUserByPairedHumanAgentTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)

	sut := repositories.NewSlackUserRepository(dbClient)

	return &WhenGettingSlackUserByPairedHumanAgentTestingSuite{
		sut: sut,

		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenGettingSlackUserByPairedHumanAgentAfterEach(suite *WhenGettingSlackUserByPairedHumanAgentTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingSlackUserByPairedHumanAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given slack user does not exist for human agent", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackUserByPairedHumanAgentBeforeEach(t)
			defer WhenGettingSlackUserByPairedHumanAgentAfterEach(suite)

			collectionAddress := common_tests.RandomEthAddress().Hex()
			collectionTokenId := uint64(1)

			pairedSlackUserToHumanAgent, err := suite.sut.GetByPairedHumanAgent(
				context.Background(),
				collectionAddress,
				collectionTokenId,
			)

			assert.NoError(t, err)
			assert.Nil(t, pairedSlackUserToHumanAgent)
		})
	})

	t.Run("Given slack user exists for this collection address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackUserByPairedHumanAgentBeforeEach(t)
			defer WhenGettingSlackUserByPairedHumanAgentAfterEach(suite)

			collectionAddress := common_tests.RandomEthAddress().Hex()
			collectionTokenId := uint64(1)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(collectionAddress).
				WithCollectionTokenId(collectionTokenId).
				Build()

			err := suite.dbClient.Create(slackAgent).Error
			if err != nil {
				t.Fatal(err)
			}

			slackUser := model_builders.NewSlackUserBuilder().
				WithPairedHumanAgentCollectionAddress(&collectionAddress).
				WithPairedHumanAgentCollectionTokenId(&collectionTokenId).
				Build()

			_, err = suite.sut.Save(context.Background(), slackUser)
			if err != nil {
				t.Fatal(err)
			}

			slackUser, err = suite.sut.GetByPairedHumanAgent(context.Background(),
				collectionAddress,
				9999,
			)

			assert.NoError(t, err)
			assert.Nil(t, slackUser)
		})
	})

	t.Run("Given slack user exists for this collection token id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackUserByPairedHumanAgentBeforeEach(t)
			defer WhenGettingSlackUserByPairedHumanAgentAfterEach(suite)

			collectionAddress := common_tests.RandomEthAddress().Hex()
			collectionTokenId := uint64(1)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(collectionAddress).
				WithCollectionTokenId(collectionTokenId).
				Build()

			err := suite.dbClient.Create(slackAgent).Error
			if err != nil {
				t.Fatal(err)
			}

			slackUser := model_builders.NewSlackUserBuilder().
				WithPairedHumanAgentCollectionAddress(&collectionAddress).
				WithPairedHumanAgentCollectionTokenId(&collectionTokenId).
				Build()

			_, err = suite.sut.Save(context.Background(), slackUser)
			if err != nil {
				t.Fatal(err)
			}

			slackUser, err = suite.sut.GetByPairedHumanAgent(context.Background(),
				common_tests.RandomEthAddress().Hex(),
				collectionTokenId,
			)

			assert.NoError(t, err)
			assert.Nil(t, slackUser)
		})
	})

	t.Run("Given slack user exist for human agent", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return slack user", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackUserByPairedHumanAgentBeforeEach(t)
			defer WhenGettingSlackUserByPairedHumanAgentAfterEach(suite)

			collectionAddress := common_tests.RandomEthAddress().Hex()
			collectionTokenId := uint64(1)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(collectionAddress).
				WithCollectionTokenId(collectionTokenId).
				Build()

			err := suite.dbClient.Create(slackAgent).Error
			if err != nil {
				t.Fatal(err)
			}

			slackUser := model_builders.NewSlackUserBuilder().
				WithPairedHumanAgentCollectionAddress(&collectionAddress).
				WithPairedHumanAgentCollectionTokenId(&collectionTokenId).
				Build()

			_, err = suite.sut.Save(context.Background(), slackUser)
			if err != nil {
				t.Fatal(err)
			}

			result, err := suite.sut.GetByPairedHumanAgent(context.Background(),
				collectionAddress,
				collectionTokenId,
			)

			assert.Nil(t, err)
			assert.Equal(t, slackUser, result)
		})
	})
}
