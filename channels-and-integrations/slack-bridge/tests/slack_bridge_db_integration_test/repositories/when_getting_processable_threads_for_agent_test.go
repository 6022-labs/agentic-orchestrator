package repositories_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	db_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingProcessableThreadsForAgentTestingSuite struct {
	dbName   string
	dbClient *slack_bridge_db.DatabaseClient

	slackUserRepository    *repositories.SlackUserRepository
	slackAgentRepository   *repositories.SlackAgentRepository
	slackMessageRepository *repositories.SlackMessageRepository
	slackChannelRepository *repositories.SlackChannelRepository

	sut *repositories.SlackThreadRepository
}

func WhenGettingProcessableThreadsForAgentBeforeEach() *WhenGettingProcessableThreadsForAgentTestingSuite {
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

	sut := repositories.NewSlackThreadRepository(dbClient)

	slackUserRepository := repositories.NewSlackUserRepository(dbClient)
	slackAgentRepository := repositories.NewSlackAgentRepository(dbClient)
	slackMessageRepository := repositories.NewSlackMessageRepository(dbClient)
	slackChannelRepository := repositories.NewSlackChannelRepository(dbClient)

	return &WhenGettingProcessableThreadsForAgentTestingSuite{
		sut: sut,

		dbName:   dbName,
		dbClient: dbClient,

		slackUserRepository:    slackUserRepository,
		slackAgentRepository:   slackAgentRepository,
		slackMessageRepository: slackMessageRepository,
		slackChannelRepository: slackChannelRepository,
	}
}

