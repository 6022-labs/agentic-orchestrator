package use_cases_test

import (
	"context"

	"errors"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenSavingAgentPromptTemplateTestingSuite struct {
	sut *use_cases.SaveAgentPromptTemplate

	mockConversationBrokerRequester *interfaces_mocks.MockConversationBrokerRequesterInterface
}

func WhenSavingAgentPromptTemplateBeforeEach(t *testing.T) *WhenSavingAgentPromptTemplateTestingSuite {
	mockController := gomock.NewController(t)

	mockConversationBrokerRequester := interfaces_mocks.NewMockConversationBrokerRequesterInterface(mockController)

	sut := use_cases.NewSaveAgentPromptTemplate(mockConversationBrokerRequester)

	return &WhenSavingAgentPromptTemplateTestingSuite{
		sut: sut,

		mockConversationBrokerRequester: mockConversationBrokerRequester,
	}
}

func generateValidSaveAgentPromptTemplateRequest() *requests.SaveAgentPromptTemplateRequest {
	return &requests.SaveAgentPromptTemplateRequest{
		GetAgentPromptTemplateRequest: requests.GetAgentPromptTemplateRequest{
			AgentCollectionAddress: "0x39DD75a025145770a8b8278Cc8BA705116443Ea7",
			AgentCollectionTokenId: "1",
			Key:                    "clone-system",
		},
		Content: "You are an agent.",
	}
}

func TestWhenSavingAgentPromptTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because content is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			request.Content = " "

			_, err := suite.sut.Save(context.Background(), request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given the broker fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite) {
			suite.mockConversationBrokerRequester.EXPECT().
				SaveAgentPromptTemplate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Save(context.Background(), generateValidSaveAgentPromptTemplateRequest())

			assert.Error(t, err)
		})
	})

	t.Run("Given the broker saves the template", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite, request *requests.SaveAgentPromptTemplateRequest) {
			suite.mockConversationBrokerRequester.EXPECT().
				SaveAgentPromptTemplate(
					gomock.Any(),
					"0x39dd75a025145770a8b8278cc8ba705116443ea7", uint64(1), "clone-system", request.Content,
				).
				Return(&responses.AgentPromptTemplateResponse{Key: "clone-system", Content: request.Content}, nil).
				Times(1)
		}

		t.Run("Should return the saved template", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			initSuite(suite, request)

			response, err := suite.sut.Save(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, request.Content, response.Content)
		})
	})
}
