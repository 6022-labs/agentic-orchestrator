package use_cases_test

import (
	"errors"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type WhenRegisteringAgentTestingSuite struct {
	sut *use_cases.RegisterAgent

	mockAgentRegistryRequester *interfaces_mocks.MockAgentRegistryRequesterInterface
	mockParticipantRepository  *repositories_mocks.MockParticipantRepositoryInterface
}

func WhenRegisteringAgentBeforeEach(t *testing.T) *WhenRegisteringAgentTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentRegistryRequester := interfaces_mocks.NewMockAgentRegistryRequesterInterface(mockController)
	mockParticipantRepository := repositories_mocks.NewMockParticipantRepositoryInterface(mockController)

	sut := use_cases.NewRegisterAgent(mockAgentRegistryRequester, mockParticipantRepository)

	return &WhenRegisteringAgentTestingSuite{
		sut: sut,

		mockAgentRegistryRequester: mockAgentRegistryRequester,
		mockParticipantRepository:  mockParticipantRepository,
	}
}

func generateValidRegisterAgentRequest() *requests.RegisterAgentRequest {
	return &requests.RegisterAgentRequest{
		CollectionAddress: "0xabc",
		CollectionTokenId: "1",
	}
}

func TestWhenRegisteringAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because collectionAddress is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringAgentBeforeEach(t)
			request := generateValidRegisterAgentRequest()
			request.CollectionAddress = ""

			_, _, err := suite.sut.Register(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "collectionAddress")
		})

		t.Run("Because collectionTokenId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringAgentBeforeEach(t)
			request := generateValidRegisterAgentRequest()
			request.CollectionTokenId = ""

			_, _, err := suite.sut.Register(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "collectionTokenId")
		})

		t.Run("Because collectionTokenId is not a number", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringAgentBeforeEach(t)
			request := generateValidRegisterAgentRequest()
			request.CollectionTokenId = "abc"

			_, _, err := suite.sut.Register(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "collectionTokenId")
		})
	})

	t.Run("Given the agent participant already exists", func(t *testing.T) {
		t.Parallel()

		existing := model_builders.NewAgentParticipantBuilder().Build()

		initSuite := func(suite *WhenRegisteringAgentTestingSuite, request *requests.RegisterAgentRequest) {
			suite.mockParticipantRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.CollectionAddress, uint64(1)).
				Return(existing, nil).
				Times(1)
		}

		t.Run("Should return the existing participant with existed=true and not call the registry or save", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringAgentBeforeEach(t)
			request := generateValidRegisterAgentRequest()
			initSuite(suite, request)

			participant, existed, err := suite.sut.Register(request)

			assert.NoError(t, err)
			assert.True(t, existed)
			assert.Equal(t, existing.ID, participant.ID)
		})
	})

	t.Run("Given the lookup fails with an unexpected error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringAgentTestingSuite, request *requests.RegisterAgentRequest) {
			suite.mockParticipantRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.CollectionAddress, uint64(1)).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringAgentBeforeEach(t)
			request := generateValidRegisterAgentRequest()
			initSuite(suite, request)

			_, _, err := suite.sut.Register(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given the agent does not exist yet and the registry call fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringAgentTestingSuite, request *requests.RegisterAgentRequest) {
			suite.mockParticipantRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.CollectionAddress, uint64(1)).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(request.CollectionAddress, uint64(1)).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringAgentBeforeEach(t)
			request := generateValidRegisterAgentRequest()
			initSuite(suite, request)

			_, _, err := suite.sut.Register(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given the agent does not exist yet and the registry returns no agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringAgentTestingSuite, request *requests.RegisterAgentRequest) {
			suite.mockParticipantRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.CollectionAddress, uint64(1)).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(request.CollectionAddress, uint64(1)).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringAgentBeforeEach(t)
			request := generateValidRegisterAgentRequest()
			initSuite(suite, request)

			_, _, err := suite.sut.Register(request)

			assert.IsType(t, &cb_errors.NotFoundError{}, err)
		})
	})

	t.Run("Given the agent does not exist yet and the registry returns an agent", func(t *testing.T) {
		t.Parallel()

		fromRegistry := &models.RegistryAgent{
			ID:                uuid.New(),
			Name:              "Ada",
			Role:              "expert",
			EnsDomain:         "ada.experts.80002.protocol6022.eth",
			CollectionAddress: "0xabc",
			CollectionTokenId: lo.ToPtr(uint64(1)),
		}

		initSuite := func(suite *WhenRegisteringAgentTestingSuite, request *requests.RegisterAgentRequest, registryAgent *models.RegistryAgent) *models.Participant {
			expected, err := registryAgent.ToParticipant()
			if err != nil {
				t.Fatalf("failed to map registry agent: %v", err)
			}

			suite.mockParticipantRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.CollectionAddress, uint64(1)).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(request.CollectionAddress, uint64(1)).
				Return(registryAgent, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				Save(expected).
				Return(expected, nil).
				Times(1)

			return expected
		}

		t.Run("Should persist the mapped participant and return existed=false", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringAgentBeforeEach(t)
			request := generateValidRegisterAgentRequest()
			initSuite(suite, request, fromRegistry)

			participant, existed, err := suite.sut.Register(request)

			assert.NoError(t, err)
			assert.False(t, existed)
			assert.Equal(t, fromRegistry.ID, participant.ID)
		})
	})
}
