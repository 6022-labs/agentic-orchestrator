package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenUpdatingAgentOwnerTestSuite struct {
	sut *use_cases.UpdateAgentOwner

	agentRepositoryMock *repositories_mocks.MockAgentRepositoryInterface
}

func WhenUpdatingAgentOwnerBeforeEach(t *testing.T) *WhenUpdatingAgentOwnerTestSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)

	sut := use_cases.NewUpdateAgentOwner(agentRepositoryMock)

	return &WhenUpdatingAgentOwnerTestSuite{
		sut: sut,

		agentRepositoryMock: agentRepositoryMock,
	}
}

func TestWhenUpdatingAgentOwner(t *testing.T) {
	t.Parallel()

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			suite := WhenUpdatingAgentOwnerBeforeEach(t)

			collectionTokenId := uint64(1)
			collectionAddress := common_tests.RandomEthAddress().Hex()
			newOwner := common_tests.RandomEthAddress().Hex()
			request := &requests.UpdateAgentOwnerRequest{
				CollectionTokenId: collectionTokenId,
				CollectionAddress: collectionAddress,
				NewOwner:          newOwner,
			}

			suite.agentRepositoryMock.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.CollectionAddress, request.CollectionTokenId).
				Return(nil, assert.AnError)

			err := suite.sut.Update(request)

			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given entity exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should update entity", func(t *testing.T) {
			suite := WhenUpdatingAgentOwnerBeforeEach(t)

			agent := model_builders.NewAgentBuilder().Build()

			collectionAddress := common_tests.RandomEthAddress().Hex()
			newOwner := common_tests.RandomEthAddress().Hex()
			request := &requests.UpdateAgentOwnerRequest{
				CollectionTokenId: *agent.CollectionTokenId,
				CollectionAddress: collectionAddress,
				NewOwner:          newOwner,
			}

			suite.agentRepositoryMock.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.CollectionAddress, request.CollectionTokenId).
				Return(agent, nil)

			suite.agentRepositoryMock.EXPECT().
				Save(agent).
				Return(agent, nil)

			err := suite.sut.Update(request)

			assert.Nil(t, err)
		})
	})
}
