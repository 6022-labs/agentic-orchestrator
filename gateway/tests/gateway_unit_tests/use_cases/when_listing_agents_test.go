package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenListingAgentsTestingSuite struct {
	sut *use_cases.ListAgents

	agentRegistryRequesterMock *interfaces_mocks.MockAgentRegistryRequesterInterface
}

func WhenListingAgentsBeforeEach(t *testing.T) *WhenListingAgentsTestingSuite {
	mockController := gomock.NewController(t)

	agentRegistryRequesterMock := interfaces_mocks.NewMockAgentRegistryRequesterInterface(mockController)

	sut := use_cases.NewListAgents(agentRegistryRequesterMock)

	return &WhenListingAgentsTestingSuite{
		sut: sut,

		agentRegistryRequesterMock: agentRegistryRequesterMock,
	}
}

func TestWhenListingAgents(t *testing.T) {
	t.Parallel()

	t.Run("Given error happens while fetching agents from the registry", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentsTestingSuite) *requests.ListAgentsRequest {
			suite.agentRegistryRequesterMock.
				EXPECT().
				ListAgents(gomock.Any()).
				Return(nil, assert.AnError)

			return &requests.ListAgentsRequest{}
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentsBeforeEach(t)
			req := initSuite(suite)

			_, err := suite.sut.List(req)

			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given agents are fetched successfully", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentsTestingSuite) []*responses.AgentResponse {
			agents := []*responses.AgentResponse{
				model_builders.NewAgentBuilder().Build(),
				model_builders.NewAgentBuilder().Build(),
			}

			suite.agentRegistryRequesterMock.
				EXPECT().
				ListAgents(gomock.Any()).
				Return(&responses.ListAgentsResponse{
					PaginateResponse: responses.PaginateResponse{
						Page:            1,
						ElementsPerPage: 10,
						TotalElements:   uint(len(agents)),
					},
					Agents: agents,
				}, nil)

			return agents
		}

		t.Run("Should return the paginated agents", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentsBeforeEach(t)
			agents := initSuite(suite)

			resp, err := suite.sut.List(&requests.ListAgentsRequest{})

			assert.Nil(t, err)
			assert.Len(t, resp.Agents, len(agents))
			assert.Equal(t, uint(len(agents)), resp.TotalElements)
		})
	})

	t.Run("Given a minted filter is provided", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentsTestingSuite) {
			suite.agentRegistryRequesterMock.
				EXPECT().
				ListAgents(gomock.Any()).
				DoAndReturn(func(request *interfaces.ListRegistryAgentsRequest) (*responses.ListAgentsResponse, error) {
					assert.Equal(t, lo.ToPtr(true), request.IsMinted)
					return &responses.ListAgentsResponse{}, nil
				})
		}

		t.Run("Should forward it to the registry", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentsBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.List(&requests.ListAgentsRequest{IsMinted: lo.ToPtr(true)})

			assert.Nil(t, err)
		})
	})

	t.Run("Given sort parameters are provided", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentsTestingSuite) {
			suite.agentRegistryRequesterMock.
				EXPECT().
				ListAgents(gomock.Any()).
				DoAndReturn(func(request *interfaces.ListRegistryAgentsRequest) (*responses.ListAgentsResponse, error) {
					assert.Equal(t, lo.ToPtr("name"), request.SortBy)
					assert.Equal(t, lo.ToPtr("desc"), request.SortDirection)
					return &responses.ListAgentsResponse{}, nil
				})
		}

		t.Run("Should forward them to the registry", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentsBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.List(&requests.ListAgentsRequest{
				PaginateRequest: requests.PaginateRequest{
					SortDirection: lo.ToPtr(requests.SortDirectionDesc),
				},
				SortBy: lo.ToPtr(requests.ListAgentSortByName),
			})

			assert.Nil(t, err)
		})
	})
}
