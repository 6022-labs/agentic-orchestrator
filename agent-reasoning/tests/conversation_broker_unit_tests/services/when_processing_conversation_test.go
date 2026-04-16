package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WhenProcessingConversationTestingSuite struct {
	sut *services.ConversationProcessor

	mockConversationConclude          *interfaces_mocks.MockConversationConcludeInterface
	mockBridgeCallbackRequester       *interfaces_mocks.MockBridgeCallbackRequesterInterface
	mockAgentCompletionsRequester     *interfaces_mocks.MockAgentCompletionsRequesterInterface
	mockBridgeRepository              *repositories_mocks.MockBridgeRepositoryInterface
	mockConversationRepository        *repositories_mocks.MockConversationRepositoryInterface
	mockConversationMessageRepository *repositories_mocks.MockConversationMessageRepositoryInterface
	mockParticipantRepository         *repositories_mocks.MockParticipantRepositoryInterface
}

func WhenProcessingConversationBeforeEach(t *testing.T) *WhenProcessingConversationTestingSuite {
	mockController := gomock.NewController(t)

	mockConversationConclude := interfaces_mocks.NewMockConversationConcludeInterface(mockController)
	mockBridgeCallbackRequester := interfaces_mocks.NewMockBridgeCallbackRequesterInterface(mockController)
	mockAgentCompletionsRequester := interfaces_mocks.NewMockAgentCompletionsRequesterInterface(mockController)
	mockBridgeRepository := repositories_mocks.NewMockBridgeRepositoryInterface(mockController)
	mockConversationRepository := repositories_mocks.NewMockConversationRepositoryInterface(mockController)
	mockConversationMessageRepository := repositories_mocks.NewMockConversationMessageRepositoryInterface(mockController)
	mockParticipantRepository := repositories_mocks.NewMockParticipantRepositoryInterface(mockController)

	processingSettings := &settings.ProcessingSettings{
		ParticipantProcessingInterval:               20,
		FacilitatorTickInterval:                     30,
		ConversationForceConcludeByFacilitatorAfter: 300,
		ConversationForceConcludeAfter:              600,
	}

	sut := services.NewConversationProcessor(
		zap.NewNop(),
		processingSettings,
		mockConversationConclude,
		mockBridgeCallbackRequester,
		mockAgentCompletionsRequester,
		mockBridgeRepository,
		mockConversationRepository,
		mockConversationMessageRepository,
		mockParticipantRepository,
	)

	return &WhenProcessingConversationTestingSuite{
		sut: sut,

		mockConversationConclude:          mockConversationConclude,
		mockBridgeCallbackRequester:       mockBridgeCallbackRequester,
		mockAgentCompletionsRequester:     mockAgentCompletionsRequester,
		mockBridgeRepository:              mockBridgeRepository,
		mockConversationRepository:        mockConversationRepository,
		mockConversationMessageRepository: mockConversationMessageRepository,
		mockParticipantRepository:         mockParticipantRepository,
	}
}

