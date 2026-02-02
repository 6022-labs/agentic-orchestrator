package repositories_test

import (
	"fmt"
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	db_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingProcessableMessagesForAgentAndThreadTestingSuite struct {
	dbName string

	dbClient *slack_bridge_db.DatabaseClient

	sut *repositories.SlackMessageRepository

	slackUserRepository    *repositories.SlackUserRepository
	slackAgentRepository   *repositories.SlackAgentRepository
	slackThreadRepository  *repositories.SlackThreadRepository
	slackChannelRepository *repositories.SlackChannelRepository
}

func WhenGettingProcessableMessagesForAgentAndThreadBeforeEach() *WhenGettingProcessableMessagesForAgentAndThreadTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbClient, _ := slack_bridge_db.NewDatabaseClient(
		connectionString,
		zap.NewNop(),
		&db_settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
	)
	err = dbClient.Migrate()
	if err != nil {
		panic(err)
	}

	sut := repositories.NewSlackMessageRepository(dbClient)

	slackUserRepository := repositories.NewSlackUserRepository(dbClient)
	slackAgentRepository := repositories.NewSlackAgentRepository(dbClient)
	slackThreadRepository := repositories.NewSlackThreadRepository(dbClient)
	slackChannelRepository := repositories.NewSlackChannelRepository(dbClient)

	return &WhenGettingProcessableMessagesForAgentAndThreadTestingSuite{
		sut: sut,

		dbName:   dbName,
		dbClient: dbClient,

		slackUserRepository:    slackUserRepository,
		slackAgentRepository:   slackAgentRepository,
		slackThreadRepository:  slackThreadRepository,
		slackChannelRepository: slackChannelRepository,
	}
}

