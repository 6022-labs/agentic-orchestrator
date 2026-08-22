package services_test

import (
	"context"

	stderrors "errors"
	"net/http"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/services"
	clients_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_http_agent_swarm_conversation_broker_mocks/clients_mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingAgentPromptTemplateTestingSuite struct {
	sut *services.ConversationBrokerRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenGettingAgentPromptTemplateBeforeEach(t *testing.T) *WhenGettingAgentPromptTemplateTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	return &WhenGettingAgentPromptTemplateTestingSuite{
		sut:        services.NewConversationBrokerRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenGettingAgentPromptTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Given the prompt template exists", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite) {
			suite.mockClient.EXPECT().
				GetAgentsCollectionAddressCollectionTokenIdPromptTemplatesKeyWithResponse(
					gomock.Any(),
					"0xabc",
					"7",
					"clone-system",
				).
				Return(&clients.GetAgentsCollectionAddressCollectionTokenIdPromptTemplatesKeyResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusOK},
					JSON200: &clients.ResponsesAgentPromptTemplateResponse{
						CollectionAddress: lo.ToPtr("0xabc"),
						CollectionTokenId: lo.ToPtr(7),
						Key:               lo.ToPtr(clients.EnumsAgentPromptTemplateKey("clone-system")),
						Content:           lo.ToPtr("You are {{name}}"),
					},
				}, nil).
				Times(1)
		}

		t.Run("Should return the mapped template", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.GetAgentPromptTemplate(context.Background(), "0xabc", 7, "clone-system")

			assert.NoError(t, err)
			if assert.NotNil(t, result) {
				assert.Equal(t, "0xabc", result.CollectionAddress)
				assert.Equal(t, uint64(7), result.CollectionTokenId)
				assert.Equal(t, "clone-system", result.Key)
				assert.Equal(t, "You are {{name}}", result.Content)
			}
		})
	})

	t.Run("Given the prompt template does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite) {
			suite.mockClient.EXPECT().
				GetAgentsCollectionAddressCollectionTokenIdPromptTemplatesKeyWithResponse(
					gomock.Any(),
					"0xabc",
					"7",
					"clone-system",
				).
				Return(&clients.GetAgentsCollectionAddressCollectionTokenIdPromptTemplatesKeyResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusNotFound},
				}, nil).
				Times(1)
		}

		t.Run("Should return a NotFoundError", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.GetAgentPromptTemplate(context.Background(), "0xabc", 7, "clone-system")

			var target *errors.NotFoundError
			assert.ErrorAs(t, err, &target)
			assert.Nil(t, result)
		})
	})

	t.Run("Given the broker returns an unexpected status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite) {
			suite.mockClient.EXPECT().
				GetAgentsCollectionAddressCollectionTokenIdPromptTemplatesKeyWithResponse(
					gomock.Any(),
					"0xabc",
					"7",
					"clone-system",
				).
				Return(&clients.GetAgentsCollectionAddressCollectionTokenIdPromptTemplatesKeyResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusInternalServerError},
					Body:         []byte("boom"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.GetAgentPromptTemplate(context.Background(), "0xabc", 7, "clone-system")

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite) {
			suite.mockClient.EXPECT().
				GetAgentsCollectionAddressCollectionTokenIdPromptTemplatesKeyWithResponse(
					gomock.Any(),
					"0xabc",
					"7",
					"clone-system",
				).
				Return(nil, stderrors.New("broker down")).
				Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.GetAgentPromptTemplate(context.Background(), "0xabc", 7, "clone-system")

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
}
