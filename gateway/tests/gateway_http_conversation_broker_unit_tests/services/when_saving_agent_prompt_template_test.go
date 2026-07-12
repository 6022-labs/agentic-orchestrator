package services_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_conversation_broker/services"
	clients_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_http_conversation_broker_mocks/clients_mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenSavingAgentPromptTemplateTestingSuite struct {
	sut *services.ConversationBrokerRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenSavingAgentPromptTemplateBeforeEach(t *testing.T) *WhenSavingAgentPromptTemplateTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	return &WhenSavingAgentPromptTemplateTestingSuite{
		sut:        services.NewConversationBrokerRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenSavingAgentPromptTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Given the broker saves the template", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite) *string {
			var sentBody string

			suite.mockClient.EXPECT().PostAgentsCollectionAddressCollectionTokenIdPromptTemplatesKeyWithBodyWithResponse(gomock.Any(), "0xabc", "7", "clone-system", "application/json", gomock.Any()).DoAndReturn(
				func(_ any, _, _, _, _ string, body io.Reader, _ ...clients.RequestEditorFn) (*clients.PostAgentsCollectionAddressCollectionTokenIdPromptTemplatesKeyResponse, error) {
					raw, _ := io.ReadAll(body)
					sentBody = string(raw)
					return &clients.PostAgentsCollectionAddressCollectionTokenIdPromptTemplatesKeyResponse{
						HTTPResponse: &http.Response{StatusCode: http.StatusOK},
						JSON200: &clients.ResponsesAgentPromptTemplateResponse{
							CollectionAddress: lo.ToPtr("0xabc"),
							CollectionTokenId: lo.ToPtr(7),
							Key:               lo.ToPtr(clients.EnumsAgentPromptTemplateKey("clone-system")),
							Content:           lo.ToPtr("You are {{name}}"),
						},
					}, nil
				}).Times(1)

			return &sentBody
		}

		t.Run("Should post the content and return the mapped template", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			sentBody := initSuite(suite)

			result, err := suite.sut.SaveAgentPromptTemplate("0xabc", 7, "clone-system", "You are {{name}}")

			assert.NoError(t, err)
			assert.JSONEq(t, `{"content":"You are {{name}}"}`, *sentBody)
			if assert.NotNil(t, result) {
				assert.Equal(t, "clone-system", result.Key)
				assert.Equal(t, "You are {{name}}", result.Content)
			}
		})
	})

	t.Run("Given the broker fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite) {
			suite.mockClient.EXPECT().PostAgentsCollectionAddressCollectionTokenIdPromptTemplatesKeyWithBodyWithResponse(gomock.Any(), "0xabc", "7", "clone-system", "application/json", gomock.Any()).Return(&clients.PostAgentsCollectionAddressCollectionTokenIdPromptTemplatesKeyResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusInternalServerError},
				Body:         []byte("boom"),
			}, nil).Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.SaveAgentPromptTemplate("0xabc", 7, "clone-system", "You are {{name}}")

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
}
