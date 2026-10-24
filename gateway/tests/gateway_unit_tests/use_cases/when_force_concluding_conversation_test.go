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
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenForceConcludingConversationTestingSuite struct {
	sut *use_cases.ForceConcludeConversation

	mockConversationAccessVerifier  *interfaces_mocks.MockConversationAccessVerifierInterface
	mockConversationBrokerRequester *interfaces_mocks.MockConversationBrokerRequesterInterface
}

func WhenForceConcludingConversationBeforeEach(t *testing.T) *WhenForceConcludingConversationTestingSuite {
	mockController := gomock.NewController(t)

	mockConversationAccessVerifier := interfaces_mocks.NewMockConversationAccessVerifierInterface(mockController)
	mockConversationBrokerRequester := interfaces_mocks.NewMockConversationBrokerRequesterInterface(mockController)

	sut := use_cases.NewForceConcludeConversation(mockConversationAccessVerifier, mockConversationBrokerRequester)

	return &WhenForceConcludingConversationTestingSuite{
		sut: sut,

		mockConversationAccessVerifier:  mockConversationAccessVerifier,
		mockConversationBrokerRequester: mockConversationBrokerRequester,
	}
}

func TestWhenForceConcludingConversation(t *testing.T) {
	t.Parallel()

	generateValidConversationByIdRequest := func() *requests.ConversationByIdRequest {
		return &requests.ConversationByIdRequest{ConversationId: uuid.New()}
	}

	t.Run("Given request is invalid because conversation id is empty", func(t *testing.T) {
		t.Parallel()

		request := generateValidConversationByIdRequest()
		request.ConversationId = uuid.Nil

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenForceConcludingConversationBeforeEach(t)

			err := suite.sut.Conclude(context.Background(), request)

			var validationErr *cb_errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "id", validationErr.Field)
		})
	})

	t.Run("Given no caller id", func(t *testing.T) {
		t.Parallel()

		request := generateValidConversationByIdRequest()

		t.Run("Should return an unauthorized error", func(t *testing.T) {
			t.Parallel()

			suite := WhenForceConcludingConversationBeforeEach(t)

			err := suite.sut.Conclude(context.Background(), request)

			var unauthorizedErr *cb_errors.UnauthorizedError
			assert.ErrorAs(t, err, &unauthorizedErr)
			assert.Equal(t, "missing_caller_id", unauthorizedErr.Code)
		})
	})

	t.Run("Given the caller has no access to the conversation", func(t *testing.T) {
		t.Parallel()

		request := generateValidConversationByIdRequest()

		initSuite := func(
			suite *WhenForceConcludingConversationTestingSuite,
			request *requests.ConversationByIdRequest,
			callerId string,
		) {
			suite.mockConversationAccessVerifier.EXPECT().
				EnsureAccess(gomock.Any(), callerId, request.ConversationId).
				Return(cb_errors.NewForbiddenError(services.ConversationAccessForbiddenCode, "forbidden")).
				Times(1)
		}

		t.Run("Should return the forbidden error without calling the broker", func(t *testing.T) {
			t.Parallel()

			suite := WhenForceConcludingConversationBeforeEach(t)
			callerId := strings.ToLower(common_tests.RandomEthAddress().Hex())
			ctx := helpers.ContextWithCallerId(context.Background(), callerId)
			initSuite(suite, request, callerId)

			err := suite.sut.Conclude(ctx, request)

			var forbiddenErr *cb_errors.ForbiddenError
			assert.ErrorAs(t, err, &forbiddenErr)
			assert.Equal(t, services.ConversationAccessForbiddenCode, forbiddenErr.Code)
		})
	})

	t.Run("Given the broker fails", func(t *testing.T) {
		t.Parallel()

		request := generateValidConversationByIdRequest()

		initSuite := func(suite *WhenForceConcludingConversationTestingSuite) {
			suite.mockConversationAccessVerifier.EXPECT().
				EnsureAccess(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).
				Times(1)
			suite.mockConversationBrokerRequester.EXPECT().
				ForceConcludeConversation(gomock.Any(), gomock.Any()).
				Return(assert.AnError).
				Times(1)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenForceConcludingConversationBeforeEach(t)
			ctx := helpers.ContextWithCallerId(
				context.Background(),
				strings.ToLower(common_tests.RandomEthAddress().Hex()),
			)
			initSuite(suite)

			err := suite.sut.Conclude(ctx, request)

			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Given the caller has access to the conversation", func(t *testing.T) {
		t.Parallel()

		request := generateValidConversationByIdRequest()

		initSuite := func(
			suite *WhenForceConcludingConversationTestingSuite,
			request *requests.ConversationByIdRequest,
			callerId string,
		) {
			suite.mockConversationAccessVerifier.EXPECT().
				EnsureAccess(gomock.Any(), callerId, request.ConversationId).
				Return(nil).
				Times(1)
			suite.mockConversationBrokerRequester.EXPECT().
				ForceConcludeConversation(gomock.Any(), request.ConversationId).
				Return(nil).
				Times(1)
		}

		t.Run("Should return the broker result", func(t *testing.T) {
			t.Parallel()

			suite := WhenForceConcludingConversationBeforeEach(t)
			callerId := strings.ToLower(common_tests.RandomEthAddress().Hex())
			ctx := helpers.ContextWithCallerId(context.Background(), callerId)
			initSuite(suite, request, callerId)

			err := suite.sut.Conclude(ctx, request)

			assert.NoError(t, err)
		})
	})
}
