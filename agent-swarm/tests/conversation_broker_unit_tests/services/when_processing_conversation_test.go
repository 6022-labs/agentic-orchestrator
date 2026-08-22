package services_test

import (
	"context"

	"errors"
	"testing"
	"time"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenProcessingConversationTestingSuite struct {
	sut *services.ConversationProcessor

	mockConversationConclude            *interfaces_mocks.MockConversationConcludeInterface
	mockBridgeCallbackRequester         *interfaces_mocks.MockBridgeCallbackRequesterInterface
	mockRemoteAgentCardResolver         *interfaces_mocks.MockRemoteAgentCardResolverInterface
	mockRemoteAgentCompletionsRequester *interfaces_mocks.MockRemoteAgentCompletionsRequesterInterface
	mockConversationPromptBuilder       *interfaces_mocks.MockConversationPromptBuilderInterface
	mockParticipantRoleResolver         *interfaces_mocks.MockParticipantRoleResolverInterface
	mockBridgeRepository                *repositories_mocks.MockBridgeRepositoryInterface
	mockConversationRepository          *repositories_mocks.MockConversationRepositoryInterface
	mockConversationMessageRepository   *repositories_mocks.MockConversationMessageRepositoryInterface
	mockParticipantRepository           *repositories_mocks.MockParticipantRepositoryInterface
	mockEmojiDecoder                    *interfaces_mocks.MockEmojiDecoderInterface
}

func WhenProcessingConversationBeforeEach(t *testing.T) *WhenProcessingConversationTestingSuite {
	mockController := gomock.NewController(t)

	mockConversationConclude := interfaces_mocks.NewMockConversationConcludeInterface(mockController)
	mockBridgeCallbackRequester := interfaces_mocks.NewMockBridgeCallbackRequesterInterface(mockController)
	mockRemoteAgentCardResolver := interfaces_mocks.NewMockRemoteAgentCardResolverInterface(mockController)
	mockRemoteAgentCompletionsRequester := interfaces_mocks.NewMockRemoteAgentCompletionsRequesterInterface(
		mockController,
	)
	mockConversationPromptBuilder := interfaces_mocks.NewMockConversationPromptBuilderInterface(mockController)
	mockParticipantRoleResolver := interfaces_mocks.NewMockParticipantRoleResolverInterface(mockController)
	mockBridgeRepository := repositories_mocks.NewMockBridgeRepositoryInterface(mockController)
	mockConversationRepository := repositories_mocks.NewMockConversationRepositoryInterface(mockController)
	mockConversationMessageRepository := repositories_mocks.NewMockConversationMessageRepositoryInterface(
		mockController,
	)
	mockParticipantRepository := repositories_mocks.NewMockParticipantRepositoryInterface(mockController)

	mockEmojiDecoder := interfaces_mocks.NewMockEmojiDecoderInterface(mockController)
	mockEmojiDecoder.EXPECT().DemojizeContent(gomock.Any()).
		DoAndReturn(func(content string) string { return content }).
		AnyTimes()

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
		mockRemoteAgentCardResolver,
		mockRemoteAgentCompletionsRequester,
		mockConversationPromptBuilder,
		mockParticipantRoleResolver,
		mockBridgeRepository,
		mockConversationRepository,
		mockConversationMessageRepository,
		mockParticipantRepository,
		mockEmojiDecoder,
	)

	return &WhenProcessingConversationTestingSuite{
		sut: sut,

		mockConversationConclude:            mockConversationConclude,
		mockBridgeCallbackRequester:         mockBridgeCallbackRequester,
		mockRemoteAgentCardResolver:         mockRemoteAgentCardResolver,
		mockRemoteAgentCompletionsRequester: mockRemoteAgentCompletionsRequester,
		mockConversationPromptBuilder:       mockConversationPromptBuilder,
		mockParticipantRoleResolver:         mockParticipantRoleResolver,
		mockBridgeRepository:                mockBridgeRepository,
		mockConversationRepository:          mockConversationRepository,
		mockConversationMessageRepository:   mockConversationMessageRepository,
		mockParticipantRepository:           mockParticipantRepository,
		mockEmojiDecoder:                    mockEmojiDecoder,
	}
}

func remoteAgentCardForTest() *models.RemoteAgentCard {
	return &models.RemoteAgentCard{
		Protocol:  "6022",
		Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
	}
}