func TestWhenProcessingConversation(t *testing.T) {
	t.Parallel()

	t.Run("Given the conversation is past the force-conclude threshold", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		conversation.CreatedAt = time.Now().Add(-2 * time.Hour)

		initSuite := func(suite *WhenProcessingConversationTestingSuite) {
			suite.mockConversationConclude.EXPECT().Conclude(conversation.ID).Return(nil).Times(1)
		}

		t.Run("Should force conclude and stop", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingConversationBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ProcessConversation(participant, conversation)

			assert.NoError(t, err)
		})
	})

	t.Run("Given the bridge does not exist", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		conversation.CreatedAt = time.Now()

		initSuite := func(suite *WhenProcessingConversationTestingSuite) {
			suite.mockBridgeRepository.EXPECT().GetById(conversation.BridgeID).Return(nil, gorm.ErrRecordNotFound).Times(1)
			suite.mockConversationConclude.EXPECT().Conclude(conversation.ID).Return(nil).Times(1)
		}

		t.Run("Should conclude the orphan conversation and stop", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingConversationBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ProcessConversation(participant, conversation)

			assert.NoError(t, err)
		})
	})

	t.Run("Given there are no processable messages for the participant", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		conversation.CreatedAt = time.Now()
		bridge := model_builders.NewBridgeBuilder().WithID(conversation.BridgeID).Build()

		initSuite := func(suite *WhenProcessingConversationTestingSuite) {
			suite.mockBridgeRepository.EXPECT().GetById(conversation.BridgeID).Return(bridge, nil).Times(1)
			suite.mockConversationMessageRepository.EXPECT().
				ListProcessableForParticipantAndConversation(participant, conversation.ID).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should return without calling agent_compute", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingConversationBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ProcessConversation(participant, conversation)

			assert.NoError(t, err)
		})
	})

	t.Run("Given the participant is human", func(t *testing.T) {
		t.Parallel()

		humanRole := enums.ParticipantRoleHuman
		participant := model_builders.NewAgentParticipantBuilder().Build()
		participant.Role = &humanRole

		conversation := model_builders.NewConversationBuilder().Build()
		conversation.CreatedAt = time.Now()
		bridge := model_builders.NewBridgeBuilder().WithID(conversation.BridgeID).Build()

		message := model_builders.NewConversationMessageBuilder().WithConversationID(conversation.ID).Build()
		messages := []*models.ConversationMessage{message}
		participants := []*models.Participant{participant}

		initSuite := func(suite *WhenProcessingConversationTestingSuite) {
			suite.mockBridgeRepository.EXPECT().GetById(conversation.BridgeID).Return(bridge, nil).Times(1)
			suite.mockConversationMessageRepository.EXPECT().
				ListProcessableForParticipantAndConversation(participant, conversation.ID).
				Return(messages, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				ListByConversationId(conversation.ID).
				Return(participants, nil).
				Times(1)
			suite.mockAgentCompletionsRequester.EXPECT().
				Completions(string(bridge.Kind), participant, conversation, messages, participants, false).
				Return(nil, false, nil).
				Times(1)
		}

		t.Run("Should call compute but skip sending completions back", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingConversationBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ProcessConversation(participant, conversation)

			assert.NoError(t, err)
		})
	})

	t.Run("Given an AI participant produces completions that conclude the conversation", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build() // role = clone by default
		conversation := model_builders.NewConversationBuilder().Build()
		conversation.CreatedAt = time.Now()
		bridge := model_builders.NewBridgeBuilder().WithID(conversation.BridgeID).Build()

		message := model_builders.NewConversationMessageBuilder().WithConversationID(conversation.ID).Build()
		messages := []*models.ConversationMessage{message}
		participants := []*models.Participant{participant}
		completions := []*string{lo.ToPtr("first reply"), lo.ToPtr("second reply")}

		initSuite := func(suite *WhenProcessingConversationTestingSuite) {
			suite.mockBridgeRepository.EXPECT().GetById(conversation.BridgeID).Return(bridge, nil).Times(1)
			suite.mockConversationMessageRepository.EXPECT().
				ListProcessableForParticipantAndConversation(participant, conversation.ID).
				Return(messages, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				ListByConversationId(conversation.ID).
				Return(participants, nil).
				Times(1)
			suite.mockAgentCompletionsRequester.EXPECT().
				Completions(string(bridge.Kind), participant, conversation, messages, participants, false).
				Return(completions, true, nil).
				Times(1)
			suite.mockConversationRepository.EXPECT().GetById(conversation.ID).Return(conversation, nil).Times(1)
			suite.mockBridgeCallbackRequester.EXPECT().
				PostCompletion(bridge, conversation, participant, "first reply").
				Return(nil).
				Times(1)
			suite.mockBridgeCallbackRequester.EXPECT().
				PostCompletion(bridge, conversation, participant, "second reply").
				Return(nil).
				Times(1)
			suite.mockConversationConclude.EXPECT().Conclude(conversation.ID).Return(nil).Times(1)
		}

		t.Run("Should post each completion to the bridge and conclude the conversation", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingConversationBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ProcessConversation(participant, conversation)

			assert.NoError(t, err)
		})
	})

	t.Run("Given an AI participant produces completions that keep the conversation going", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		conversation.CreatedAt = time.Now()
		bridge := model_builders.NewBridgeBuilder().WithID(conversation.BridgeID).Build()

		messages := []*models.ConversationMessage{model_builders.NewConversationMessageBuilder().WithConversationID(conversation.ID).Build()}
		participants := []*models.Participant{participant}
		completions := []*string{lo.ToPtr("only reply")}

		initSuite := func(suite *WhenProcessingConversationTestingSuite) {
			suite.mockBridgeRepository.EXPECT().GetById(conversation.BridgeID).Return(bridge, nil).Times(1)
			suite.mockConversationMessageRepository.EXPECT().
				ListProcessableForParticipantAndConversation(participant, conversation.ID).
				Return(messages, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().ListByConversationId(conversation.ID).Return(participants, nil).Times(1)
			suite.mockAgentCompletionsRequester.EXPECT().
				Completions(string(bridge.Kind), participant, conversation, messages, participants, false).
				Return(completions, false, nil).
				Times(1)
			suite.mockConversationRepository.EXPECT().GetById(conversation.ID).Return(conversation, nil).Times(1)
			suite.mockBridgeCallbackRequester.EXPECT().
				PostCompletion(bridge, conversation, participant, "only reply").
				Return(nil).
				Times(1)
		}

		t.Run("Should post the completion to the bridge and leave the conversation active", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingConversationBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ProcessConversation(participant, conversation)

			assert.NoError(t, err)
		})
	})

	t.Run("Given the conversation got concluded between compute and reload", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		conversation.CreatedAt = time.Now()
		bridge := model_builders.NewBridgeBuilder().WithID(conversation.BridgeID).Build()

		messages := []*models.ConversationMessage{model_builders.NewConversationMessageBuilder().WithConversationID(conversation.ID).Build()}
		participants := []*models.Participant{participant}
		completions := []*string{lo.ToPtr("a")}

		reloaded := model_builders.NewConversationBuilder().WithID(conversation.ID).WithStatus(enums.ConversationStatusConcluded).Build()

		initSuite := func(suite *WhenProcessingConversationTestingSuite) {
			suite.mockBridgeRepository.EXPECT().GetById(conversation.BridgeID).Return(bridge, nil).Times(1)
			suite.mockConversationMessageRepository.EXPECT().
				ListProcessableForParticipantAndConversation(participant, conversation.ID).
				Return(messages, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().ListByConversationId(conversation.ID).Return(participants, nil).Times(1)
			suite.mockAgentCompletionsRequester.EXPECT().
				Completions(string(bridge.Kind), participant, conversation, messages, participants, false).
				Return(completions, false, nil).
				Times(1)
			suite.mockConversationRepository.EXPECT().GetById(conversation.ID).Return(reloaded, nil).Times(1)
		}

		t.Run("Should skip posting completions and not error", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingConversationBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ProcessConversation(participant, conversation)

			assert.NoError(t, err)
		})
	})

	t.Run("Given agent_compute fails", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		conversation.CreatedAt = time.Now()
		bridge := model_builders.NewBridgeBuilder().WithID(conversation.BridgeID).Build()

		messages := []*models.ConversationMessage{model_builders.NewConversationMessageBuilder().WithConversationID(conversation.ID).Build()}

		initSuite := func(suite *WhenProcessingConversationTestingSuite) {
			suite.mockBridgeRepository.EXPECT().GetById(conversation.BridgeID).Return(bridge, nil).Times(1)
			suite.mockConversationMessageRepository.EXPECT().
				ListProcessableForParticipantAndConversation(participant, conversation.ID).
				Return(messages, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().ListByConversationId(conversation.ID).Return([]*models.Participant{participant}, nil).Times(1)
			suite.mockAgentCompletionsRequester.EXPECT().
				Completions(string(bridge.Kind), participant, conversation, messages, []*models.Participant{participant}, false).
				Return(nil, false, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingConversationBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ProcessConversation(participant, conversation)

			assert.Error(t, err)
		})
	})
}
