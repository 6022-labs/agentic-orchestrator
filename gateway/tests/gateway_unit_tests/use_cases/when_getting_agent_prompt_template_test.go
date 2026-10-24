package use_cases_test

import (
	"context"

	"strings"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingAgentPromptTemplateTestingSuite struct {
	sut *use_cases.GetAgentPromptTemplate

	mockAgentAccessVerifier         *interfaces_mocks.MockAgentAccessVerifierInterface
	mockConversationBrokerRequester *interfaces_mocks.MockConversationBrokerRequesterInterface
}

func WhenGettingAgentPromptTemplateBeforeEach(t *testing.T) *WhenGettingAgentPromptTemplateTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentAccessVerifier := interfaces_mocks.NewMockAgentAccessVerifierInterface(mockController)
	mockConversationBrokerRequester := interfaces_mocks.NewMockConversationBrokerRequesterInterface(mockController)

	sut := use_cases.NewGetAgentPromptTemplate(mockAgentAccessVerifier, mockConversationBrokerRequester)

	return &WhenGettingAgentPromptTemplateTestingSuite{
		sut: sut,

		mockAgentAccessVerifier:         mockAgentAccessVerifier,
		mockConversationBrokerRequester: mockConversationBrokerRequester,
	}
}

func generateValidGetAgentPromptTemplateRequest() *requests.GetAgentPromptTemplateRequest {
	return &requests.GetAgentPromptTemplateRequest{
		AgentCollectionAddress: common_tests.RandomEthAddress().Hex(),
		AgentCollectionTokenId: "1",
		Key:                    "clone-system",
	}
}

func TestWhenGettingAgentPromptTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid because token id is not a number", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			request := generateValidGetAgentPromptTemplateRequest()
			request.AgentCollectionTokenId = "abc"

			_, err := suite.sut.Get(context.Background(), request)

			var validationErr *cb_errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "agentCollectionTokenId", validationErr.Field)
		})
	})

	t.Run("Given no caller id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an unauthorized error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)

			_, err := suite.sut.Get(context.Background(), generateValidGetAgentPromptTemplateRequest())

			var unauthorizedErr *cb_errors.UnauthorizedError
			assert.ErrorAs(t, err, &unauthorizedErr)
			assert.Equal(t, "missing_caller_id", unauthorizedErr.Code)
		})
	})

	t.Run("Given the caller has no access to the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(
			suite *WhenGettingAgentPromptTemplateTestingSuite,
			request *requests.GetAgentPromptTemplateRequest,
			callerId string,
		) {
			suite.mockAgentAccessVerifier.EXPECT().
				EnsureAccess(gomock.Any(), callerId, strings.ToLower(request.AgentCollectionAddress), uint64(1)).
				Return(cb_errors.NewForbiddenError(services.AgentAccessForbiddenCode, "forbidden")).
				Times(1)
		}

		t.Run("Should return the forbidden error without reading the template", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			callerId := strings.ToLower(common_tests.RandomEthAddress().Hex())
			ctx := helpers.ContextWithCallerId(context.Background(), callerId)
			request := generateValidGetAgentPromptTemplateRequest()
			initSuite(suite, request, callerId)

			_, err := suite.sut.Get(ctx, request)

			var forbiddenErr *cb_errors.ForbiddenError
			assert.ErrorAs(t, err, &forbiddenErr)
			assert.Equal(t, services.AgentAccessForbiddenCode, forbiddenErr.Code)
		})
	})

	t.Run("Given the broker fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite) {
			suite.mockAgentAccessVerifier.EXPECT().
				EnsureAccess(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).
				Times(1)
			suite.mockConversationBrokerRequester.EXPECT().
				GetAgentPromptTemplate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			ctx := helpers.ContextWithCallerId(
				context.Background(),
				strings.ToLower(common_tests.RandomEthAddress().Hex()),
			)
			initSuite(suite)

			_, err := suite.sut.Get(ctx, generateValidGetAgentPromptTemplateRequest())

			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Given the caller has access to the agent and the broker returns the template", func(t *testing.T) {
		t.Parallel()

		initSuite := func(
			suite *WhenGettingAgentPromptTemplateTestingSuite,
			request *requests.GetAgentPromptTemplateRequest,
			callerId string,
		) {
			suite.mockAgentAccessVerifier.EXPECT().
				EnsureAccess(gomock.Any(), callerId, strings.ToLower(request.AgentCollectionAddress), uint64(1)).
				Return(nil).
				Times(1)
			suite.mockConversationBrokerRequester.EXPECT().
				GetAgentPromptTemplate(gomock.Any(), strings.ToLower(request.AgentCollectionAddress), uint64(1), "clone-system").
				Return(&responses.AgentPromptTemplateResponse{Key: "clone-system", Content: "tpl"}, nil).
				Times(1)
		}

		t.Run("Should return the template", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			callerId := strings.ToLower(common_tests.RandomEthAddress().Hex())
			ctx := helpers.ContextWithCallerId(context.Background(), callerId)
			request := generateValidGetAgentPromptTemplateRequest()
			initSuite(suite, request, callerId)

			response, err := suite.sut.Get(ctx, request)

			assert.NoError(t, err)
			assert.Equal(t, "tpl", response.Content)
		})
	})
}
