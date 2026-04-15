package use_cases_test

import (
	"errors"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenRegisteringNonMintedParticipantTestingSuite struct {
	sut *use_cases.RegisterNonMintedParticipant

	mockParticipantRepository *repositories_mocks.MockParticipantRepositoryInterface
}

func WhenRegisteringNonMintedParticipantBeforeEach(t *testing.T) *WhenRegisteringNonMintedParticipantTestingSuite {
	mockController := gomock.NewController(t)

	mockParticipantRepository := repositories_mocks.NewMockParticipantRepositoryInterface(mockController)

	sut := use_cases.NewRegisterNonMintedParticipant(mockParticipantRepository)

	return &WhenRegisteringNonMintedParticipantTestingSuite{
		sut: sut,

		mockParticipantRepository: mockParticipantRepository,
	}
}

func generateValidRegisterNonMintedParticipantRequest() *requests.RegisterNonMintedParticipantRequest {
	return &requests.RegisterNonMintedParticipantRequest{
		DisplayName: "Pierre",
	}
}

func TestWhenRegisteringNonMintedParticipant(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because displayName is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringNonMintedParticipantBeforeEach(t)
			request := generateValidRegisterNonMintedParticipantRequest()
			request.DisplayName = ""

			_, err := suite.sut.Register(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "displayName")
		})
	})

	t.Run("Given saving the participant fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringNonMintedParticipantTestingSuite) {
			suite.mockParticipantRepository.EXPECT().
				Save(gomock.Any()).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringNonMintedParticipantBeforeEach(t)
			request := generateValidRegisterNonMintedParticipantRequest()
			initSuite(suite)

			_, err := suite.sut.Register(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given request is valid", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringNonMintedParticipantTestingSuite, request *requests.RegisterNonMintedParticipantRequest) {
			suite.mockParticipantRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(p *models.Participant) (*models.Participant, error) {
					assert.NotNil(t, p.DisplayName)
					assert.Equal(t, request.DisplayName, *p.DisplayName)
					assert.Nil(t, p.CollectionAddress)
					assert.Nil(t, p.CollectionTokenId)
					assert.Nil(t, p.Role)
					return p, nil
				}).
				Times(1)
		}

		t.Run("Should persist participant with displayName only and minted fields nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringNonMintedParticipantBeforeEach(t)
			request := generateValidRegisterNonMintedParticipantRequest()
			initSuite(suite, request)

			_, err := suite.sut.Register(request)

			assert.NoError(t, err)
		})
	})
}
