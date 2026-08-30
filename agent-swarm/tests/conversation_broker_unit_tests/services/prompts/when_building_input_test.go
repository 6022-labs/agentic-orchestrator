package prompts_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/prompts"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenBuildingInputTestingSuite struct {
	sut *prompts.ConversationPromptBuilder

	mockAgentPromptTemplateRepository *repositories_mocks.MockAgentPromptTemplateRepositoryInterface
	mockRegistryAgentResolver         *interfaces_mocks.MockRegistryAgentResolverInterface
}

func WhenBuildingInputBeforeEach(t *testing.T) *WhenBuildingInputTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentPromptTemplateRepository := repositories_mocks.NewMockAgentPromptTemplateRepositoryInterface(
		mockController,
	)
	mockRegistryAgentResolver := interfaces_mocks.NewMockRegistryAgentResolverInterface(mockController)

	sut := prompts.NewConversationPromptBuilder(mockAgentPromptTemplateRepository, mockRegistryAgentResolver)

	return &WhenBuildingInputTestingSuite{
		sut: sut,

		mockAgentPromptTemplateRepository: mockAgentPromptTemplateRepository,
		mockRegistryAgentResolver:         mockRegistryAgentResolver,
	}
}

func TestWhenBuildingInput(t *testing.T) {
	t.Parallel()

	t.Run("Given a clone agent", func(t *testing.T) {
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

		buildInput := func(
			suite *WhenBuildingInputTestingSuite,
		) (*models.RemoteAgentInput, error) {
			return suite.sut.BuildInput(context.Background(),
				"slack",
				agent,
				lo.ToPtr(enums.ParticipantRoleClone),
				conversation,
				messages,
				participants,
				nil,
			)
		}

		t.Run("Given stored template overrides", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenBuildingInputTestingSuite) {
				suite.mockAgentPromptTemplateRepository.EXPECT().
					GetByAgentAndKey(
						gomock.Any(),
						*agent.CollectionAddress,
						*agent.CollectionTokenId,
						enums.AgentPromptTemplateKeyCloneSystem,
					).
					Return(model_builders.NewAgentPromptTemplateBuilder().WithContent("SYSTEM {{.ChannelDescription}}").Build(), nil).
					Times(1)
				suite.mockAgentPromptTemplateRepository.EXPECT().
					GetByAgentAndKey(
						gomock.Any(),
						*agent.CollectionAddress,
						*agent.CollectionTokenId,
						enums.AgentPromptTemplateKeyCloneDefault,
					).
					Return(
						model_builders.NewAgentPromptTemplateBuilder().
							WithContent("REPLY TO {{.LatestRequest.Content}}").
							Build(),
						nil,
					).
					Times(1)
				// ChannelMembers + Latest/Initial authors resolve through the registry.
				suite.mockRegistryAgentResolver.EXPECT().
					ResolveByEnsDomain(gomock.Any(), gomock.Any()).
					Return(nil, nil).
					AnyTimes()
			}

			t.Run("Should build instructions and input from the overrides", func(t *testing.T) {
				t.Parallel()

				suite := WhenBuildingInputBeforeEach(t)
				initSuite(suite)

				turn, err := buildInput(suite)

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

		t.Run("Given no stored template", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenBuildingInputTestingSuite) {
				suite.mockAgentPromptTemplateRepository.EXPECT().
					GetByAgentAndKey(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(2)
				suite.mockRegistryAgentResolver.EXPECT().
					ResolveByEnsDomain(gomock.Any(), gomock.Any()).
					Return(nil, nil).
					AnyTimes()
			}

			t.Run("Should fall back to the built-in clone templates", func(t *testing.T) {
				t.Parallel()

				suite := WhenBuildingInputBeforeEach(t)
				initSuite(suite)

				turn, err := buildInput(suite)

				assert.NoError(t, err)
				assert.Contains(t, turn.Instructions, "You are taking part in a group conversation")
				assert.Contains(t, turn.Input[len(turn.Input)-1].Content, "Reply to the conversation.")
			})
		})
	})

	t.Run("Given an expert agent", func(t *testing.T) {
		t.Parallel()

		agent := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		message := model_builders.NewConversationMessageBuilder().
			WithConversationID(conversation.ID).
			WithParticipantID(agent.ID).
			Build()
		messages := []*models.ConversationMessage{message}
		participants := []*models.Participant{agent}
		sessionId := uuid.MustParse("3f2b7c41-9a1e-4d0b-8c33-2f6a1b9d5c91")

		buildInput := func(
			suite *WhenBuildingInputTestingSuite,
			memorySessionId *uuid.UUID,
		) (*models.RemoteAgentInput, error) {
			return suite.sut.BuildInput(context.Background(),
				"slack",
				agent,
				lo.ToPtr(enums.ParticipantRoleExpert),
				conversation,
				messages,
				participants,
				memorySessionId,
			)
		}

		initSuite := func(suite *WhenBuildingInputTestingSuite) {
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil).
				Times(2)
			suite.mockRegistryAgentResolver.EXPECT().
				ResolveByEnsDomain(gomock.Any(), gomock.Any()).
				Return(nil, nil).
				AnyTimes()
		}

		t.Run("Given no session is bound yet", func(t *testing.T) {
			t.Parallel()

			t.Run("Should ask for a session and key the turn on the conversation", func(t *testing.T) {
				t.Parallel()

				suite := WhenBuildingInputBeforeEach(t)
				initSuite(suite)

				turn, err := buildInput(suite, nil)

				assert.NoError(t, err)
				assert.Contains(t, turn.Input[len(turn.Input)-1].Content, "not attached to one of your memory sessions")
				assert.Equal(t, conversation.ID.String(), turn.RemoteConversationId)
				assert.Equal(t, conversation.ID.String(), turn.Metadata["conversationId"])
				assert.Equal(t, conversation.ID.String(), turn.Metadata["memorySessionId"])
			})
		})

		t.Run("Given a bound session", func(t *testing.T) {
			t.Parallel()

			t.Run("Should key the turn on the session and render its id", func(t *testing.T) {
				t.Parallel()

				suite := WhenBuildingInputBeforeEach(t)
				initSuite(suite)

				turn, err := buildInput(suite, &sessionId)

				assert.NoError(t, err)
				assert.Contains(t, turn.Input[len(turn.Input)-1].Content, sessionId.String())
				assert.Equal(t, sessionId.String(), turn.RemoteConversationId)
				// The broker's own conversation id stays in metadata: bridges key on it.
				assert.Equal(t, conversation.ID.String(), turn.Metadata["conversationId"])
				assert.Equal(t, sessionId.String(), turn.Metadata["memorySessionId"])
			})
		})
	})

	t.Run("Given a facilitator agent", func(t *testing.T) {
		t.Parallel()

		agent := model_builders.NewAgentParticipantBuilder().Build()
		user := model_builders.NewUserParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()

		message := model_builders.NewConversationMessageBuilder().
			WithConversationID(conversation.ID).
			WithParticipantID(user.ID).
			WithContent("where do we stand?").
			Build()
		messages := []*models.ConversationMessage{message}
		participants := []*models.Participant{agent, user}
		sessionId := uuid.MustParse("6d1c0b52-4a77-4f0e-9a1c-5f2d8e3b7a04")

		buildInput := func(
			suite *WhenBuildingInputTestingSuite,
			memorySessionId *uuid.UUID,
		) (*models.RemoteAgentInput, error) {
			return suite.sut.BuildInput(context.Background(),
				"slack",
				agent,
				lo.ToPtr(enums.ParticipantRoleFacilitator),
				conversation,
				messages,
				participants,
				memorySessionId,
			)
		}

		t.Run("Given no stored template", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenBuildingInputTestingSuite) {
				suite.mockAgentPromptTemplateRepository.EXPECT().
					GetByAgentAndKey(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(2)
				suite.mockRegistryAgentResolver.EXPECT().
					ResolveByEnsDomain(gomock.Any(), gomock.Any()).
					Return(nil, nil).
					AnyTimes()
			}

			t.Run("Should fall back to the built-in facilitator templates", func(t *testing.T) {
				t.Parallel()

				suite := WhenBuildingInputBeforeEach(t)
				initSuite(suite)

				turn, err := buildInput(suite, &sessionId)

				assert.NoError(t, err)
				assert.Contains(t, turn.Instructions, "You facilitate a group conversation")
				assert.Contains(t, turn.Input[len(turn.Input)-1].Content, "where do we stand?")
				assert.Contains(t, turn.Input[len(turn.Input)-1].Content, "Facilitate the next step.")
			})

			t.Run("Should key the turn on the bound session", func(t *testing.T) {
				t.Parallel()

				suite := WhenBuildingInputBeforeEach(t)
				initSuite(suite)

				turn, err := buildInput(suite, &sessionId)

				assert.NoError(t, err)
				assert.Equal(t, sessionId.String(), turn.RemoteConversationId)
				assert.Equal(t, sessionId.String(), turn.Metadata["memorySessionId"])
				// The broker's own conversation id stays in metadata: bridges key on it.
				assert.Equal(t, conversation.ID.String(), turn.Metadata["conversationId"])
			})
		})
	})

	t.Run("Given the agent has the human role", func(t *testing.T) {
		t.Parallel()

		agent := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		message := model_builders.NewConversationMessageBuilder().WithConversationID(conversation.ID).Build()

		t.Run("Should return an error because humans have no templates", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingInputBeforeEach(t)

			_, err := suite.sut.BuildInput(context.Background(),
				"web",
				agent,
				lo.ToPtr(enums.ParticipantRoleHuman),
				conversation,
				[]*models.ConversationMessage{message},
				[]*models.Participant{agent},
				nil,
			)

			assert.Error(t, err)
		})
	})
}
