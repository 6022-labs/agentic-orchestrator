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
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenRegisteringMintedParticipantTestingSuite struct {
	sut *use_cases.RegisterMintedParticipant

	mockAgentRegistryRequester *interfaces_mocks.MockAgentRegistryRequesterInterface
	mockParticipantRepository  *repositories_mocks.MockParticipantRepositoryInterface
}

func WhenRegisteringMintedParticipantBeforeEach(t *testing.T) *WhenRegisteringMintedParticipantTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentRegistryRequester := interfaces_mocks.NewMockAgentRegistryRequesterInterface(mockController)
	mockParticipantRepository := repositories_mocks.NewMockParticipantRepositoryInterface(mockController)

	sut := use_cases.NewRegisterMintedParticipant(mockAgentRegistryRequester, mockParticipantRepository)

	return &WhenRegisteringMintedParticipantTestingSuite{
		sut: sut,

		mockAgentRegistryRequester: mockAgentRegistryRequester,
		mockParticipantRepository:  mockParticipantRepository,
	}
}

func generateValidRegisterMintedParticipantRequest() *requests.RegisterMintedParticipantRequest {
	return &requests.RegisterMintedParticipantRequest{
		CollectionAddress: "0xabc",
		CollectionTokenId: "1",
	}
}

func TestWhenRegisteringMintedParticipant(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because collectionAddress is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringMintedParticipantBeforeEach(t)
			request := generateValidRegisterMintedParticipantRequest()
			request.CollectionAddress = ""

			_, err := suite.sut.Register(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "collectionAddress")
		})

		t.Run("Because collectionTokenId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringMintedParticipantBeforeEach(t)
			request := generateValidRegisterMintedParticipantRequest()
			request.CollectionTokenId = ""

			_, err := suite.sut.Register(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "collectionTokenId")
		})

		t.Run("Because collectionTokenId is not a number", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringMintedParticipantBeforeEach(t)
			request := generateValidRegisterMintedParticipantRequest()
			request.CollectionTokenId = "abc"

			_, err := suite.sut.Register(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "collectionTokenId")
		})
	})

	t.Run("Given the registry call fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringMintedParticipantTestingSuite, request *requests.RegisterMintedParticipantRequest) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetMintedAgent(request.CollectionAddress, uint64(1)).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringMintedParticipantBeforeEach(t)
			request := generateValidRegisterMintedParticipantRequest()
			initSuite(suite, request)

			_, err := suite.sut.Register(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given the registry returns no agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringMintedParticipantTestingSuite, request *requests.RegisterMintedParticipantRequest) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetMintedAgent(request.CollectionAddress, uint64(1)).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringMintedParticipantBeforeEach(t)
			request := generateValidRegisterMintedParticipantRequest()
			initSuite(suite, request)

			_, err := suite.sut.Register(request)

			assert.IsType(t, &cb_errors.NotFoundError{}, err)
		})
	})

	t.Run("Given the registry returns an agent", func(t *testing.T) {
		t.Parallel()

		fromRegistry := model_builders.NewParticipantBuilder().Build()

		initSuite := func(suite *WhenRegisteringMintedParticipantTestingSuite, request *requests.RegisterMintedParticipantRequest, participant *models.Participant) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetMintedAgent(request.CollectionAddress, uint64(1)).
				Return(participant, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				Save(participant).
				Return(participant, nil).
				Times(1)
		}

		t.Run("Should persist the participant returned by the registry", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringMintedParticipantBeforeEach(t)
			request := generateValidRegisterMintedParticipantRequest()
			initSuite(suite, request, fromRegistry)

			participant, err := suite.sut.Register(request)

			assert.NoError(t, err)
			assert.Equal(t, fromRegistry.ID, participant.ID)
		})
	})
}
