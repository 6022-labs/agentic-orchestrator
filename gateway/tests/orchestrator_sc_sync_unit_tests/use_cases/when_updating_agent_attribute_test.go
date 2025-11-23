package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/model_builders"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenUpdatingAgentMissionTestSuite struct {
	sut *use_cases.UpdateAgentAttribute

	agentRepositoryMock *repositories_mocks.MockAgentRepositoryInterface
}

func WhenUpdatingAgentMissionBeforeEach(t *testing.T) *WhenUpdatingAgentMissionTestSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)

	sut := use_cases.NewUpdateAgentAttribute(agentRepositoryMock)

	return &WhenUpdatingAgentMissionTestSuite{
		sut: sut,

		agentRepositoryMock: agentRepositoryMock,
	}
}

func TestWhenUpdatingAgentMission(t *testing.T) {
	t.Parallel()

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			suite := WhenUpdatingAgentMissionBeforeEach(t)

			collectionTokenId := uint64(1)
			collectionAddress := common_tests.RandomEthAddress().Hex()
			attributeKey := "mission"
			newAttribute := "New attribute value"
			request := &requests.UpdateAgentAttributeRequest{
				CollectionTokenId: collectionTokenId,
				CollectionAddress: collectionAddress,
				Key:               attributeKey,
				NewValue:          newAttribute,
			}

			suite.agentRepositoryMock.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress, collectionTokenId).
				Return(nil, assert.AnError)

			err := suite.sut.Update(request)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given entity exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should update mission", func(t *testing.T) {
			suite := WhenUpdatingAgentMissionBeforeEach(t)

			agent := model_builders.NewAgentBuilder().Build()
			collectionTokenId := uint64(*agent.CollectionTokenId)
			collectionAddress := common_tests.RandomEthAddress().Hex()
			attributeKey := "mission"
			newAttribute := "New attribute value"
			request := &requests.UpdateAgentAttributeRequest{
				CollectionTokenId: collectionTokenId,
				CollectionAddress: collectionAddress,
				Key:               attributeKey,
				NewValue:          newAttribute,
			}

			suite.agentRepositoryMock.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(collectionAddress, collectionTokenId).
				Return(agent, nil)

			suite.agentRepositoryMock.EXPECT().
				Save(agent).
				DoAndReturn(func(a *models.Agent) (*models.Agent, error) {
					assert.Equal(t, newAttribute, a.Attributes[attributeKey])
					return a, nil
				})

			err := suite.sut.Update(request)
			assert.Nil(t, err)
		})
	})
}