func agentTurnRequestForTest() *models.RemoteAgentInput {
	return &models.RemoteAgentInput{
		Instructions: "instructions",
		Input:        []*models.RemoteAgentInputMessage{{Role: enums.RemoteAgentInputRoleUser, Content: "hello"}},
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
			suite.mockConversationConclude.EXPECT().Conclude(gomock.Any(), conversation.ID).Return(nil).Times(1)
		}

		t.Run("Should force conclude and stop", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingConversationBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ProcessConversation(context.Background(), participant, conversation)

			assert.NoError(t, err)
		})
	})

	t.Run("Given the participant is human", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()

		conversation := model_builders.NewConversationBuilder().Build()
		conversation.CreatedAt = time.Now()

		t.Run("Should return before any node invocation", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingConversationBeforeEach(t)
			suite.mockParticipantRoleResolver.EXPECT().
				Resolve(gomock.Any(), participant).
				Return(lo.ToPtr(enums.ParticipantRoleHuman), nil).
				Times(1)

			err := suite.sut.ProcessConversation(context.Background(), participant, conversation)

			assert.NoError(t, err)
		})
	})

	t.Run("Given the bridge does not exist", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		conversation.CreatedAt = time.Now()

		initSuite := func(suite *WhenProcessingConversationTestingSuite) {
			suite.mockParticipantRoleResolver.EXPECT().
				Resolve(gomock.Any(), participant).
				Return(lo.ToPtr(enums.ParticipantRoleClone), nil).
				Times(1)
			suite.mockBridgeRepository.EXPECT().GetById(gomock.Any(), conversation.BridgeID).Return(nil, nil).Times(1)
			suite.mockConversationConclude.EXPECT().Conclude(gomock.Any(), conversation.ID).Return(nil).Times(1)
		}

		t.Run("Should conclude the orphan conversation and stop", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingConversationBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ProcessConversation(context.Background(), participant, conversation)

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
			suite.mockParticipantRoleResolver.EXPECT().
				Resolve(gomock.Any(), participant).
				Return(lo.ToPtr(enums.ParticipantRoleClone), nil).
				Times(1)
			suite.mockBridgeRepository.EXPECT().
				GetById(gomock.Any(), conversation.BridgeID).
				Return(bridge, nil).
				Times(1)
			suite.mockConversationMessageRepository.EXPECT().
				ListProcessableForParticipantAndConversation(gomock.Any(), participant, gomock.Any(), conversation.ID).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should return without invoking the agent node", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingConversationBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ProcessConversation(context.Background(), participant, conversation)

			assert.NoError(t, err)
		})
	})

	t.Run("Given the agent node answers the turn", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build() // role = clone by default
		conversation := model_builders.NewConversationBuilder().Build()
		conversation.CreatedAt = time.Now()
		bridge := model_builders.NewBridgeBuilder().WithID(conversation.BridgeID).Build()

		messages := []*models.ConversationMessage{
			model_builders.NewConversationMessageBuilder().WithConversationID(conversation.ID).Build(),
		}
		participants := []*models.Participant{participant}
		card := remoteAgentCardForTest()
		turn := agentTurnRequestForTest()

		initSuite := func(suite *WhenProcessingConversationTestingSuite) {
			suite.mockParticipantRoleResolver.EXPECT().
				Resolve(gomock.Any(), participant).
				Return(lo.ToPtr(enums.ParticipantRoleClone), nil).
				Times(1)
			suite.mockBridgeRepository.EXPECT().
				GetById(gomock.Any(), conversation.BridgeID).
				Return(bridge, nil).
				Times(1)
			suite.mockConversationMessageRepository.EXPECT().
				ListProcessableForParticipantAndConversation(gomock.Any(), participant, gomock.Any(), conversation.ID).
				Return(messages, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				ListByConversationId(gomock.Any(), conversation.ID).
				Return(participants, nil).
				Times(1)
			suite.mockRemoteAgentCardResolver.EXPECT().
				ResolveCard(gomock.Any(), *participant.EnsDomain).
				Return(card, nil).
				Times(1)
			suite.mockConversationPromptBuilder.EXPECT().
				BuildInput(gomock.Any(), string(bridge.Kind), participant, gomock.Any(), conversation, messages, participants).
				Return(turn, nil).
				Times(1)
			suite.mockRemoteAgentCompletionsRequester.EXPECT().
				RequestCompletions(gomock.Any(), card, turn).
				Return([]string{"only reply"}, nil).
				Times(1)
			suite.mockConversationRepository.EXPECT().
				GetById(gomock.Any(), conversation.ID).
				Return(conversation, nil).
				Times(1)
			suite.mockBridgeCallbackRequester.EXPECT().
				PostCompletion(gomock.Any(), bridge, conversation, participant, "only reply").
				Return(nil).
				Times(1)
		}

		t.Run("Should post the completion to the bridge and leave the conversation active", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingConversationBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ProcessConversation(context.Background(), participant, conversation)

			assert.NoError(t, err)
		})
	})

	t.Run("Given the agent node requires payment", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		conversation.CreatedAt = time.Now()
		bridge := model_builders.NewBridgeBuilder().WithID(conversation.BridgeID).Build()

		messages := []*models.ConversationMessage{
			model_builders.NewConversationMessageBuilder().WithConversationID(conversation.ID).Build(),
		}
		participants := []*models.Participant{participant}
		card := remoteAgentCardForTest()
		turn := agentTurnRequestForTest()

		initSuite := func(suite *WhenProcessingConversationTestingSuite) {
			suite.mockParticipantRoleResolver.EXPECT().
				Resolve(gomock.Any(), participant).
				Return(lo.ToPtr(enums.ParticipantRoleClone), nil).
				Times(1)
			suite.mockBridgeRepository.EXPECT().
				GetById(gomock.Any(), conversation.BridgeID).
				Return(bridge, nil).
				Times(1)
			suite.mockConversationMessageRepository.EXPECT().
				ListProcessableForParticipantAndConversation(gomock.Any(), participant, gomock.Any(), conversation.ID).
				Return(messages, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				ListByConversationId(gomock.Any(), conversation.ID).
				Return(participants, nil).
				Times(1)
			suite.mockRemoteAgentCardResolver.EXPECT().
				ResolveCard(gomock.Any(), *participant.EnsDomain).
				Return(card, nil).
				Times(1)
			suite.mockConversationPromptBuilder.EXPECT().
				BuildInput(gomock.Any(), string(bridge.Kind), participant, gomock.Any(), conversation, messages, participants).
				Return(turn, nil).
				Times(1)
			suite.mockRemoteAgentCompletionsRequester.EXPECT().
				RequestCompletions(gomock.Any(), card, turn).
				Return(nil, cb_errors.NewPaymentRequiredError(*participant.EnsDomain, "payload")).
				Times(1)
		}

		t.Run("Should skip the turn without error and post nothing", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingConversationBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ProcessConversation(context.Background(), participant, conversation)

			assert.NoError(t, err)
		})
	})

	t.Run("Given endpoint resolution fails", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		conversation.CreatedAt = time.Now()
		bridge := model_builders.NewBridgeBuilder().WithID(conversation.BridgeID).Build()

		messages := []*models.ConversationMessage{
			model_builders.NewConversationMessageBuilder().WithConversationID(conversation.ID).Build(),
		}

		initSuite := func(suite *WhenProcessingConversationTestingSuite) {
			suite.mockParticipantRoleResolver.EXPECT().
				Resolve(gomock.Any(), participant).
				Return(lo.ToPtr(enums.ParticipantRoleClone), nil).
				Times(1)
			suite.mockBridgeRepository.EXPECT().
				GetById(gomock.Any(), conversation.BridgeID).
				Return(bridge, nil).
				Times(1)
			suite.mockConversationMessageRepository.EXPECT().
				ListProcessableForParticipantAndConversation(gomock.Any(), participant, gomock.Any(), conversation.ID).
				Return(messages, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				ListByConversationId(gomock.Any(), conversation.ID).
				Return([]*models.Participant{participant}, nil).
				Times(1)
			suite.mockRemoteAgentCardResolver.EXPECT().
				ResolveCard(gomock.Any(), *participant.EnsDomain).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingConversationBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ProcessConversation(context.Background(), participant, conversation)

			assert.Error(t, err)
		})
	})

	t.Run("Given a facilitator completion carries the conclude marker", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		conversation.CreatedAt = time.Now()
		bridge := model_builders.NewBridgeBuilder().WithID(conversation.BridgeID).Build()

		messages := []*models.ConversationMessage{
			model_builders.NewConversationMessageBuilder().WithConversationID(conversation.ID).Build(),
		}
		participants := []*models.Participant{participant}
		card := remoteAgentCardForTest()
		turn := agentTurnRequestForTest()
		concludeTurn := agentTurnRequestForTest()

		initSuite := func(suite *WhenProcessingConversationTestingSuite) {
			suite.mockParticipantRoleResolver.EXPECT().
				Resolve(gomock.Any(), participant).
				Return(lo.ToPtr(enums.ParticipantRoleFacilitator), nil).
				Times(1)
			suite.mockBridgeRepository.EXPECT().
				GetById(gomock.Any(), conversation.BridgeID).
				Return(bridge, nil).
				Times(1)
			suite.mockConversationMessageRepository.EXPECT().
				ListProcessableForParticipantAndConversation(gomock.Any(), participant, gomock.Any(), conversation.ID).
				Return(messages, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				ListByConversationId(gomock.Any(), conversation.ID).
				Return(participants, nil).
				Times(1)
			suite.mockRemoteAgentCardResolver.EXPECT().
				ResolveCard(gomock.Any(), *participant.EnsDomain).
				Return(card, nil).
				Times(1)
			suite.mockConversationPromptBuilder.EXPECT().
				BuildInput(gomock.Any(), string(bridge.Kind), participant, gomock.Any(), conversation, messages, participants).
				Return(turn, nil).
				Times(1)
			suite.mockRemoteAgentCompletionsRequester.EXPECT().
				RequestCompletions(gomock.Any(), card, turn).
				Return([]string{"we are done :end:"}, nil).
				Times(1)
			suite.mockConversationPromptBuilder.EXPECT().
				BuildConcludeInput(gomock.Any(),
					gomock.Any(),
					[]string{"we are done :end:"},
					participant,
					gomock.Any(),
					conversation,
					messages,
					participants,
				).
				Return(concludeTurn, nil).
				Times(1)
			suite.mockRemoteAgentCompletionsRequester.EXPECT().
				RequestCompletions(gomock.Any(), card, concludeTurn).
				Return([]string{"final summary"}, nil).
				Times(1)
			suite.mockConversationRepository.EXPECT().
				GetById(gomock.Any(), conversation.ID).
				Return(conversation, nil).
				Times(1)
			suite.mockBridgeCallbackRequester.EXPECT().
				PostCompletion(gomock.Any(), bridge, conversation, participant, "we are done :end:").
				Return(nil).
				Times(1)
			suite.mockBridgeCallbackRequester.EXPECT().
				PostCompletion(gomock.Any(), bridge, conversation, participant, "final summary").
				Return(nil).
				Times(1)
			suite.mockConversationConclude.EXPECT().Conclude(gomock.Any(), conversation.ID).Return(nil).Times(1)
		}

		t.Run("Should run the conclude turn, post both completions and conclude", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingConversationBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ProcessConversation(context.Background(), participant, conversation)

			assert.NoError(t, err)
		})
	})

	t.Run("Given the latest inbound message carries the conclude marker", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		conversation.CreatedAt = time.Now()
		bridge := model_builders.NewBridgeBuilder().WithID(conversation.BridgeID).Build()

		messages := []*models.ConversationMessage{
			model_builders.NewConversationMessageBuilder().
				WithConversationID(conversation.ID).
				WithContent("thanks all, :end:").
				Build(),
		}
		participants := []*models.Participant{participant}
		card := remoteAgentCardForTest()
		turn := agentTurnRequestForTest()
		concludeTurn := agentTurnRequestForTest()

		initSuite := func(suite *WhenProcessingConversationTestingSuite) {
			suite.mockParticipantRoleResolver.EXPECT().
				Resolve(gomock.Any(), participant).
				Return(lo.ToPtr(enums.ParticipantRoleFacilitator), nil).
				Times(1)
			suite.mockBridgeRepository.EXPECT().
				GetById(gomock.Any(), conversation.BridgeID).
				Return(bridge, nil).
				Times(1)
			suite.mockConversationMessageRepository.EXPECT().
				ListProcessableForParticipantAndConversation(gomock.Any(), participant, gomock.Any(), conversation.ID).
				Return(messages, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				ListByConversationId(gomock.Any(), conversation.ID).
				Return(participants, nil).
				Times(1)
			suite.mockRemoteAgentCardResolver.EXPECT().
				ResolveCard(gomock.Any(), *participant.EnsDomain).
				Return(card, nil).
				Times(1)
			suite.mockConversationPromptBuilder.EXPECT().
				BuildInput(gomock.Any(), string(bridge.Kind), participant, gomock.Any(), conversation, messages, participants).
				Return(turn, nil).
				Times(1)
			suite.mockRemoteAgentCompletionsRequester.EXPECT().
				RequestCompletions(gomock.Any(), card, turn).
				Return([]string{"wrapping up"}, nil).
				Times(1)
			suite.mockConversationPromptBuilder.EXPECT().
				BuildConcludeInput(gomock.Any(),
					gomock.Any(),
					[]string{"wrapping up"},
					participant,
					gomock.Any(),
					conversation,
					messages,
					participants,
				).
				Return(concludeTurn, nil).
				Times(1)
			suite.mockRemoteAgentCompletionsRequester.EXPECT().
				RequestCompletions(gomock.Any(), card, concludeTurn).
				Return([]string{"final summary"}, nil).
				Times(1)
			suite.mockConversationRepository.EXPECT().
				GetById(gomock.Any(), conversation.ID).
				Return(conversation, nil).
				Times(1)
			suite.mockBridgeCallbackRequester.EXPECT().
				PostCompletion(gomock.Any(), bridge, conversation, participant, "wrapping up").
				Return(nil).
				Times(1)
			suite.mockBridgeCallbackRequester.EXPECT().
				PostCompletion(gomock.Any(), bridge, conversation, participant, "final summary").
				Return(nil).
				Times(1)
			suite.mockConversationConclude.EXPECT().Conclude(gomock.Any(), conversation.ID).Return(nil).Times(1)
		}

		t.Run("Should conclude even though the facilitator completion has no marker", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingConversationBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ProcessConversation(context.Background(), participant, conversation)

			assert.NoError(t, err)
		})
	})

	t.Run("Given the conversation got concluded between the turn and reload", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		conversation.CreatedAt = time.Now()
		bridge := model_builders.NewBridgeBuilder().WithID(conversation.BridgeID).Build()

		messages := []*models.ConversationMessage{
			model_builders.NewConversationMessageBuilder().WithConversationID(conversation.ID).Build(),
		}
		participants := []*models.Participant{participant}
		card := remoteAgentCardForTest()
		turn := agentTurnRequestForTest()

		reloaded := model_builders.NewConversationBuilder().
			WithID(conversation.ID).
			WithStatus(enums.ConversationStatusConcluded).
			Build()

		initSuite := func(suite *WhenProcessingConversationTestingSuite) {
			suite.mockParticipantRoleResolver.EXPECT().
				Resolve(gomock.Any(), participant).
				Return(lo.ToPtr(enums.ParticipantRoleClone), nil).
				Times(1)
			suite.mockBridgeRepository.EXPECT().
				GetById(gomock.Any(), conversation.BridgeID).
				Return(bridge, nil).
				Times(1)
			suite.mockConversationMessageRepository.EXPECT().
				ListProcessableForParticipantAndConversation(gomock.Any(), participant, gomock.Any(), conversation.ID).
				Return(messages, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				ListByConversationId(gomock.Any(), conversation.ID).
				Return(participants, nil).
				Times(1)
			suite.mockRemoteAgentCardResolver.EXPECT().
				ResolveCard(gomock.Any(), *participant.EnsDomain).
				Return(card, nil).
				Times(1)
			suite.mockConversationPromptBuilder.EXPECT().
				BuildInput(gomock.Any(), string(bridge.Kind), participant, gomock.Any(), conversation, messages, participants).
				Return(turn, nil).
				Times(1)
			suite.mockRemoteAgentCompletionsRequester.EXPECT().
				RequestCompletions(gomock.Any(), card, turn).
				Return([]string{"a"}, nil).
				Times(1)
			suite.mockConversationRepository.EXPECT().
				GetById(gomock.Any(), conversation.ID).
				Return(reloaded, nil).
				Times(1)
		}

		t.Run("Should skip posting completions and not error", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingConversationBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ProcessConversation(context.Background(), participant, conversation)

			assert.NoError(t, err)
		})
	})
}
