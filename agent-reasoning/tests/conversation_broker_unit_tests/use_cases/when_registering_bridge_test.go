package use_cases_test

import (
	"errors"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/google/uuid"
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
		Kind:    enums.BridgeKindSlack,
		BaseUrl: "http://slack-bridge.local/internal/broker/callback",
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

			_, _, err := suite.sut.Register(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "kind")
		})

		t.Run("Because baseUrl is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringBridgeBeforeEach(t)
			request := generateValidRegisterBridgeRequest()
			request.BaseUrl = ""

			_, _, err := suite.sut.Register(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "baseUrl")
		})

		t.Run("Because id is the zero UUID", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringBridgeBeforeEach(t)
			request := generateValidRegisterBridgeRequest()
			zero := uuid.Nil
			request.ID = &zero

			_, _, err := suite.sut.Register(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "id")
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

			_, _, err := suite.sut.Register(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given request is valid and no id is provided", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringBridgeTestingSuite, request *requests.RegisterBridgeRequest) {
			suite.mockBridgeRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(b *models.Bridge) (*models.Bridge, error) {
					assert.Equal(t, request.Kind, b.Kind)
					assert.Equal(t, request.BaseUrl, b.BaseUrl)
					return b, nil
				}).
				Times(1)
		}

		t.Run("Should create a new bridge and signal it did not exist", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringBridgeBeforeEach(t)
			request := generateValidRegisterBridgeRequest()
			initSuite(suite, request)

			bridge, existed, err := suite.sut.Register(request)

			assert.NoError(t, err)
			assert.False(t, existed)
			assert.Equal(t, request.Kind, bridge.Kind)
		})
	})

	t.Run("Given request is valid and id is provided for an already-known bridge with matching fields", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringBridgeTestingSuite, existing *models.Bridge) {
			suite.mockBridgeRepository.EXPECT().GetById(existing.ID).Return(existing, nil).Times(1)
			// No Save expectation — idempotent when nothing changed.
		}

		t.Run("Should return the existing bridge and signal it already existed", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringBridgeBeforeEach(t)
			request := generateValidRegisterBridgeRequest()
			existing := model_builders.NewBridgeBuilder().
				WithKind(request.Kind).
				WithBaseUrl(request.BaseUrl).
				Build()
			id := existing.ID
			request.ID = &id
			initSuite(suite, existing)

			bridge, existed, err := suite.sut.Register(request)

			assert.NoError(t, err)
			assert.True(t, existed)
			assert.Equal(t, existing.ID, bridge.ID)
		})
	})

	t.Run("Given request is valid and id is provided for an already-known bridge with a different callback url", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringBridgeTestingSuite, existing *models.Bridge, newBaseUrl string) {
			suite.mockBridgeRepository.EXPECT().GetById(existing.ID).Return(existing, nil).Times(1)
			suite.mockBridgeRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(b *models.Bridge) (*models.Bridge, error) {
					assert.Equal(t, existing.ID, b.ID)
					assert.Equal(t, newBaseUrl, b.BaseUrl)
					return b, nil
				}).
				Times(1)
		}

		t.Run("Should update the callback url and signal it already existed", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringBridgeBeforeEach(t)
			existing := model_builders.NewBridgeBuilder().
				WithKind(enums.BridgeKindSlack).
				WithBaseUrl("http://old.local/cb").
				Build()
			request := generateValidRegisterBridgeRequest()
			request.BaseUrl = "http://new.local/cb"
			id := existing.ID
			request.ID = &id
			initSuite(suite, existing, request.BaseUrl)

			bridge, existed, err := suite.sut.Register(request)

			assert.NoError(t, err)
			assert.True(t, existed)
			assert.Equal(t, existing.ID, bridge.ID)
			assert.Equal(t, "http://new.local/cb", bridge.BaseUrl)
		})
	})

	t.Run("Given request is valid and id is provided for an unknown bridge", func(t *testing.T) {
		t.Parallel()

		newId := uuid.New()

		initSuite := func(suite *WhenRegisteringBridgeTestingSuite) {
			suite.mockBridgeRepository.EXPECT().GetById(newId).Return(nil, nil).Times(1)
			suite.mockBridgeRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(b *models.Bridge) (*models.Bridge, error) {
					assert.Equal(t, newId, b.ID)
					return b, nil
				}).
				Times(1)
		}

		t.Run("Should create the bridge using the supplied id", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringBridgeBeforeEach(t)
			request := generateValidRegisterBridgeRequest()
			request.ID = &newId
			initSuite(suite)

			bridge, existed, err := suite.sut.Register(request)

			assert.NoError(t, err)
			assert.False(t, existed)
			assert.Equal(t, newId, bridge.ID)
		})
	})
}
