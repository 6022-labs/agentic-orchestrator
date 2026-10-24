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

type WhenSavingAgentPromptTemplateTestingSuite struct {
	sut *use_cases.SaveAgentPromptTemplate

	mockAgentAccessVerifier         *interfaces_mocks.MockAgentAccessVerifierInterface
	mockConversationBrokerRequester *interfaces_mocks.MockConversationBrokerRequesterInterface
}

func WhenSavingAgentPromptTemplateBeforeEach(t *testing.T) *WhenSavingAgentPromptTemplateTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentAccessVerifier := interfaces_mocks.NewMockAgentAccessVerifierInterface(mockController)
	mockConversationBrokerRequester := interfaces_mocks.NewMockConversationBrokerRequesterInterface(mockController)

	sut := use_cases.NewSaveAgentPromptTemplate(mockAgentAccessVerifier, mockConversationBrokerRequester)

	return &WhenSavingAgentPromptTemplateTestingSuite{
		sut: sut,

		mockAgentAccessVerifier:         mockAgentAccessVerifier,
		mockConversationBrokerRequester: mockConversationBrokerRequester,
	}
}

func generateValidSaveAgentPromptTemplateRequest() *requests.SaveAgentPromptTemplateRequest {
	return &requests.SaveAgentPromptTemplateRequest{
		GetAgentPromptTemplateRequest: requests.GetAgentPromptTemplateRequest{
			AgentCollectionAddress: common_tests.RandomEthAddress().Hex(),
			AgentCollectionTokenId: "1",
			Key:                    "clone-system",
		},
		Content: "You are an agent.",
	}
}

func TestWhenSavingAgentPromptTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid because content is empty", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			request.Content = " "

			_, err := suite.sut.Save(context.Background(), request)

			var validationErr *cb_errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "content", validationErr.Field)
		})
	})

	t.Run("Given no caller id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an unauthorized error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)

			_, err := suite.sut.Save(context.Background(), generateValidSaveAgentPromptTemplateRequest())

			var unauthorizedErr *cb_errors.UnauthorizedError
			assert.ErrorAs(t, err, &unauthorizedErr)
			assert.Equal(t, "missing_caller_id", unauthorizedErr.Code)
		})
	})

	t.Run("Given the caller has no access to the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(
			suite *WhenSavingAgentPromptTemplateTestingSuite,
			request *requests.SaveAgentPromptTemplateRequest,
			callerId string,
		) {
			suite.mockAgentAccessVerifier.EXPECT().
				EnsureAccess(gomock.Any(), callerId, strings.ToLower(request.AgentCollectionAddress), uint64(1)).
				Return(cb_errors.NewForbiddenError(services.AgentAccessForbiddenCode, "forbidden")).
				Times(1)
		}

		t.Run("Should return the forbidden error without saving", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			callerId := strings.ToLower(common_tests.RandomEthAddress().Hex())
			ctx := helpers.ContextWithCallerId(context.Background(), callerId)
			request := generateValidSaveAgentPromptTemplateRequest()
			initSuite(suite, request, callerId)

			_, err := suite.sut.Save(ctx, request)

			var forbiddenErr *cb_errors.ForbiddenError
			assert.ErrorAs(t, err, &forbiddenErr)
			assert.Equal(t, services.AgentAccessForbiddenCode, forbiddenErr.Code)
		})
	})

	t.Run("Given the broker fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite) {
			suite.mockAgentAccessVerifier.EXPECT().
				EnsureAccess(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).
				Times(1)
			suite.mockConversationBrokerRequester.EXPECT().
				SaveAgentPromptTemplate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			ctx := helpers.ContextWithCallerId(
				context.Background(),
				strings.ToLower(common_tests.RandomEthAddress().Hex()),
			)
			initSuite(suite)

			_, err := suite.sut.Save(ctx, generateValidSaveAgentPromptTemplateRequest())

			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Given the caller has access to the agent and the broker saves the template", func(t *testing.T) {
		t.Parallel()

		initSuite := func(
			suite *WhenSavingAgentPromptTemplateTestingSuite,
			request *requests.SaveAgentPromptTemplateRequest,
			callerId string,
		) {
			suite.mockAgentAccessVerifier.EXPECT().
				EnsureAccess(gomock.Any(), callerId, strings.ToLower(request.AgentCollectionAddress), uint64(1)).
				Return(nil).
				Times(1)
			suite.mockConversationBrokerRequester.EXPECT().
				SaveAgentPromptTemplate(
					gomock.Any(),
					strings.ToLower(request.AgentCollectionAddress),
					uint64(1),
					"clone-system",
					request.Content,
				).
				Return(&responses.AgentPromptTemplateResponse{Key: "clone-system", Content: request.Content}, nil).
				Times(1)
		}

		t.Run("Should return the saved template", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			callerId := strings.ToLower(common_tests.RandomEthAddress().Hex())
			ctx := helpers.ContextWithCallerId(context.Background(), callerId)
			request := generateValidSaveAgentPromptTemplateRequest()
			initSuite(suite, request, callerId)

			response, err := suite.sut.Save(ctx, request)

			assert.NoError(t, err)
			assert.Equal(t, request.Content, response.Content)
		})
	})
}
