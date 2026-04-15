package use_cases_test

import (
	"errors"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenRegisteringBridgeTestingSuite struct {
	sut *use_cases.RegisterBridge

	mockBridgeRepository *repositories_mocks.MockBridgeRepositoryInterface
}

func WhenRegisteringBridgeBeforeEach(t *testing.T) *WhenRegisteringBridgeTestingSuite {
	mockController := gomock.NewController(t)

	mockBridgeRepository := repositories_mocks.NewMockBridgeRepositoryInterface(mockController)

	sut := use_cases.NewRegisterBridge(mockBridgeRepository)

	return &WhenRegisteringBridgeTestingSuite{
		sut: sut,

		mockBridgeRepository: mockBridgeRepository,
	}
}

func generateValidRegisterBridgeRequest() *requests.RegisterBridgeRequest {
	return &requests.RegisterBridgeRequest{
		Kind:        enums.BridgeKindSlack,
		CallbackUrl: "http://slack-bridge.local/internal/broker/callback",
	}
}

func TestWhenRegisteringBridge(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because kind is not a known bridge kind", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringBridgeBeforeEach(t)
			request := generateValidRegisterBridgeRequest()
			request.Kind = enums.BridgeKind("unknown")

			_, err := suite.sut.Register(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "kind")
		})

		t.Run("Because callbackUrl is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringBridgeBeforeEach(t)
			request := generateValidRegisterBridgeRequest()
			request.CallbackUrl = ""

			_, err := suite.sut.Register(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "callbackUrl")
		})
	})

	t.Run("Given saving the bridge fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringBridgeTestingSuite) {
			suite.mockBridgeRepository.EXPECT().
				Save(gomock.Any()).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringBridgeBeforeEach(t)
			request := generateValidRegisterBridgeRequest()
			initSuite(suite)

			_, err := suite.sut.Register(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given request is valid", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringBridgeTestingSuite, request *requests.RegisterBridgeRequest) {
			suite.mockBridgeRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(b *models.Bridge) (*models.Bridge, error) {
					assert.Equal(t, request.Kind, b.Kind)
					assert.Equal(t, request.CallbackUrl, b.CallbackUrl)
					return b, nil
				}).
				Times(1)
		}

		t.Run("Should persist and return the bridge", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringBridgeBeforeEach(t)
			request := generateValidRegisterBridgeRequest()
			initSuite(suite, request)

			bridge, err := suite.sut.Register(request)

			assert.NoError(t, err)
			assert.Equal(t, request.Kind, bridge.Kind)
			assert.Equal(t, request.CallbackUrl, bridge.CallbackUrl)
		})
	})
}
