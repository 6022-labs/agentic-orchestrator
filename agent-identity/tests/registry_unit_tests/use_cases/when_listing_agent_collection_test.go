package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenListingAgentCollectionTestingSuite struct {
	sut *use_cases.ListAgentCollections

	mockAgentCollectionRepository *repositories_mocks.MockAgentCollectionRepositoryInterface
}

func WhenListingAgentCollectionBeforeEach(t *testing.T) *WhenListingAgentCollectionTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentCollectionRepository := repositories_mocks.NewMockAgentCollectionRepositoryInterface(mockController)

	sut := use_cases.NewListAgentCollections(mockAgentCollectionRepository)

	return &WhenListingAgentCollectionTestingSuite{
		sut:                           sut,
		mockAgentCollectionRepository: mockAgentCollectionRepository,
	}
}

func TestWhenListingAgentCollection(t *testing.T) {
	t.Parallel()

	t.Run("Given error happens while fetching collections from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentCollectionTestingSuite) {
			suite.mockAgentCollectionRepository.EXPECT().
				GetAll().
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenListingAgentCollectionBeforeEach(t)
			initSuite(basics)

			_, err := basics.sut.List()

			assert.NotNil(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given collections are fetched successfully", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return existing collections", func(t *testing.T) {
			t.Parallel()
			basics := WhenListingAgentCollectionBeforeEach(t)

			collections := []*models.AgentCollection{}
			for i := 0; i < 2; i++ {
				collection := model_builders.NewAgentCollectionBuilder().Build()
				collections = append(collections, collection)
			}

			basics.mockAgentCollectionRepository.EXPECT().
				GetAll().
				Return(collections, nil).
				Times(1)

			result, err := basics.sut.List()

			assert.Nil(t, err)
			assert.Len(t, result, len(collections))

			for i, collection := range collections {
				assert.Equal(t, collection.Address, result[i].Address)
				assert.Equal(t, collection.Name, result[i].Name)
				assert.Equal(t, collection.Symbol, result[i].Symbol)
			}
		})
	})

	t.Run("Given no collections exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an empty list", func(t *testing.T) {
			t.Parallel()
			basics := WhenListingAgentCollectionBeforeEach(t)

			basics.mockAgentCollectionRepository.EXPECT().
				GetAll().
				Return([]*models.AgentCollection{}, nil).
				Times(1)

			result, err := basics.sut.List()

			assert.Nil(t, err)
			assert.Empty(t, result)
		})
	})
}
