package services_test

import (
	"errors"
	"testing"

	domain_models "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	http_compute_models "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_agent_reasoning_agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_agent_reasoning_agent_compute/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/conversation_broker_mocks/http_agent_reasoning_agent_compute_mocks/clients_mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenRequestingAgentCompletionsTestingSuite struct {
	sut *services.AgentCompletionsRequester

	mockClient *clients_mocks.MockAgentReasoningAgentComputeClientInterface
}

func WhenRequestingAgentCompletionsBeforeEach(t *testing.T) *WhenRequestingAgentCompletionsTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockAgentReasoningAgentComputeClientInterface(mockController)

	sut := services.NewAgentCompletionsRequester(mockClient)

	return &WhenRequestingAgentCompletionsTestingSuite{
		sut: sut,

		mockClient: mockClient,
	}
}

func TestWhenRequestingAgentCompletions(t *testing.T) {
	t.Parallel()

	t.Run("Given the agent participant is not minted", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingAgentCompletionsBeforeEach(t)
			agent := model_builders.NewNonMintedParticipantBuilder().Build()
			conversation := model_builders.NewConversationBuilder().Build()

			_, _, err := suite.sut.Completions("slack", agent, conversation, nil, nil, false)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "must be minted")
		})
	})

	t.Run("Given the underlying client fails", func(t *testing.T) {
		t.Parallel()

		agent := model_builders.NewParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()

		initSuite := func(suite *WhenRequestingAgentCompletionsTestingSuite) {
			suite.mockClient.EXPECT().
				Completions(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingAgentCompletionsBeforeEach(t)
			initSuite(suite)

			_, _, err := suite.sut.Completions("slack", agent, conversation, nil, []*domain_models.Participant{agent}, false)

			assert.Error(t, err)
		})
	})

	t.Run("Given a conversation with messages and mixed participants", func(t *testing.T) {
		t.Parallel()

		agent := model_builders.NewParticipantBuilder().Build()
		otherMintedAgent := model_builders.NewParticipantBuilder().Build()
		human := model_builders.NewNonMintedParticipantBuilder().Build()

		descriptionValue := "engineering channel"
		conversation := model_builders.NewConversationBuilder().Build()
		conversation.Description = &descriptionValue

		humanMessage := model_builders.NewConversationMessageBuilder().
			WithConversationID(conversation.ID).
			WithParticipantID(human.ID).
			WithContent("hello team").
			Build()
		agentMessage := model_builders.NewConversationMessageBuilder().
			WithConversationID(conversation.ID).
			WithParticipantID(otherMintedAgent.ID).
			WithContent("hi pierre").
			Build()
		messages := []*domain_models.ConversationMessage{humanMessage, agentMessage}
		participants := []*domain_models.Participant{agent, otherMintedAgent, human}

		isFinalReturn := true
		responsesReturn := []*string{lo.ToPtr("first"), lo.ToPtr("second")}

		initSuite := func(suite *WhenRequestingAgentCompletionsTestingSuite) {
			suite.mockClient.EXPECT().
				Completions(*agent.CollectionAddress, gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ string, tokenId uint64, body *http_compute_models.ExternalAgentComputeRequest) (*http_compute_models.ExternalAgentComputeResponses, error) {
					assert.Equal(t, "slack", body.Platform)
					assert.True(t, body.ForceConclusion)
					assert.Equal(t, conversation.ID.String(), body.ConversationContext.Id)
					assert.Len(t, body.ConversationContext.Agents, 2)
					assert.Equal(t, descriptionValue, body.PromptParameters["ChannelDescription"])

					channelMembers, ok := body.PromptParameters["ChannelMembers"].([]string)
					assert.True(t, ok)
					assert.Contains(t, channelMembers, *agent.EnsDomain)
					assert.Contains(t, channelMembers, *otherMintedAgent.EnsDomain)
					assert.Contains(t, channelMembers, *human.DisplayName)

					assert.Len(t, body.Inputs, 2)
					assert.Nil(t, body.Inputs[0].Author)
					assert.Equal(t, "hello team", body.Inputs[0].Content)
					assert.NotNil(t, body.Inputs[1].Author)
					assert.Equal(t, *otherMintedAgent.EnsDomain, *body.Inputs[1].Author)

					return &http_compute_models.ExternalAgentComputeResponses{
						IsFinal:   isFinalReturn,
						Responses: responsesReturn,
					}, nil
				}).
				Times(1)
		}

		t.Run("Should map the request and return the responses", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingAgentCompletionsBeforeEach(t)
			initSuite(suite)

			responses, isFinal, err := suite.sut.Completions("slack", agent, conversation, messages, participants, true)

			assert.NoError(t, err)
			assert.True(t, isFinal)
			assert.Equal(t, responsesReturn, responses)
		})
	})
}