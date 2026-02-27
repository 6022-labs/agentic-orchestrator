package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenListingAgentCollectionsTestSuite struct {
	agentCollectionRepositoryMock *repositories_mocks.MockAgentCollectionRepositoryInterface

	sut *use_cases.ListAgentCollections
}

func WhenListingAgentCollectionsTestSuiteBeforeEach(t *testing.T) *WhenListingAgentCollectionsTestSuite {
	// Create mocks
	mockController := gomock.NewController(t, gomock.WithOverridableExpectations())

	agentCollectionRepositoryMock := repositories_mocks.NewMockAgentCollectionRepositoryInterface(mockController)

	sut := use_cases.NewListAgentCollections(
		agentCollectionRepositoryMock,
	)

	return &WhenListingAgentCollectionsTestSuite{
		sut: sut,

		agentCollectionRepositoryMock: agentCollectionRepositoryMock,
	}
}

func TestWhenListingAgentCollections(t *testing.T) {
	t.Parallel()

	t.Run("Given error happens while fetching collections from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentCollectionsTestSuite) {
			suite.agentCollectionRepositoryMock.EXPECT().GetAll().Return(nil, assert.AnError)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			basics := WhenListingAgentCollectionsTestSuiteBeforeEach(t)
			initSuite(basics)

			_, err := basics.sut.List()

			assert.NotNil(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given collections are fetched successfully", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentCollectionsTestSuite) []*models.AgentCollection {
			collections := []*models.AgentCollection{}

			for i := 0; i < 2; i++ {
				collection := model_builders.NewAgentCollectionBuilder().Build()
				collections = append(collections, collection)
			}

			suite.agentCollectionRepositoryMock.EXPECT().GetAll().Return(collections, nil)

			return collections
		}

		t.Run("Should return existing collections", func(t *testing.T) {
			t.Parallel()

			basics := WhenListingAgentCollectionsTestSuiteBeforeEach(t)
			existingCollections := initSuite(basics)

			fetchedCollections, err := basics.sut.List()

			assert.Nil(t, err)
			assert.Len(t, fetchedCollections, len(existingCollections))

			for i, existingCollection := range existingCollections {
				fetchedCollection := fetchedCollections[i]
				assert.Equal(t, existingCollection.Name, fetchedCollection.Name)
				assert.Equal(t, existingCollection.Symbol, fetchedCollection.Symbol)
				assert.Equal(t, existingCollection.Address, fetchedCollection.Address)
			}
		})
	})
}
