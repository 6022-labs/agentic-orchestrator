package services_test

import (
	"net/http"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/services"
	clients_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_http_agent_identity_registry_mocks/clients_mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenListingAgentCollectionsTestingSuite struct {
	sut *services.AgentRegistryRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenListingAgentCollectionsBeforeEach(t *testing.T) *WhenListingAgentCollectionsTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	return &WhenListingAgentCollectionsTestingSuite{
		sut:        services.NewAgentRegistryRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenListingAgentCollections(t *testing.T) {
	t.Parallel()

	t.Run("Given the registry returns collections", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentCollectionsTestingSuite) {
			suite.mockClient.EXPECT().GetAgentCollectionsWithResponse(gomock.Any()).Return(&clients.GetAgentCollectionsResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &[]clients.ResponsesAgentCollectionResponse{{
					Address: lo.ToPtr("0xabc"),
					Name:    lo.ToPtr("Experts"),
					Symbol:  lo.ToPtr("EXP"),
				}},
			}, nil).Times(1)
		}

		t.Run("Should return the mapped collections", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentCollectionsBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ListAgentCollections()

			assert.NoError(t, err)
			if assert.Len(t, result, 1) {
				assert.Equal(t, "0xabc", result[0].Address)
				assert.Equal(t, "Experts", result[0].Name)
				assert.Equal(t, "EXP", result[0].Symbol)
			}
		})
	})

	t.Run("Given the registry fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentCollectionsTestingSuite) {
			suite.mockClient.EXPECT().GetAgentCollectionsWithResponse(gomock.Any()).Return(&clients.GetAgentCollectionsResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusInternalServerError},
				Body:         []byte("boom"),
			}, nil).Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentCollectionsBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ListAgentCollections()

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
}