func WhenGettingProcessableMessagesForAgentAndThreadAfterEach(suite *WhenGettingProcessableMessagesForAgentAndThreadTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingProcessableMessagesForAgentAndThread(t *testing.T) {
	t.Parallel()

	generateSlackThreadWithAgent := func(suite *WhenGettingProcessableMessagesForAgentAndThreadTestingSuite, agentRole enums.AgentRole, agentEnsDomain string) (*models.SlackChannel, *models.SlackThread, *models.SlackAgent) {
		slackChannel := model_builders.NewSlackChannelBuilder().Build()
		_, err := suite.slackChannelRepository.Save(slackChannel)
		if err != nil {
			t.Fatal(err)
		}

		slackThread := model_builders.NewSlackThreadBuilder().
			WithChannelId(slackChannel.Id).
			WithIsConcluded(false).
			Build()
		_, err = suite.slackThreadRepository.Save(slackThread)
		if err != nil {
			t.Fatal(err)
		}

		slackAgent := model_builders.NewSlackAgentBuilder().
			WithRole(agentRole).
			WithEnsDomain(agentEnsDomain).
			Build()
		_, err = suite.slackAgentRepository.Save(slackAgent)
		if err != nil {
			t.Fatal(err)
		}

		_, err = suite.slackChannelRepository.AddAgent(slackChannel, slackAgent)
		if err != nil {
			t.Fatal(err)
		}

		return slackChannel, slackThread, slackAgent
	}

	t.Run("Given no messages for the thread", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return empty array", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingProcessableMessagesForAgentAndThreadBeforeEach()
			defer WhenGettingProcessableMessagesForAgentAndThreadAfterEach(suite)

			agent := model_builders.NewSlackAgentBuilder().Build()

			messages, err := suite.sut.GetProcessableMessagesForAgentAndThread(agent, "thread-id", "channel-id")

			assert.NoError(t, err)
			assert.Empty(t, messages, "Expected no messages for the thread")
		})
	})

	t.Run("Given agent is facilitator", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return all messages of the thread", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingProcessableMessagesForAgentAndThreadBeforeEach()
			defer WhenGettingProcessableMessagesForAgentAndThreadAfterEach(suite)

			slackChannel, slackThread, slackAgent := generateSlackThreadWithAgent(suite, enums.AgentRoleFacilitator, "facilitator.ens")

			message1 := model_builders.NewSlackMessageBuilder().
				WithAuthorCollectionAddress(slackAgent.CollectionAddress).
				WithAuthorCollectionTokenId(slackAgent.CollectionTokenId).
				WithIsFromAgent(true).
				WithIsFromUser(false).
				WithThreadId(slackThread.Id).
				WithChannelId(slackChannel.Id).
				WithContent("First message in thread").
				Build()

			slackUser := model_builders.NewSlackUserBuilder().
				WithId("author-user-id").
				WithPairedHumanAgentCollectionAddress(&slackAgent.CollectionAddress).
				WithPairedHumanAgentCollectionTokenId(&slackAgent.CollectionTokenId).
				Build()
			_, err := suite.slackUserRepository.Save(slackUser)
			if err != nil {
				t.Fatal(err)
			}

			message2 := model_builders.NewSlackMessageBuilder().
				WithAuthorUserId(slackUser.Id).
				WithThreadId(slackThread.Id).
				WithChannelId(slackChannel.Id).
				WithContent(fmt.Sprintf("mentioned agent %s", slackAgent.EnsDomain)).
				Build()

			_, err = suite.sut.SaveMany([]*models.SlackMessage{message1, message2})
			assert.NoError(t, err)

			messages, err := suite.sut.GetProcessableMessagesForAgentAndThread(slackAgent, slackThread.Id, slackChannel.Id)
			assert.NoError(t, err)
			assert.Len(t, messages, 2, "Expected two messages in the thread")
		})
	})

	t.Run("Given agent is expert", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return all messages where agent is mentioned", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingProcessableMessagesForAgentAndThreadBeforeEach()
			defer WhenGettingProcessableMessagesForAgentAndThreadAfterEach(suite)

			slackChannel, slackThread, slackAgent := generateSlackThreadWithAgent(suite, enums.AgentRoleExpert, "agent.ens")

			message1 := model_builders.NewSlackMessageBuilder().
				WithAuthorCollectionAddress(slackAgent.CollectionAddress).
				WithAuthorCollectionTokenId(slackAgent.CollectionTokenId).
				WithThreadId(slackThread.Id).
				WithChannelId(slackChannel.Id).
				WithContent("First message in thread").
				Build()

			slackUser := model_builders.NewSlackUserBuilder().
				WithId("another-author").
				Build()
			_, err := suite.slackUserRepository.Save(slackUser)
			if err != nil {
				t.Fatal(err)
			}

			message2 := model_builders.NewSlackMessageBuilder().
				WithAuthorUserId(slackUser.Id).
				WithThreadId(slackThread.Id).
				WithChannelId(slackChannel.Id).
				WithContent(fmt.Sprintf("mentioned agent %s", slackAgent.EnsDomain)).
				Build()

			message3 := model_builders.NewSlackMessageBuilder().
				WithAuthorUserId(slackUser.Id).
				WithThreadId(slackThread.Id).
				WithChannelId(slackChannel.Id).
				WithContent("This message does not mention the agent").
				Build()

			_, err = suite.sut.SaveMany([]*models.SlackMessage{message1, message2, message3})
			assert.NoError(t, err)

			messages, err := suite.sut.GetProcessableMessagesForAgentAndThread(slackAgent, slackThread.Id, slackChannel.Id)

			assert.NoError(t, err)
			assert.Len(t, messages, 1, "Expected one messages where agent is mentioned")
		})
	})

	t.Run("Given agent is human", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return all messages from this human agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingProcessableMessagesForAgentAndThreadBeforeEach()
			defer WhenGettingProcessableMessagesForAgentAndThreadAfterEach(suite)

			slackChannel, slackThread, slackAgent := generateSlackThreadWithAgent(suite, enums.AgentRoleHuman, "human.ens")

			slackUser := model_builders.NewSlackUserBuilder().
				WithId("author-user-id").
				WithPairedHumanAgentCollectionAddress(&slackAgent.CollectionAddress).
				WithPairedHumanAgentCollectionTokenId(&slackAgent.CollectionTokenId).
				Build()
			_, err := suite.slackUserRepository.Save(slackUser)
			if err != nil {
				t.Fatal(err)
			}

			otherAgent := model_builders.NewSlackAgentBuilder().
				WithRole(enums.AgentRoleExpert).
				WithEnsDomain("other.ens").
				Build()

			_, err = suite.slackAgentRepository.Save(otherAgent)
			if err != nil {
				t.Fatal(err)
			}

			message1 := model_builders.NewSlackMessageBuilder().
				WithAuthorUserId(slackUser.Id).
				WithAuthorCollectionAddress(slackAgent.CollectionAddress).
				WithAuthorCollectionTokenId(slackAgent.CollectionTokenId).
				WithIsFromAgent(true).
				WithIsFromUser(true).
				WithThreadId(slackThread.Id).
				WithChannelId(slackChannel.Id).
				WithContent("First message in thread").
				Build()

			message2 := model_builders.NewSlackMessageBuilder().
				WithAuthorUserId(slackUser.Id).
				WithAuthorCollectionAddress(slackAgent.CollectionAddress).
				WithAuthorCollectionTokenId(slackAgent.CollectionTokenId).
				WithIsFromAgent(true).
				WithIsFromUser(true).
				WithThreadId(slackThread.Id).
				WithChannelId(slackChannel.Id).
				WithContent("Third message in thread from a user").
				Build()

			message3 := model_builders.NewSlackMessageBuilder().
				WithAuthorCollectionAddress(otherAgent.CollectionAddress).
				WithAuthorCollectionTokenId(otherAgent.CollectionTokenId).
				WithIsFromAgent(true).
				WithIsFromUser(false).
				WithThreadId(slackThread.Id).
				WithChannelId(slackChannel.Id).
				WithContent("Second message in thread from another agent").
				Build()

			_, err = suite.sut.SaveMany([]*models.SlackMessage{message1, message2, message3})
			assert.NoError(t, err)

			messages, err := suite.sut.GetProcessableMessagesForAgentAndThread(slackAgent, slackThread.Id, slackChannel.Id)

			assert.NoError(t, err)
			assert.Len(t, messages, 2, "Expected two messages from this human agent")
			if len(messages) == 2 {
				assert.True(t, lo.ContainsBy(messages, func(m *models.SlackMessage) bool {
					return m.Id == message1.Id
				}), "Expected first message to be in the result")
				assert.True(t, lo.ContainsBy(messages, func(m *models.SlackMessage) bool {
					return m.Id == message2.Id
				}), "Expected second message to be in the result")
			}
		})

		t.Run("Should return all messages where human agent is mentioned", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingProcessableMessagesForAgentAndThreadBeforeEach()
			defer WhenGettingProcessableMessagesForAgentAndThreadAfterEach(suite)

			slackChannel, slackThread, slackAgent := generateSlackThreadWithAgent(suite, enums.AgentRoleHuman, "human.ens")

			slackUser := model_builders.NewSlackUserBuilder().
				WithId("author-user-id").
				WithPairedHumanAgentCollectionAddress(&slackAgent.CollectionAddress).
				WithPairedHumanAgentCollectionTokenId(&slackAgent.CollectionTokenId).
				Build()
			_, err := suite.slackUserRepository.Save(slackUser)
			if err != nil {
				t.Fatal(err)
			}

			otherUser := model_builders.NewSlackUserBuilder().
				WithId("other-author").
				Build()

			_, err = suite.slackUserRepository.Save(otherUser)
			if err != nil {
				t.Fatal(err)
			}

			message1 := model_builders.NewSlackMessageBuilder().
				WithAuthorUserId(otherUser.Id).
				WithThreadId(slackThread.Id).
				WithChannelId(slackChannel.Id).
				WithContent(fmt.Sprintf("mentioned agent %s", slackAgent.EnsDomain)).
				Build()

			message2 := model_builders.NewSlackMessageBuilder().
				WithAuthorUserId(otherUser.Id).
				WithThreadId(slackThread.Id).
				WithChannelId(slackChannel.Id).
				WithContent("This message does not mention the agent").
				Build()

			_, err = suite.sut.SaveMany([]*models.SlackMessage{message1, message2})
			assert.NoError(t, err)

			messages, err := suite.sut.GetProcessableMessagesForAgentAndThread(slackAgent, slackThread.Id, slackChannel.Id)

			assert.NoError(t, err)
			assert.Len(t, messages, 1, "Expected one messages where human agent is mentioned")

			if len(messages) == 1 {
				assert.Equal(t, message1.Id, messages[0].Id, "Expected the correct message to be returned")
			}
		})
	})
}
