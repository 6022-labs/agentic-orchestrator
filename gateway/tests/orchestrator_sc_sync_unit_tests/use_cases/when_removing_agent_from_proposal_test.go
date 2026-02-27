package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type WhenRemovingAgentFromProposalTestingSuite struct {
	sut *use_cases.RemoveAgentFromProposal

	agentRepositoryMock *repositories_mocks.MockAgentRepositoryInterface
}

func WhenRemovingAgentFromProposalBeforeEach(t *testing.T) *WhenRemovingAgentFromProposalTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)

	sut := use_cases.NewRemoveAgentFromProposal(
		agentRepositoryMock,
	)

	return &WhenRemovingAgentFromProposalTestingSuite{
		sut: sut,

		agentRepositoryMock: agentRepositoryMock,
	}
}

func generateValidRemoveAgentFromProposalRequest() *requests.RemoveAgentFromProposalRequest {
	collectionAddress := common_tests.RandomEthAddress().Hex()

	return &requests.RemoveAgentFromProposalRequest{
		CollectionAddress: collectionAddress,
		ProposalId:        1,
	}
}

func TestWhenRemovingAgentFromProposal(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because collection address is not a valid Ethereum address", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromProposalBeforeEach(t)
			request := generateValidRemoveAgentFromProposalRequest()
			request.CollectionAddress = "invalid_address"

			err := suite.sut.Remove(request)

			assert.Error(t, err)
		})

		t.Run("Because proposal ID is not valid", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromProposalBeforeEach(t)
			request := generateValidRemoveAgentFromProposalRequest()
			request.ProposalId = 0

			err := suite.sut.Remove(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while fetching agent from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromProposalTestingSuite) {
			suite.agentRepositoryMock.
				EXPECT().
				GetNotYetMintedByCollectionAddressAndProposalId(gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromProposalBeforeEach(t)
			initSuite(suite)

			request := generateValidRemoveAgentFromProposalRequest()
			err := suite.sut.Remove(request)

			assert.Error(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given agent is not found for the provided collection address and proposal ID", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromProposalTestingSuite) {
			suite.agentRepositoryMock.
				EXPECT().
				GetNotYetMintedByCollectionAddressAndProposalId(gomock.Any(), gomock.Any()).
				Return(nil, gorm.ErrRecordNotFound)
		}

		t.Run("Should return an error indicating no agent found", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromProposalBeforeEach(t)
			initSuite(suite)

			request := generateValidRemoveAgentFromProposalRequest()
			err := suite.sut.Remove(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while removing agent from the proposal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromProposalTestingSuite) {
			agent := model_builders.NewAgentBuilder().Build()
			suite.agentRepositoryMock.
				EXPECT().
				GetNotYetMintedByCollectionAddressAndProposalId(gomock.Any(), gomock.Any()).
				Return(agent, nil)

			suite.agentRepositoryMock.
				EXPECT().
				Remove(agent).
				Return(assert.AnError)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromProposalBeforeEach(t)
			initSuite(suite)

			request := generateValidRemoveAgentFromProposalRequest()
			err := suite.sut.Remove(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent is successfully removed from the proposal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromProposalTestingSuite) {
			agent := model_builders.NewAgentBuilder().Build()
			suite.agentRepositoryMock.
				EXPECT().
				GetNotYetMintedByCollectionAddressAndProposalId(gomock.Any(), gomock.Any()).
				Return(agent, nil)

			suite.agentRepositoryMock.
				EXPECT().
				Remove(agent).
				Return(nil)
		}

		t.Run("Should successfully remove the agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromProposalBeforeEach(t)
			initSuite(suite)

			request := generateValidRemoveAgentFromProposalRequest()
			err := suite.sut.Remove(request)

			assert.NoError(t, err)
		})
	})
}