func WhenGettingProcessableThreadsForAgentAfterEach(suite *WhenGettingProcessableThreadsForAgentTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingProcessableThreadsForAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given agent has no associated channels", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return empty list", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingProcessableThreadsForAgentBeforeEach()
			defer WhenGettingProcessableThreadsForAgentAfterEach(suite)

			agent := model_builders.NewSlackAgentBuilder().Build()
			agent, err := suite.slackAgentRepository.Save(agent)
			assert.NoError(t, err)

			processableThreads, err := suite.sut.GetProcessableThreadsForAgent(agent, 30)

			assert.NoError(t, err)
			assert.Empty(t, processableThreads)
		})
	})

	t.Run("Given agent has no associated threads", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return empty list", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingProcessableThreadsForAgentBeforeEach()
			defer WhenGettingProcessableThreadsForAgentAfterEach(suite)

			channel := model_builders.NewSlackChannelBuilder().Build()
			agent := model_builders.NewSlackAgentBuilder().Build()

			_, err := suite.slackChannelRepository.Save(channel)
			assert.NoError(t, err)

			agent, err = suite.slackAgentRepository.Save(agent)
			assert.NoError(t, err)

			_, err = suite.slackChannelRepository.AddAgent(channel, agent)
			assert.NoError(t, err)

			processableThreads, err := suite.sut.GetProcessableThreadsForAgent(agent, 30)

			assert.NoError(t, err)
			assert.Empty(t, processableThreads)
		})
	})

	t.Run("Given agent has concluded associated threads", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return empty list", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingProcessableThreadsForAgentBeforeEach()
			defer WhenGettingProcessableThreadsForAgentAfterEach(suite)

			channel := model_builders.NewSlackChannelBuilder().Build()
			agent := model_builders.NewSlackAgentBuilder().Build()

			_, err := suite.slackChannelRepository.Save(channel)
			assert.NoError(t, err)

			agent, err = suite.slackAgentRepository.Save(agent)
			assert.NoError(t, err)

			_, err = suite.slackChannelRepository.AddAgent(channel, agent)
			assert.NoError(t, err)

			thread := model_builders.NewSlackThreadBuilder().
				WithChannelId(channel.Id).
				WithIsConcluded(true).
				Build()

			_, err = suite.sut.Save(thread)
			assert.NoError(t, err)

			processableThreads, err := suite.sut.GetProcessableThreadsForAgent(agent, 30)

			assert.NoError(t, err)
			assert.Empty(t, processableThreads)
		})
	})

	t.Run("Given agent is facilitator", func(t *testing.T) {
		t.Parallel()

		t.Run("And facilitator never send a message", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return processable threads", func(t *testing.T) {
				t.Parallel()

				suite := WhenGettingProcessableThreadsForAgentBeforeEach()
				defer WhenGettingProcessableThreadsForAgentAfterEach(suite)

				slackChannel := model_builders.NewSlackChannelBuilder().Build()
				slackAgent := model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()

				_, err := suite.slackChannelRepository.Save(slackChannel)
				assert.NoError(t, err)

				slackAgent, err = suite.slackAgentRepository.Save(slackAgent)
				assert.NoError(t, err)

				_, err = suite.slackChannelRepository.AddAgent(slackChannel, slackAgent)
				assert.NoError(t, err)

				slackThread := model_builders.NewSlackThreadBuilder().
					WithChannelId(slackChannel.Id).
					WithIsConcluded(false).
					Build()

				_, err = suite.sut.Save(slackThread)
				assert.NoError(t, err)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId("author-user-id").
					Build()

				_, err = suite.slackUserRepository.Save(slackUser)
				assert.NoError(t, err)

				message := model_builders.NewSlackMessageBuilder().
					WithId(model_builders.GenerateSlackMessageID(time.Now())).
					WithThreadId(slackThread.Id).
					WithChannelId(slackChannel.Id).
					WithAuthorUserId(slackUser.Id).
					Build()

				_, err = suite.slackMessageRepository.Save(message)
				assert.NoError(t, err)

				processableThreads, err := suite.sut.GetProcessableThreadsForAgent(slackAgent, 30)

				assert.NoError(t, err)
				assert.Len(t, processableThreads, 1)
				assert.False(t, processableThreads[0].IsConcluded)
				assert.Equal(t, slackThread.Id, processableThreads[0].Id)
				assert.Equal(t, slackThread.ChannelId, processableThreads[0].ChannelId)
			})
		})

		t.Run("And facilitator latest message is older than processing interval", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return processable threads", func(t *testing.T) {
				t.Parallel()

				suite := WhenGettingProcessableThreadsForAgentBeforeEach()
				defer WhenGettingProcessableThreadsForAgentAfterEach(suite)

				slackChannel := model_builders.NewSlackChannelBuilder().Build()
				slackAgent := model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()

				_, err := suite.slackChannelRepository.Save(slackChannel)
				assert.NoError(t, err)

				slackAgent, err = suite.slackAgentRepository.Save(slackAgent)
				assert.NoError(t, err)

				_, err = suite.slackChannelRepository.AddAgent(slackChannel, slackAgent)
				assert.NoError(t, err)

				slackThread := model_builders.NewSlackThreadBuilder().
					WithChannelId(slackChannel.Id).
					WithIsConcluded(false).
					Build()

				_, err = suite.sut.Save(slackThread)
				assert.NoError(t, err)

				baseInterval := uint64(30)
				interval := baseInterval + 1

				message := model_builders.NewSlackMessageBuilder().
					WithId(model_builders.GenerateSlackMessageID(time.Now().Add(-time.Duration(interval) * time.Second))).
					WithThreadId(slackThread.Id).
					WithChannelId(slackChannel.Id).
					WithAuthorCollectionAddress(slackAgent.CollectionAddress).
					WithAuthorCollectionTokenId(slackAgent.CollectionTokenId).
					Build()

				_, err = suite.slackMessageRepository.Save(message)
				assert.NoError(t, err)

				processableThreads, err := suite.sut.GetProcessableThreadsForAgent(slackAgent, baseInterval)

				assert.NoError(t, err)
				assert.Len(t, processableThreads, 1)
				assert.False(t, processableThreads[0].IsConcluded)
				assert.Equal(t, slackThread.Id, processableThreads[0].Id)
				assert.Equal(t, slackThread.ChannelId, processableThreads[0].ChannelId)
			})
		})

		t.Run("And messages are newer than threshold", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return empty list", func(t *testing.T) {
				t.Parallel()

				suite := WhenGettingProcessableThreadsForAgentBeforeEach()
				defer WhenGettingProcessableThreadsForAgentAfterEach(suite)

				slackChannel := model_builders.NewSlackChannelBuilder().Build()
				slackAgent := model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()

				_, err := suite.slackChannelRepository.Save(slackChannel)
				assert.NoError(t, err)

				slackAgent, err = suite.slackAgentRepository.Save(slackAgent)
				assert.NoError(t, err)

				_, err = suite.slackChannelRepository.AddAgent(slackChannel, slackAgent)
				assert.NoError(t, err)

				slackThread := model_builders.NewSlackThreadBuilder().
					WithChannelId(slackChannel.Id).
					WithIsConcluded(false).
					Build()

				_, err = suite.sut.Save(slackThread)
				assert.NoError(t, err)

				message := model_builders.NewSlackMessageBuilder().
					WithId(model_builders.GenerateSlackMessageID(time.Now())).
					WithThreadId(slackThread.Id).
					WithChannelId(slackChannel.Id).
					WithAuthorCollectionAddress(slackAgent.CollectionAddress).
					WithAuthorCollectionTokenId(slackAgent.CollectionTokenId).
					Build()
				_, err = suite.slackMessageRepository.Save(message)
				assert.NoError(t, err)

				processableThreads, err := suite.sut.GetProcessableThreadsForAgent(slackAgent, 30)

				assert.NoError(t, err)
				assert.Empty(t, processableThreads)
			})
		})
	})

	t.Run("Given agent is human", func(t *testing.T) {
		t.Parallel()

		t.Run("And agent is mentioned", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return processable threads", func(t *testing.T) {
				t.Parallel()

				suite := WhenGettingProcessableThreadsForAgentBeforeEach()
				defer WhenGettingProcessableThreadsForAgentAfterEach(suite)

				slackChannel := model_builders.NewSlackChannelBuilder().Build()
				slackAgent := model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleHuman).
					Build()

				_, err := suite.slackChannelRepository.Save(slackChannel)
				assert.NoError(t, err)

				slackAgent, err = suite.slackAgentRepository.Save(slackAgent)
				assert.NoError(t, err)

				_, err = suite.slackChannelRepository.AddAgent(slackChannel, slackAgent)
				assert.NoError(t, err)

				slackThread := model_builders.NewSlackThreadBuilder().
					WithId(model_builders.GenerateRandomSlackMessageID()).
					WithChannelId(slackChannel.Id).
					WithIsConcluded(false).
					Build()

				_, err = suite.sut.Save(slackThread)
				assert.NoError(t, err)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId("author-user-id").
					Build()

				_, err = suite.slackUserRepository.Save(slackUser)
				assert.NoError(t, err)

				// Message mentioning the human agent's ENS domain
				message := model_builders.NewSlackMessageBuilder().
					WithId(model_builders.GenerateRandomSlackMessageID()).
					WithThreadId(slackThread.Id).
					WithChannelId(slackChannel.Id).
					WithContent(fmt.Sprintf("Hey %s, can you help?", slackAgent.EnsDomain)).
					WithAuthorUserId(slackUser.Id).
					Build()

				_, err = suite.slackMessageRepository.Save(message)
				assert.NoError(t, err)

				processableThreads, err := suite.sut.GetProcessableThreadsForAgent(slackAgent, 30)

				assert.NoError(t, err)
				assert.Len(t, processableThreads, 1)
				assert.False(t, processableThreads[0].IsConcluded)
				assert.Equal(t, slackThread.Id, processableThreads[0].Id)
				assert.Equal(t, slackThread.ChannelId, processableThreads[0].ChannelId)
			})
		})

		t.Run("And underlying physical human send a message", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return processable threads", func(t *testing.T) {
				t.Parallel()

				suite := WhenGettingProcessableThreadsForAgentBeforeEach()
				defer WhenGettingProcessableThreadsForAgentAfterEach(suite)

				slackChannel := model_builders.NewSlackChannelBuilder().Build()
				slackAgent := model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleHuman).
					Build()

				_, err := suite.slackChannelRepository.Save(slackChannel)
				assert.NoError(t, err)

				slackAgent, err = suite.slackAgentRepository.Save(slackAgent)
				assert.NoError(t, err)

				_, err = suite.slackChannelRepository.AddAgent(slackChannel, slackAgent)
				assert.NoError(t, err)

				slackThread := model_builders.NewSlackThreadBuilder().
					WithId(model_builders.GenerateRandomSlackMessageID()).
					WithChannelId(slackChannel.Id).
					WithIsConcluded(false).
					Build()

				_, err = suite.sut.Save(slackThread)
				assert.NoError(t, err)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId("author-user-id").
					Build()

				_, err = suite.slackUserRepository.Save(slackUser)
				assert.NoError(t, err)

				// Message authored by the human agent (physical human behind the agent)
				message := model_builders.NewSlackMessageBuilder().
					WithId(model_builders.GenerateRandomSlackMessageID()).
					WithIsFromUser(true).
					WithIsFromAgent(true).
					WithThreadId(slackThread.Id).
					WithChannelId(slackChannel.Id).
					WithAuthorUserId(slackUser.Id).
					WithAuthorCollectionAddress(slackAgent.CollectionAddress).
					WithAuthorCollectionTokenId(slackAgent.CollectionTokenId).
					WithContent("General chat without mentions").
					Build()

				_, err = suite.slackMessageRepository.Save(message)
				assert.NoError(t, err)

				processableThreads, err := suite.sut.GetProcessableThreadsForAgent(slackAgent, 30)

				assert.NoError(t, err)
				assert.Len(t, processableThreads, 1)
				assert.False(t, processableThreads[0].IsConcluded)
				assert.Equal(t, slackThread.Id, processableThreads[0].Id)
				assert.Equal(t, slackThread.ChannelId, processableThreads[0].ChannelId)
			})
		})

		t.Run("And agent is mentioned + the underlying physical human send a message", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return processable threads", func(t *testing.T) {
				t.Parallel()

				suite := WhenGettingProcessableThreadsForAgentBeforeEach()
				defer WhenGettingProcessableThreadsForAgentAfterEach(suite)

				slackChannel := model_builders.NewSlackChannelBuilder().Build()
				slackAgent := model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleHuman).
					Build()

				_, err := suite.slackChannelRepository.Save(slackChannel)
				assert.NoError(t, err)

				slackAgent, err = suite.slackAgentRepository.Save(slackAgent)
				assert.NoError(t, err)

				_, err = suite.slackChannelRepository.AddAgent(slackChannel, slackAgent)
				assert.NoError(t, err)

				slackThread := model_builders.NewSlackThreadBuilder().
					WithId(model_builders.GenerateRandomSlackMessageID()).
					WithChannelId(slackChannel.Id).
					WithIsConcluded(false).
					Build()

				_, err = suite.sut.Save(slackThread)
				assert.NoError(t, err)

				slackUser1 := model_builders.NewSlackUserBuilder().
					WithId("author-user-id").
					WithPairedHumanAgentCollectionAddress(&slackAgent.CollectionAddress).
					WithPairedHumanAgentCollectionTokenId(&slackAgent.CollectionTokenId).
					Build()

				_, err = suite.slackUserRepository.Save(slackUser1)
				assert.NoError(t, err)

				// Message authored by the human agent (physical human behind the agent)
				message1 := model_builders.NewSlackMessageBuilder().
					WithId(model_builders.GenerateRandomSlackMessageID()).
					WithIsFromUser(true).
					WithIsFromAgent(true).
					WithThreadId(slackThread.Id).
					WithChannelId(slackChannel.Id).
					WithAuthorUserId(slackUser1.Id).
					WithAuthorCollectionAddress(slackAgent.CollectionAddress).
					WithAuthorCollectionTokenId(slackAgent.CollectionTokenId).
					WithContent("General chat without mentions").
					Build()

				_, err = suite.slackMessageRepository.Save(message1)
				assert.NoError(t, err)

				slackUser2 := model_builders.NewSlackUserBuilder().
					WithId("author-user-id-2").
					Build()

				_, err = suite.slackUserRepository.Save(slackUser2)
				assert.NoError(t, err)

				// Message mentioning the human agent's ENS domain
				message2 := model_builders.NewSlackMessageBuilder().
					WithId(model_builders.GenerateRandomSlackMessageID()).
					WithThreadId(slackThread.Id).
					WithChannelId(slackChannel.Id).
					WithContent(fmt.Sprintf("Hey %s, can you help?", slackAgent.EnsDomain)).
					WithAuthorUserId(slackUser2.Id).
					Build()

				_, err = suite.slackMessageRepository.Save(message2)
				assert.NoError(t, err)

				processableThreads, err := suite.sut.GetProcessableThreadsForAgent(slackAgent, 30)

				assert.NoError(t, err)
				assert.Len(t, processableThreads, 1)
				assert.False(t, processableThreads[0].IsConcluded)
				assert.Equal(t, slackThread.Id, processableThreads[0].Id)
				assert.Equal(t, slackThread.ChannelId, processableThreads[0].ChannelId)
			})
		})

		t.Run("And agent is not mentioned and underlying physical human didn't send a message", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return empty list", func(t *testing.T) {
				t.Parallel()

				suite := WhenGettingProcessableThreadsForAgentBeforeEach()
				defer WhenGettingProcessableThreadsForAgentAfterEach(suite)

				slackChannel := model_builders.NewSlackChannelBuilder().Build()
				slackAgent := model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleHuman).
					Build()

				_, err := suite.slackChannelRepository.Save(slackChannel)
				assert.NoError(t, err)

				slackAgent, err = suite.slackAgentRepository.Save(slackAgent)
				assert.NoError(t, err)

				_, err = suite.slackChannelRepository.AddAgent(slackChannel, slackAgent)
				assert.NoError(t, err)

				slackThread := model_builders.NewSlackThreadBuilder().
					WithId(model_builders.GenerateRandomSlackMessageID()).
					WithChannelId(slackChannel.Id).
					WithIsConcluded(false).
					Build()

				_, err = suite.sut.Save(slackThread)
				assert.NoError(t, err)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId("author-user-id").
					Build()

				_, err = suite.slackUserRepository.Save(slackUser)
				assert.NoError(t, err)

				// Message without mention and not authored by the human agent
				message := model_builders.NewSlackMessageBuilder().
					WithId(model_builders.GenerateRandomSlackMessageID()).
					WithThreadId(slackThread.Id).
					WithChannelId(slackChannel.Id).
					WithContent("General chat without mentions").
					WithAuthorUserId(slackUser.Id).
					Build()

				_, err = suite.slackMessageRepository.Save(message)
				assert.NoError(t, err)

				processableThreads, err := suite.sut.GetProcessableThreadsForAgent(slackAgent, 30)

				assert.NoError(t, err)
				assert.Empty(t, processableThreads)
			})
		})
	})

	t.Run("Given agent is expert", func(t *testing.T) {
		t.Parallel()

		t.Run("And agent is mentioned since his last message", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return processable threads", func(t *testing.T) {
				t.Parallel()

				suite := WhenGettingProcessableThreadsForAgentBeforeEach()
				defer WhenGettingProcessableThreadsForAgentAfterEach(suite)

				slackChannel := model_builders.NewSlackChannelBuilder().Build()

				slackAgent := model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleExpert).
					Build()

				_, err := suite.slackChannelRepository.Save(slackChannel)
				assert.NoError(t, err)

				slackAgent, err = suite.slackAgentRepository.Save(slackAgent)
				assert.NoError(t, err)

				_, err = suite.slackChannelRepository.AddAgent(slackChannel, slackAgent)
				assert.NoError(t, err)

				slackThread := model_builders.NewSlackThreadBuilder().
					WithId(model_builders.GenerateRandomSlackMessageID()).
					WithChannelId(slackChannel.Id).
					WithIsConcluded(false).
					Build()

				_, err = suite.sut.Save(slackThread)
				assert.NoError(t, err)

				firstMessage := model_builders.NewSlackMessageBuilder().
					WithId(model_builders.GenerateRandomSlackMessageID()).
					WithThreadId(slackThread.Id).
					WithChannelId(slackChannel.Id).
					WithContent("Hello, how are you ?").
					WithAuthorCollectionAddress(slackAgent.CollectionAddress).
					WithAuthorCollectionTokenId(slackAgent.CollectionTokenId).
					Build()

				_, err = suite.slackMessageRepository.Save(firstMessage)
				assert.NoError(t, err)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId("author-user-id").
					Build()

				_, err = suite.slackUserRepository.Save(slackUser)
				assert.NoError(t, err)

				secondMessage := model_builders.NewSlackMessageBuilder().
					WithId(model_builders.GenerateRandomSlackMessageID()).
					WithThreadId(slackThread.Id).
					WithChannelId(slackChannel.Id).
					WithContent(fmt.Sprintf("Hello %s, how are you?", slackAgent.EnsDomain)).
					WithAuthorUserId(slackUser.Id).
					Build()

				_, err = suite.slackMessageRepository.Save(secondMessage)
				assert.NoError(t, err)

				processableThreads, err := suite.sut.GetProcessableThreadsForAgent(slackAgent, 30)
				assert.NoError(t, err)

				assert.Len(t, processableThreads, 1)
			})
		})

		t.Run("And last mention is from agent himself", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return empty list", func(t *testing.T) {
				t.Parallel()

				suite := WhenGettingProcessableThreadsForAgentBeforeEach()
				defer WhenGettingProcessableThreadsForAgentAfterEach(suite)

				slackChannel := model_builders.NewSlackChannelBuilder().Build()

				slackAgent := model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleExpert).
					Build()

				_, err := suite.slackChannelRepository.Save(slackChannel)
				assert.NoError(t, err)

				slackAgent, err = suite.slackAgentRepository.Save(slackAgent)
				assert.NoError(t, err)

				_, err = suite.slackChannelRepository.AddAgent(slackChannel, slackAgent)
				assert.NoError(t, err)

				slackThread := model_builders.NewSlackThreadBuilder().
					WithId(model_builders.GenerateRandomSlackMessageID()).
					WithChannelId(slackChannel.Id).
					WithIsConcluded(false).
					Build()
				_, err = suite.sut.Save(slackThread)
				assert.NoError(t, err)

				firstMessage := model_builders.NewSlackMessageBuilder().
					WithId(model_builders.GenerateRandomSlackMessageID()).
					WithThreadId(slackThread.Id).
					WithChannelId(slackChannel.Id).
					WithContent(fmt.Sprintf("Hello @%s, how are you?", slackAgent.EnsDomain)).
					WithAuthorCollectionAddress(slackAgent.CollectionAddress).
					WithAuthorCollectionTokenId(slackAgent.CollectionTokenId).
					Build()

				_, err = suite.slackMessageRepository.Save(firstMessage)
				assert.NoError(t, err)

				secondMessage := model_builders.NewSlackMessageBuilder().
					WithId(model_builders.GenerateRandomSlackMessageID()).
					WithThreadId(slackThread.Id).
					WithChannelId(slackChannel.Id).
					WithContent(fmt.Sprintf("Hello, I'm @%s, I'm fine and you?", slackAgent.EnsDomain)).
					WithAuthorCollectionAddress(slackAgent.CollectionAddress).
					WithAuthorCollectionTokenId(slackAgent.CollectionTokenId).
					Build()

				_, err = suite.slackMessageRepository.Save(secondMessage)
				assert.NoError(t, err)

				processableThreads, err := suite.sut.GetProcessableThreadsForAgent(slackAgent, 30)

				assert.NoError(t, err)
				assert.Empty(t, processableThreads)
			})
		})

		t.Run("And agent is not mentioned since his last message", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return empty list", func(t *testing.T) {
				t.Parallel()

				suite := WhenGettingProcessableThreadsForAgentBeforeEach()
				defer WhenGettingProcessableThreadsForAgentAfterEach(suite)

				slackChannel := model_builders.NewSlackChannelBuilder().Build()

				slackAgent := model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleExpert).
					Build()

				_, err := suite.slackChannelRepository.Save(slackChannel)
				assert.NoError(t, err)

				slackAgent, err = suite.slackAgentRepository.Save(slackAgent)
				assert.NoError(t, err)

				_, err = suite.slackChannelRepository.AddAgent(slackChannel, slackAgent)
				assert.NoError(t, err)

				slackThread := model_builders.NewSlackThreadBuilder().
					WithId(model_builders.GenerateRandomSlackMessageID()).
					WithChannelId(slackChannel.Id).
					WithIsConcluded(false).
					Build()

				_, err = suite.sut.Save(slackThread)
				assert.NoError(t, err)

				firstMessage := model_builders.NewSlackMessageBuilder().
					WithId(model_builders.GenerateRandomSlackMessageID()).
					WithThreadId(slackThread.Id).
					WithChannelId(slackChannel.Id).
					WithContent("Hello, how are you ?").
					WithAuthorCollectionAddress(slackAgent.CollectionAddress).
					WithAuthorCollectionTokenId(slackAgent.CollectionTokenId).
					Build()

				_, err = suite.slackMessageRepository.Save(firstMessage)
				assert.NoError(t, err)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId("author-user-id").
					Build()

				_, err = suite.slackUserRepository.Save(slackUser)
				assert.NoError(t, err)

				secondMessage := model_builders.NewSlackMessageBuilder().
					WithId(model_builders.GenerateRandomSlackMessageID()).
					WithThreadId(slackThread.Id).
					WithChannelId(slackChannel.Id).
					WithContent("I am fine, thanks!").
					WithAuthorUserId(slackUser.Id).
					Build()

				_, err = suite.slackMessageRepository.Save(secondMessage)
				assert.NoError(t, err)

				processableThreads, err := suite.sut.GetProcessableThreadsForAgent(slackAgent, 30)

				assert.NoError(t, err)
				assert.Empty(t, processableThreads)
			})
		})
	})
}
