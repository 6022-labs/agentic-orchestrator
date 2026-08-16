package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenListingAgentsTestingSuite struct {
	sut *use_cases.ListAgents

	mockAgentRepository        *repositories_mocks.MockAgentRepositoryInterface
	mockAgentAddressRepository *repositories_mocks.MockAgentAddressRepositoryInterface
}

func WhenListingAgentsBeforeEach(t *testing.T) *WhenListingAgentsTestingSuite {
	mockController := gomock.NewController(t, gomock.WithOverridableExpectations())

	mockAgentRepository := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	mockAgentAddressRepository := repositories_mocks.NewMockAgentAddressRepositoryInterface(mockController)

	sut := use_cases.NewListAgents(mockAgentRepository, mockAgentAddressRepository)

	return &WhenListingAgentsTestingSuite{
		sut:                        sut,
		mockAgentRepository:        mockAgentRepository,
		mockAgentAddressRepository: mockAgentAddressRepository,
	}
}

func TestWhenListingAgents(t *testing.T) {
	t.Parallel()

	t.Run("Given error happens while searching agents", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentsTestingSuite) {
			suite.mockAgentRepository.EXPECT().
				Search(gomock.Any()).
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenListingAgentsBeforeEach(t)
			initSuite(basics)

			req := &requests.ListAgentsRequest{}
			_, err := basics.sut.List(req)

			assert.NotNil(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given error happens while counting agents", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentsTestingSuite) {
			suite.mockAgentRepository.EXPECT().
				Search(gomock.Any()).
				Return([]*models.Agent{}, nil).
				Times(1)
			suite.mockAgentRepository.EXPECT().
				SearchCount(gomock.Any()).
				Return(uint(0), assert.AnError).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenListingAgentsBeforeEach(t)
			initSuite(basics)

			req := &requests.ListAgentsRequest{}
			_, err := basics.sut.List(req)

			assert.NotNil(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given agents are fetched successfully", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return agents with pagination info", func(t *testing.T) {
			t.Parallel()
			basics := WhenListingAgentsBeforeEach(t)

			agents := []*models.Agent{}
			for i := 0; i < 2; i++ {
				agent := model_builders.NewAgentBuilder().Build()
				agents = append(agents, agent)
			}

			basics.mockAgentRepository.EXPECT().
				Search(gomock.Any()).
				Return(agents, nil).
				Times(1)
			basics.mockAgentRepository.EXPECT().
				SearchCount(gomock.Any()).
				Return(uint(len(agents)), nil).
				Times(1)

			basics.mockAgentAddressRepository.EXPECT().
				GetAllByAgentIds(gomock.Any()).
				Return([]*models.AgentAddress{
					model_builders.NewAgentAddressBuilder().
						WithAgentId(agents[0].ID).
						WithAddress("0x2f0ddd93040c19863527b0553feceb696862ea7d").
						Build(),
				}, nil).
				Times(1)

			page := uint(1)
			elementsPerPage := uint(10)
			req := &requests.ListAgentsRequest{
				PaginateRequest: requests.PaginateRequest{
					Page:            &page,
					ElementsPerPage: &elementsPerPage,
				},
			}

			response, err := basics.sut.List(req)

			assert.Nil(t, err)
			assert.NotNil(t, response)
			assert.Len(t, response.Agents, len(agents))
			assert.Equal(t, page, response.Page)
			assert.Equal(t, elementsPerPage, response.ElementsPerPage)
			assert.Equal(t, uint(len(agents)), response.TotalElements)

			expectedByID := make(map[uuid.UUID]*models.Agent, len(agents))
			for _, agent := range agents {
				expectedByID[agent.ID] = agent
			}

			for _, fetchedAgent := range response.Agents {
				existingAgent := expectedByID[fetchedAgent.ID]
				if assert.NotNil(t, existingAgent) {
					assert.Equal(t, existingAgent.Name, fetchedAgent.Name)
					assert.Equal(t, existingAgent.Owner, fetchedAgent.Owner)
				}
			}
		})
	})

	t.Run("Given no pagination parameters are provided", func(t *testing.T) {
		t.Parallel()

		t.Run("Should use default pagination values", func(t *testing.T) {
			t.Parallel()
			basics := WhenListingAgentsBeforeEach(t)

			basics.mockAgentRepository.EXPECT().
				Search(gomock.Any()).
				DoAndReturn(func(params repositories.AgentSearchParamsWithPagination) ([]*models.Agent, error) {
					assert.Equal(t, lo.ToPtr(uint(1)), params.Page)
					assert.Equal(t, lo.ToPtr(uint(10)), params.ElementsPerPage)
					return []*models.Agent{}, nil
				}).
				Times(1)
			basics.mockAgentRepository.EXPECT().
				SearchCount(gomock.Any()).
				Return(uint(0), nil).
				Times(1)

			basics.mockAgentAddressRepository.EXPECT().
				GetAllByAgentIds(gomock.Any()).
				Return(nil, nil).
				Times(1)

			req := &requests.ListAgentsRequest{}

			response, err := basics.sut.List(req)

			assert.Nil(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, uint(1), response.Page)
			assert.Equal(t, uint(10), response.ElementsPerPage)
		})
	})

	t.Run("Given a minted filter is provided", func(t *testing.T) {
		t.Parallel()

		t.Run("Should forward it to the repository search and count", func(t *testing.T) {
			t.Parallel()
			basics := WhenListingAgentsBeforeEach(t)

			basics.mockAgentRepository.EXPECT().
				Search(gomock.Any()).
				DoAndReturn(func(params repositories.AgentSearchParamsWithPagination) ([]*models.Agent, error) {
					assert.Equal(t, lo.ToPtr(true), params.IsMinted)
					return []*models.Agent{}, nil
				}).
				Times(1)
			basics.mockAgentRepository.EXPECT().
				SearchCount(gomock.Any()).
				DoAndReturn(func(params repositories.AgentSearchParams) (uint, error) {
					assert.Equal(t, lo.ToPtr(true), params.IsMinted)
					return uint(0), nil
				}).
				Times(1)

			basics.mockAgentAddressRepository.EXPECT().
				GetAllByAgentIds(gomock.Any()).
				Return(nil, nil).
				Times(1)

			req := &requests.ListAgentsRequest{IsMinted: lo.ToPtr(true)}

			_, err := basics.sut.List(req)

			assert.Nil(t, err)
		})
	})

	t.Run("Given sort parameters are provided", func(t *testing.T) {
		t.Parallel()

		t.Run("Should forward them to the repository search", func(t *testing.T) {
			t.Parallel()
			basics := WhenListingAgentsBeforeEach(t)

			basics.mockAgentRepository.EXPECT().
				Search(gomock.Any()).
				DoAndReturn(func(params repositories.AgentSearchParamsWithPagination) ([]*models.Agent, error) {
					assert.Equal(t, lo.ToPtr("name"), params.SortBy)
					assert.Equal(t, lo.ToPtr("desc"), params.SortDirection)
					return []*models.Agent{}, nil
				}).
				Times(1)
			basics.mockAgentRepository.EXPECT().
				SearchCount(gomock.Any()).
				Return(uint(0), nil).
				Times(1)

			basics.mockAgentAddressRepository.EXPECT().
				GetAllByAgentIds(gomock.Any()).
				Return(nil, nil).
				Times(1)

			req := &requests.ListAgentsRequest{
				PaginateRequest: requests.PaginateRequest{
					SortDirection: lo.ToPtr(requests.SortDirectionDesc),
				},
				SortBy: lo.ToPtr(requests.ListAgentSortByName),
			}

			_, err := basics.sut.List(req)

			assert.Nil(t, err)
		})
	})
}
