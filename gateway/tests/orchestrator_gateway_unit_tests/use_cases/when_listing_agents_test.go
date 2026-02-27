package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenListingAgentsTestSuite struct {
	agentRepositoryMock             *repositories_mocks.MockAgentRepositoryInterface
	llmConnectionRequesterMock      *interfaces_mocks.MockLLMConnectionRequesterInterface
	databaseConnectionRequesterMock *interfaces_mocks.MockDatabaseConnectionRequesterInterface

	sut *use_cases.ListAgents
}

func WhenListingAgentsTestSuiteBeforeEach(t *testing.T) *WhenListingAgentsTestSuite {
	// Create mocks
	mockController := gomock.NewController(t, gomock.WithOverridableExpectations())

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	llmConnectionRequesterMock := interfaces_mocks.NewMockLLMConnectionRequesterInterface(mockController)
	databaseConnectionRequesterMock := interfaces_mocks.NewMockDatabaseConnectionRequesterInterface(mockController)

	sut := use_cases.NewListAgents(
		agentRepositoryMock,
		llmConnectionRequesterMock,
		databaseConnectionRequesterMock,
	)

	return &WhenListingAgentsTestSuite{
		sut: sut,

		agentRepositoryMock:             agentRepositoryMock,
		llmConnectionRequesterMock:      llmConnectionRequesterMock,
		databaseConnectionRequesterMock: databaseConnectionRequesterMock,
	}
}

func TestWhenListingAgents(t *testing.T) {
	t.Parallel()

	t.Run("Given error happens while fetching agents from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentsTestSuite) *requests.ListAgentsRequest {
			suite.agentRepositoryMock.
				EXPECT().
				Search(gomock.Any()).
				Return(nil, assert.AnError)

			return &requests.ListAgentsRequest{}
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			basics := WhenListingAgentsTestSuiteBeforeEach(t)
			req := initSuite(basics)

			_, err := basics.sut.List(req)

			assert.NotNil(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given error happens while fetching count of agents from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentsTestSuite) *requests.ListAgentsRequest {
			suite.agentRepositoryMock.
				EXPECT().
				Search(gomock.Any()).
				Return(nil, nil)
			suite.agentRepositoryMock.
				EXPECT().
				SearchCount(gomock.Any()).
				Return(uint(0), assert.AnError)

			return &requests.ListAgentsRequest{}
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			basics := WhenListingAgentsTestSuiteBeforeEach(t)
			req := initSuite(basics)

			_, err := basics.sut.List(req)

			assert.NotNil(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given agents are fetched successfully", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentsTestSuite) []*models.Agent {
			agents := []*models.Agent{}

			for i := 0; i < 2; i++ {
				agent := model_builders.NewAgentBuilder().Build()
				agents = append(agents, agent)

			}

			suite.agentRepositoryMock.
				EXPECT().
				Search(gomock.Any()).
				Return(agents, nil)
			suite.agentRepositoryMock.
				EXPECT().
				SearchCount(gomock.Any()).
				Return(uint(len(agents)), nil)

			llmStatus := "llm-ok"
			llmConnectionStatus := interfaces.LLMConnectionStatus{
				Status:      &llmStatus,
				IsCompleted: true,
			}
			dbStatus := "db-ok"
			dbConnectionStatus := interfaces.DatabaseConnectionStatus{
				Status:      &dbStatus,
				IsCompleted: true,
			}

			suite.llmConnectionRequesterMock.
				EXPECT().
				GetSetupStatus(gomock.Any()).
				Return(&llmConnectionStatus, nil).
				AnyTimes()
			suite.databaseConnectionRequesterMock.
				EXPECT().
				GetSetupStatus(gomock.Any()).
				Return(&dbConnectionStatus, nil).
				AnyTimes()

			return agents
		}

		t.Run("Should return existing agents", func(t *testing.T) {
			t.Parallel()

			basics := WhenListingAgentsTestSuiteBeforeEach(t)
			existingAgents := initSuite(basics)

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
			assert.Len(t, response.Agents, len(existingAgents))
			assert.Equal(t, page, response.Page)
			assert.Equal(t, elementsPerPage, response.ElementsPerPage)
			assert.Equal(t, uint(len(existingAgents)), response.TotalElements)

			for i, existingAgent := range existingAgents {
				fetchedAgent := response.Agents[i]
				assert.Equal(t, existingAgent.ID, fetchedAgent.ID)
				assert.Equal(t, existingAgent.Name, fetchedAgent.Name)
				assert.Equal(t, existingAgent.Owner, fetchedAgent.Owner)
				assert.Equal(t, existingAgent.WalletAddress, fetchedAgent.WalletAddress)
				assert.NotNil(t, fetchedAgent.LLMSetupStatus)
				assert.Equal(t, "llm-ok", *fetchedAgent.LLMSetupStatus)
				assert.NotNil(t, fetchedAgent.DatabaseSetupStatus)
				assert.Equal(t, "db-ok", *fetchedAgent.DatabaseSetupStatus)
			}
		})
	})
}
