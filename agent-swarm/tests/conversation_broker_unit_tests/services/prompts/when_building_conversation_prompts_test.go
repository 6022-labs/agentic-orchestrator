package prompts_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/prompts"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenBuildingConversationPromptsTestingSuite struct {
	sut *prompts.ConversationPromptBuilder

	mockAgentPromptTemplateRepository *repositories_mocks.MockAgentPromptTemplateRepositoryInterface
	mockRegistryAgentResolver         *interfaces_mocks.MockRegistryAgentResolverInterface
}

func WhenBuildingConversationPromptsBeforeEach(t *testing.T) *WhenBuildingConversationPromptsTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentPromptTemplateRepository := repositories_mocks.NewMockAgentPromptTemplateRepositoryInterface(mockController)
	mockRegistryAgentResolver := interfaces_mocks.NewMockRegistryAgentResolverInterface(mockController)

	sut := prompts.NewConversationPromptBuilder(mockAgentPromptTemplateRepository, mockRegistryAgentResolver)

	return &WhenBuildingConversationPromptsTestingSuite{
		sut: sut,

		mockAgentPromptTemplateRepository: mockAgentPromptTemplateRepository,
		mockRegistryAgentResolver:         mockRegistryAgentResolver,
	}
}

func TestWhenBuildingConversationPrompts(t *testing.T) {
	t.Parallel()

	t.Run("Given a clone agent with stored template overrides", func(t *testing.T) {
		t.Parallel()

		agent := model_builders.NewAgentParticipantBuilder().Build()
		user := model_builders.NewUserParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()

		userMessage := model_builders.NewConversationMessageBuilder().
			WithConversationID(conversation.ID).
			WithParticipantID(user.ID).
			WithContent("hello agent").
			Build()
		agentMessage := model_builders.NewConversationMessageBuilder().
			WithConversationID(conversation.ID).
			WithParticipantID(agent.ID).
			WithContent("previous reply").
			Build()
		messages := []*models.ConversationMessage{userMessage, agentMessage}
		participants := []*models.Participant{agent, user}

		initSuite := func(suite *WhenBuildingConversationPromptsTestingSuite) {
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(*agent.CollectionAddress, *agent.CollectionTokenId, enums.AgentPromptTemplateKeyCloneSystem).
				Return(model_builders.NewAgentPromptTemplateBuilder().WithContent("SYSTEM {{.ChannelDescription}}").Build(), nil).
				Times(1)
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(*agent.CollectionAddress, *agent.CollectionTokenId, enums.AgentPromptTemplateKeyCloneDefault).
				Return(model_builders.NewAgentPromptTemplateBuilder().WithContent("REPLY TO {{.LatestRequest.Content}}").Build(), nil).
				Times(1)
			// ChannelMembers + Latest/Initial authors resolve through the registry.
			suite.mockRegistryAgentResolver.EXPECT().
				ResolveByEnsDomain(gomock.Any()).
				Return(nil, nil).
				AnyTimes()
		}

		t.Run("Should build instructions and input from the overrides", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingConversationPromptsBeforeEach(t)
			initSuite(suite)

			turn, err := suite.sut.BuildInput("slack", agent, lo.ToPtr(enums.ParticipantRoleClone), conversation, messages, participants)

			assert.NoError(t, err)
			assert.Equal(t, "SYSTEM "+*conversation.Description, turn.Instructions)
			// mapped history + final turn prompt
			assert.Len(t, turn.Input, 3)
			assert.Equal(t, enums.RemoteAgentInputRoleUser, turn.Input[0].Role)
			assert.Contains(t, turn.Input[0].Content, "hello agent")
			assert.Contains(t, turn.Input[0].Content, *user.DisplayName)
			assert.Equal(t, enums.RemoteAgentInputRoleAssistant, turn.Input[1].Role)
			assert.Equal(t, "previous reply", turn.Input[1].Content)
			assert.Equal(t, enums.RemoteAgentInputRoleUser, turn.Input[2].Role)
			assert.Equal(t, "REPLY TO previous reply", turn.Input[2].Content)
			assert.Equal(t, conversation.ID.String(), turn.Metadata["conversationId"])
			assert.Equal(t, "slack", turn.Metadata["platform"])
			assert.Equal(t, *agent.EnsDomain, turn.Metadata["ensDomain"])
		})
	})

	t.Run("Given no stored templates exist", func(t *testing.T) {
		t.Parallel()

		agent := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		message := model_builders.NewConversationMessageBuilder().
			WithConversationID(conversation.ID).
			WithParticipantID(agent.ID).
			Build()

		initSuite := func(suite *WhenBuildingConversationPromptsTestingSuite) {
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil).
				Times(2)
			suite.mockRegistryAgentResolver.EXPECT().
				ResolveByEnsDomain(gomock.Any()).
				Return(nil, nil).
				AnyTimes()
		}

		t.Run("Should fall back to the built-in role defaults", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingConversationPromptsBeforeEach(t)
			initSuite(suite)

			turn, err := suite.sut.BuildInput("web", agent, lo.ToPtr(enums.ParticipantRoleClone), conversation, []*models.ConversationMessage{message}, []*models.Participant{agent})

			assert.NoError(t, err)
			assert.NotEmpty(t, turn.Instructions)
			assert.NotEmpty(t, turn.Input)
		})
	})

	t.Run("Given a facilitator conclude turn", func(t *testing.T) {
		t.Parallel()

		agent := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		message := model_builders.NewConversationMessageBuilder().
			WithConversationID(conversation.ID).
			WithParticipantID(agent.ID).
			Build()
		messages := []*models.ConversationMessage{message}
		participants := []*models.Participant{agent}

		baseTurn := &models.RemoteAgentInput{
			Instructions: "instructions",
			Input:        []*models.RemoteAgentInputMessage{{Role: enums.RemoteAgentInputRoleUser, Content: "prompt"}},
			Metadata:     map[string]string{"conversationId": conversation.ID.String()},
		}

		initSuite := func(suite *WhenBuildingConversationPromptsTestingSuite) {
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(*agent.CollectionAddress, *agent.CollectionTokenId, enums.AgentPromptTemplateKeyFacilitatorConclude).
				Return(nil, nil).
				Times(1)
			suite.mockRegistryAgentResolver.EXPECT().
				ResolveByEnsDomain(gomock.Any()).
				Return(nil, nil).
				AnyTimes()
		}

		t.Run("Should append the completions and the conclude prompt to the input", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingConversationPromptsBeforeEach(t)
			initSuite(suite)

			concludeTurn, err := suite.sut.BuildConcludeInput(baseTurn, []string{"first completion"}, agent, lo.ToPtr(enums.ParticipantRoleFacilitator), conversation, messages, participants)

			assert.NoError(t, err)
			assert.Equal(t, baseTurn.Instructions, concludeTurn.Instructions)
			assert.Len(t, concludeTurn.Input, 3)
			assert.Equal(t, enums.RemoteAgentInputRoleAssistant, concludeTurn.Input[1].Role)
			assert.Equal(t, "first completion", concludeTurn.Input[1].Content)
			assert.Equal(t, enums.RemoteAgentInputRoleUser, concludeTurn.Input[2].Role)
			assert.NotEmpty(t, concludeTurn.Input[2].Content)
		})
	})

	t.Run("Given the agent has the human role", func(t *testing.T) {
		t.Parallel()

		agent := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		message := model_builders.NewConversationMessageBuilder().WithConversationID(conversation.ID).Build()

		t.Run("Should return an error because humans have no templates", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingConversationPromptsBeforeEach(t)

			_, err := suite.sut.BuildInput("web", agent, lo.ToPtr(enums.ParticipantRoleHuman), conversation, []*models.ConversationMessage{message}, []*models.Participant{agent})

			assert.Error(t, err)
		})
	})
}
