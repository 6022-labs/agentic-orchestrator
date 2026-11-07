package use_cases_test

import (
	"context"
	"strings"

	"errors"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
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
		CollectionAddress: "0xabc00000000000000000000000000000000000de",
		CollectionTokenId: 1,
	}
}

func TestWhenRegisteringAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because collectionAddress has no 0x prefix", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringAgentBeforeEach(t)
			request := generateValidRegisterAgentRequest()
			request.CollectionAddress = strings.TrimPrefix(request.CollectionAddress, "0x")

			_, _, err := suite.sut.Register(context.Background(), request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given a checksummed collection address", func(t *testing.T) {
		t.Parallel()

		fromRegistry := model_builders.NewRegistryAgentBuilder().
			WithCollection(" 0xAbC00000000000000000000000000000000000dE ", 1).
			Build()

		initSuite := func(suite *WhenRegisteringAgentTestingSuite) {
			suite.mockParticipantRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), "0xabc00000000000000000000000000000000000de", uint64(1)).
				Return(nil, nil).
				Times(1)
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), "0xabc00000000000000000000000000000000000de", uint64(1)).
				Return(fromRegistry, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, participant *models.Participant) (*models.Participant, error) {
					return participant, nil
				}).
				Times(1)
		}

		t.Run("Should look up, fetch and persist the lowercased address", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringAgentBeforeEach(t)
			request := generateValidRegisterAgentRequest()
			request.CollectionAddress = "0xAbC00000000000000000000000000000000000dE"
			initSuite(suite)

			participant, _, err := suite.sut.Register(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, "0xabc00000000000000000000000000000000000de", *participant.CollectionAddress)
		})
	})

	t.Run("Given the lookup fails with an unexpected error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringAgentTestingSuite, request *requests.RegisterAgentRequest) {
			suite.mockParticipantRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), request.CollectionAddress, uint64(1)).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringAgentBeforeEach(t)
			request := generateValidRegisterAgentRequest()
			initSuite(suite, request)

			_, _, err := suite.sut.Register(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given the agent participant already exists", func(t *testing.T) {
		t.Parallel()

		existing := model_builders.NewAgentParticipantBuilder().Build()

		initSuite := func(suite *WhenRegisteringAgentTestingSuite, request *requests.RegisterAgentRequest) {
			suite.mockParticipantRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), request.CollectionAddress, uint64(1)).
				Return(existing, nil).
				Times(1)
		}

		t.Run(
			"Should return the existing participant with existed=true and not call the registry or save",
			func(t *testing.T) {
				t.Parallel()

				suite := WhenRegisteringAgentBeforeEach(t)
				request := generateValidRegisterAgentRequest()
				initSuite(suite, request)

				participant, existed, err := suite.sut.Register(context.Background(), request)

				assert.NoError(t, err)
				assert.True(t, existed)
				assert.Equal(t, existing.ID, participant.ID)
			},
		)
	})

	t.Run("Given the agent does not exist yet and the registry call fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringAgentTestingSuite, request *requests.RegisterAgentRequest) {
			suite.mockParticipantRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), request.CollectionAddress, uint64(1)).
				Return(nil, nil).
				Times(1)
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), request.CollectionAddress, uint64(1)).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringAgentBeforeEach(t)
			request := generateValidRegisterAgentRequest()
			initSuite(suite, request)

			_, _, err := suite.sut.Register(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given the agent does not exist yet and the registry returns no agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringAgentTestingSuite, request *requests.RegisterAgentRequest) {
			suite.mockParticipantRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), request.CollectionAddress, uint64(1)).
				Return(nil, nil).
				Times(1)
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), request.CollectionAddress, uint64(1)).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringAgentBeforeEach(t)
			request := generateValidRegisterAgentRequest()
			initSuite(suite, request)

			_, _, err := suite.sut.Register(context.Background(), request)

			var notFoundErr *cb_errors.NotFoundError
			require.ErrorAs(t, err, &notFoundErr)
			assert.Equal(t, "agent_not_found", notFoundErr.Code)
		})
	})

	t.Run("Given the agent does not exist yet and the registry returns an agent", func(t *testing.T) {
		t.Parallel()

		fromRegistry := model_builders.NewRegistryAgentBuilder().
			WithName("Ada").
			WithRole("expert").
			WithEnsDomain("ada.experts.80002.6022.eth").
			WithCollection("0xabc00000000000000000000000000000000000de", 1).
			Build()

		initSuite := func(suite *WhenRegisteringAgentTestingSuite, request *requests.RegisterAgentRequest, registryAgent *models.RegistryAgent) *models.Participant {
			expected := registryAgent.ToParticipant()

			suite.mockParticipantRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), request.CollectionAddress, uint64(1)).
				Return(nil, nil).
				Times(1)
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), request.CollectionAddress, uint64(1)).
				Return(registryAgent, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				Save(gomock.Any(), expected).
				Return(expected, nil).
				Times(1)

			return expected
		}

		t.Run("Should persist the mapped participant and return existed=false", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringAgentBeforeEach(t)
			request := generateValidRegisterAgentRequest()
			initSuite(suite, request, fromRegistry)

			participant, existed, err := suite.sut.Register(context.Background(), request)

			assert.NoError(t, err)
			assert.False(t, existed)
			assert.Equal(t, fromRegistry.ID, participant.ID)
		})
	})
}
