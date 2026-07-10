package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	http_node_models "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_remote_agent/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_remote_agent/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/conversation_broker_http_remote_agent_mocks/clients_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenRequestingRemoteAgentCompletionsTestingSuite struct {
	sut *services.RemoteAgentCompletionsRequester

	mockRemoteAgentClient *clients_mocks.MockRemoteAgentClientInterface
}

func WhenRequestingRemoteAgentCompletionsBeforeEach(t *testing.T) *WhenRequestingRemoteAgentCompletionsTestingSuite {
	mockController := gomock.NewController(t)

	mockRemoteAgentClient := clients_mocks.NewMockRemoteAgentClientInterface(mockController)

	sut := services.NewRemoteAgentCompletionsRequester(mockRemoteAgentClient)

	return &WhenRequestingRemoteAgentCompletionsTestingSuite{
		sut: sut,

		mockRemoteAgentClient: mockRemoteAgentClient,
	}
}

func turnForTest() *models.RemoteAgentInput {
	return &models.RemoteAgentInput{
		Instructions: "instructions",
		Input: []*models.RemoteAgentInputMessage{
			{Role: enums.RemoteAgentInputRoleUser, Content: "hello"},
			{Role: enums.RemoteAgentInputRoleAssistant, Content: "hi"},
		},
		Metadata: map[string]string{"ensDomain": "agent.eth"},
	}
}

func TestWhenRequestingRemoteAgentCompletions(t *testing.T) {
	t.Parallel()

	t.Run("Given a responses endpoint", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{Protocol: "6022", Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"}}

		initSuite := func(suite *WhenRequestingRemoteAgentCompletionsTestingSuite, turn *models.RemoteAgentInput) {
			suite.mockRemoteAgentClient.EXPECT().
				CreateResponse(card.Endpoints.Responses, "agent.eth", gomock.Any()).
				DoAndReturn(func(_, _ string, request *http_node_models.ExternalCreateResponseRequest) (*http_node_models.ExternalResponse, error) {
					assert.Equal(t, turn.Instructions, request.Instructions)
					assert.Len(t, request.Input, 2)
					assert.Equal(t, "message", request.Input[0].Type)
					assert.Equal(t, "user", request.Input[0].Role)
					assert.Equal(t, "assistant", request.Input[1].Role)
					return &http_node_models.ExternalResponse{
						Status: "completed",
						Output: []*http_node_models.ExternalResponseOutput{
							{Type: "reasoning"},
							{Type: "message", Role: "assistant", Content: []*http_node_models.ExternalResponseOutputContent{{Type: "output_text", Text: "answer"}}},
						},
					}, nil
				}).
				Times(1)
		}

		t.Run("Should map the turn and extract message output texts", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
			turn := turnForTest()
			initSuite(suite, turn)

			completions, err := suite.sut.RequestCompletions(card, turn)

			assert.NoError(t, err)
			assert.Equal(t, []string{"answer"}, completions)
		})
	})

	t.Run("Given the response carries an error payload", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{Protocol: "6022", Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"}}

		initSuite := func(suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			suite.mockRemoteAgentClient.EXPECT().
				CreateResponse(card.Endpoints.Responses, "agent.eth", gomock.Any()).
				Return(&http_node_models.ExternalResponse{Error: &http_node_models.ExternalResponseError{Code: "server_error", Message: "boom"}}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.RequestCompletions(card, turnForTest())

			assert.Error(t, err)
		})
	})

	t.Run("Given a chat completions endpoint", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{Protocol: "6022", Endpoints: models.RemoteAgentCardEndpoints{ChatCompletions: "https://node.example/api/chat"}}

		initSuite := func(suite *WhenRequestingRemoteAgentCompletionsTestingSuite, turn *models.RemoteAgentInput) {
			suite.mockRemoteAgentClient.EXPECT().
				CreateChatCompletion(card.Endpoints.ChatCompletions, "agent.eth", gomock.Any()).
				DoAndReturn(func(_, _ string, request *http_node_models.ExternalChatCompletionRequest) (*http_node_models.ExternalChatCompletionResponse, error) {
					assert.Len(t, request.Messages, 3)
					assert.Equal(t, "system", request.Messages[0].Role)
					assert.Equal(t, turn.Instructions, request.Messages[0].Content)
					return &http_node_models.ExternalChatCompletionResponse{
						Choices: []*http_node_models.ExternalChatCompletionChoice{
							{Message: &http_node_models.ExternalChatCompletionMessage{Role: "assistant", Content: "chat answer"}},
						},
					}, nil
				}).
				Times(1)
		}

		t.Run("Should prepend instructions as a system message and extract choices", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
			turn := turnForTest()
			initSuite(suite, turn)

			completions, err := suite.sut.RequestCompletions(card, turn)

			assert.NoError(t, err)
			assert.Equal(t, []string{"chat answer"}, completions)
		})
	})

	t.Run("Given the response carries no message output", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{Protocol: "6022", Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"}}

		initSuite := func(suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			suite.mockRemoteAgentClient.EXPECT().
				CreateResponse(card.Endpoints.Responses, "agent.eth", gomock.Any()).
				Return(&http_node_models.ExternalResponse{Status: "completed"}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.RequestCompletions(card, turnForTest())

			assert.Error(t, err)
		})
	})
}
