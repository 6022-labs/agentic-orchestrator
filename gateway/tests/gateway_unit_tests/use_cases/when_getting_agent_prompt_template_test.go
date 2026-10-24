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

type WhenGettingAgentPromptTemplateTestingSuite struct {
	sut *use_cases.GetAgentPromptTemplate

	mockConversationBrokerRequester *interfaces_mocks.MockConversationBrokerRequesterInterface
}

func WhenGettingAgentPromptTemplateBeforeEach(t *testing.T) *WhenGettingAgentPromptTemplateTestingSuite {
	mockController := gomock.NewController(t)

	mockConversationBrokerRequester := interfaces_mocks.NewMockConversationBrokerRequesterInterface(mockController)

	sut := use_cases.NewGetAgentPromptTemplate(mockConversationBrokerRequester)

	return &WhenGettingAgentPromptTemplateTestingSuite{
		sut: sut,

		mockConversationBrokerRequester: mockConversationBrokerRequester,
	}
}

func generateValidGetAgentPromptTemplateRequest() *requests.GetAgentPromptTemplateRequest {
	return &requests.GetAgentPromptTemplateRequest{
		AgentCollectionAddress: "0x39DD75a025145770a8b8278Cc8BA705116443Ea7",
		AgentCollectionTokenId: "1",
		Key:                    "clone-system",
	}
}

func TestWhenGettingAgentPromptTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because token id is not a number", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			request := generateValidGetAgentPromptTemplateRequest()
			request.AgentCollectionTokenId = "abc"

			_, err := suite.sut.Get(context.Background(), request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given the broker fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite) {
			suite.mockConversationBrokerRequester.EXPECT().
				GetAgentPromptTemplate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Get(context.Background(), generateValidGetAgentPromptTemplateRequest())

			assert.Error(t, err)
		})
	})

	t.Run("Given the broker returns the template", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite) {
			suite.mockConversationBrokerRequester.EXPECT().
				GetAgentPromptTemplate(gomock.Any(), "0x39dd75a025145770a8b8278cc8ba705116443ea7", uint64(1), "clone-system").
				Return(&responses.AgentPromptTemplateResponse{Key: "clone-system", Content: "tpl"}, nil).
				Times(1)
		}

		t.Run("Should return the template", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Get(context.Background(), generateValidGetAgentPromptTemplateRequest())

			assert.NoError(t, err)
			assert.Equal(t, "tpl", response.Content)
		})
	})
}
