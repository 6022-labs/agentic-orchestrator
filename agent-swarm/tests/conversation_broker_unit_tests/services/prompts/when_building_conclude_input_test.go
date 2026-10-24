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

type WhenBuildingConcludeInputTestingSuite struct {
	sut *prompts.ConversationPromptBuilder

	mockAgentPromptTemplateRepository *repositories_mocks.MockAgentPromptTemplateRepositoryInterface
	mockRegistryAgentResolver         *interfaces_mocks.MockRegistryAgentResolverInterface
}

func WhenBuildingConcludeInputBeforeEach(t *testing.T) *WhenBuildingConcludeInputTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentPromptTemplateRepository := repositories_mocks.NewMockAgentPromptTemplateRepositoryInterface(
		mockController,
	)
	mockRegistryAgentResolver := interfaces_mocks.NewMockRegistryAgentResolverInterface(mockController)

	return &WhenBuildingConcludeInputTestingSuite{
		sut: prompts.NewConversationPromptBuilder(mockAgentPromptTemplateRepository, mockRegistryAgentResolver),

		mockAgentPromptTemplateRepository: mockAgentPromptTemplateRepository,
		mockRegistryAgentResolver:         mockRegistryAgentResolver,
	}
}

func TestWhenBuildingConcludeInput(t *testing.T) {
	t.Parallel()

	t.Run("Given an agent stored with a checksummed collection address", func(t *testing.T) {
		t.Parallel()

		agent := model_builders.NewAgentParticipantBuilder().Build()
		agent.CollectionAddress = lo.ToPtr("0xAbC")
		conversation := model_builders.NewConversationBuilder().Build()
		message := model_builders.NewConversationMessageBuilder().
			WithConversationID(conversation.ID).
			WithParticipantID(agent.ID).
			Build()
		baseTurn := &models.RemoteAgentInput{Instructions: "instructions"}

		initSuite := func(suite *WhenBuildingConcludeInputTestingSuite) {
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(gomock.Any(),
					"0xabc",
					*agent.CollectionTokenId,
					enums.AgentPromptTemplateKeyFacilitatorConclude,
				).
				Return(nil, nil).
				Times(1)
			suite.mockRegistryAgentResolver.EXPECT().
				ResolveByEnsDomain(gomock.Any(), gomock.Any()).
				Return(nil, nil).
				AnyTimes()
		}

		t.Run("Should look its prompt template up by the lowercased address", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingConcludeInputBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.BuildConcludeInput(context.Background(),
				baseTurn,
				nil,
				agent,
				lo.ToPtr(enums.ParticipantRoleFacilitator),
				conversation,
				[]*models.ConversationMessage{message},
				[]*models.Participant{agent},
				nil,
			)

			assert.NoError(t, err)
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
			Input: []*models.RemoteAgentInputMessage{
				{Role: enums.RemoteAgentInputRoleUser, Content: "prompt"},
			},
			Metadata:             map[string]string{"conversationId": conversation.ID.String()},
			RemoteConversationId: "3f2b7c41-9a1e-4d0b-8c33-2f6a1b9d5c91",
		}
		concludeSessionId := uuid.MustParse("3f2b7c41-9a1e-4d0b-8c33-2f6a1b9d5c91")

		initSuite := func(suite *WhenBuildingConcludeInputTestingSuite) {
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(gomock.Any(),
					*agent.CollectionAddress,
					*agent.CollectionTokenId,
					enums.AgentPromptTemplateKeyFacilitatorConclude,
				).
				Return(nil, nil).
				Times(1)
			suite.mockRegistryAgentResolver.EXPECT().
				ResolveByEnsDomain(gomock.Any(), gomock.Any()).
				Return(nil, nil).
				AnyTimes()
		}

		t.Run("Should append the completions and the conclude prompt to the input", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingConcludeInputBeforeEach(t)
			initSuite(suite)

			concludeTurn, err := suite.sut.BuildConcludeInput(context.Background(),
				baseTurn,
				[]string{"first completion"},
				agent,
				lo.ToPtr(enums.ParticipantRoleFacilitator),
				conversation,
				messages,
				participants,
				&concludeSessionId,
			)

			assert.NoError(t, err)
			assert.Equal(t, baseTurn.Instructions, concludeTurn.Instructions)
			assert.Equal(t, baseTurn.RemoteConversationId, concludeTurn.RemoteConversationId)
			assert.Len(t, concludeTurn.Input, 3)
			assert.Equal(t, enums.RemoteAgentInputRoleAssistant, concludeTurn.Input[1].Role)
			assert.Equal(t, "first completion", concludeTurn.Input[1].Content)
			assert.Equal(t, enums.RemoteAgentInputRoleUser, concludeTurn.Input[2].Role)
			assert.NotEmpty(t, concludeTurn.Input[2].Content)
		})
	})
}
