package use_cases_test

import (
	"fmt"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenListingAgentCollectionsTestSuite struct {
	agentRegistryRequesterMock *interfaces_mocks.MockAgentRegistryRequesterInterface

	sut *use_cases.ListAgentCollections
}

func WhenListingAgentCollectionsTestSuiteBeforeEach(t *testing.T) *WhenListingAgentCollectionsTestSuite {
	// Create mocks
	mockController := gomock.NewController(t, gomock.WithOverridableExpectations())

	agentRegistryRequesterMock := interfaces_mocks.NewMockAgentRegistryRequesterInterface(mockController)

	sut := use_cases.NewListAgentCollections(
		agentRegistryRequesterMock,
	)

	return &WhenListingAgentCollectionsTestSuite{
		sut: sut,

		agentRegistryRequesterMock: agentRegistryRequesterMock,
	}
}

func TestWhenListingAgentCollections(t *testing.T) {
	t.Parallel()

	t.Run("Given error happens while fetching collections from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentCollectionsTestSuite) {
			suite.agentRegistryRequesterMock.EXPECT().ListAgentCollections().Return(nil, assert.AnError)
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

		initSuite := func(suite *WhenListingAgentCollectionsTestSuite) []*responses.AgentCollectionResponse {
			collections := []*responses.AgentCollectionResponse{}

			for i := 0; i < 2; i++ {
				collection := &responses.AgentCollectionResponse{
					Address: fmt.Sprintf("0x%040x", i+1),
					Name:    fmt.Sprintf("Collection%d", i+1),
					Symbol:  fmt.Sprintf("COL%d", i+1),
				}
				collections = append(collections, collection)
			}

			suite.agentRegistryRequesterMock.EXPECT().ListAgentCollections().Return(collections, nil)

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
