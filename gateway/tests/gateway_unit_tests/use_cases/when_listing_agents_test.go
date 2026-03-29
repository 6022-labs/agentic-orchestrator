package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenListingAgentsTestSuite struct {
	agentRegistryRequesterMock      *interfaces_mocks.MockAgentRegistryRequesterInterface
	llmConnectionRequesterMock      *interfaces_mocks.MockLLMConnectionRequesterInterface
	databaseConnectionRequesterMock *interfaces_mocks.MockDatabaseConnectionRequesterInterface

	sut *use_cases.ListAgents
}

func WhenListingAgentsTestSuiteBeforeEach(t *testing.T) *WhenListingAgentsTestSuite {
	// Create mocks
	mockController := gomock.NewController(t, gomock.WithOverridableExpectations())

	agentRegistryRequesterMock := interfaces_mocks.NewMockAgentRegistryRequesterInterface(mockController)
	llmConnectionRequesterMock := interfaces_mocks.NewMockLLMConnectionRequesterInterface(mockController)
	databaseConnectionRequesterMock := interfaces_mocks.NewMockDatabaseConnectionRequesterInterface(mockController)

	sut := use_cases.NewListAgents(
		agentRegistryRequesterMock,
		llmConnectionRequesterMock,
		databaseConnectionRequesterMock,
	)

	return &WhenListingAgentsTestSuite{
		sut: sut,

		agentRegistryRequesterMock:      agentRegistryRequesterMock,
		llmConnectionRequesterMock:      llmConnectionRequesterMock,
		databaseConnectionRequesterMock: databaseConnectionRequesterMock,
	}
}

func TestWhenListingAgents(t *testing.T) {
	t.Parallel()

	t.Run("Given error happens while fetching agents from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentsTestSuite) *requests.ListAgentsRequest {
			suite.agentRegistryRequesterMock.
				EXPECT().
				ListAgents(gomock.Any()).
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

	t.Run("Given agents are fetched successfully", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentsTestSuite) []*responses.AgentResponse {
			agents := []*responses.AgentResponse{}
			registryAgents := []*responses.AgentResponse{}

			for i := 0; i < 2; i++ {
				agent := model_builders.NewAgentBuilder().Build()
				agents = append(agents, agent)
				registryAgents = append(registryAgents, &responses.AgentResponse{
					ID:                agent.ID,
					Name:              agent.Name,
					Role:              agent.Role,
					Owner:             agent.Owner,
					Creator:           agent.Creator,
					Images:            agent.Images,
					CloneOf:           agent.CloneOf,
					Attributes:        agent.Attributes,
					EnsDomain:         agent.EnsDomain,
					IsInternal:        agent.IsInternal,
					WalletAddress:     agent.WalletAddress,
					ProposalId:        agent.ProposalId,
					CollectionTokenId: agent.CollectionTokenId,
					CollectionAddress: agent.CollectionAddress,
				})

			}

			suite.agentRegistryRequesterMock.
				EXPECT().
				ListAgents(gomock.Any()).
				Return(&responses.ListAgentsResponse{
					PaginateResponse: responses.PaginateResponse{
						Page:            1,
						ElementsPerPage: 10,
						TotalElements:   uint(len(registryAgents)),
					},
					Agents: registryAgents,
				}, nil)

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

			expectedByID := make(map[uuid.UUID]*responses.AgentResponse, len(existingAgents))
			for _, existingAgent := range existingAgents {
				expectedByID[existingAgent.ID] = existingAgent
			}

			for _, fetchedAgent := range response.Agents {
				existingAgent := expectedByID[fetchedAgent.ID]
				if assert.NotNil(t, existingAgent) {
					assert.Equal(t, existingAgent.Name, fetchedAgent.Name)
					assert.Equal(t, existingAgent.Owner, fetchedAgent.Owner)
					assert.Equal(t, existingAgent.WalletAddress, fetchedAgent.WalletAddress)
					assert.NotNil(t, fetchedAgent.LLMSetupStatus)
					assert.Equal(t, "llm-ok", *fetchedAgent.LLMSetupStatus)
					assert.NotNil(t, fetchedAgent.DatabaseSetupStatus)
					assert.Equal(t, "db-ok", *fetchedAgent.DatabaseSetupStatus)
				}
			}
		})
	})

	t.Run("Given setup status retrieval fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentsTestSuite) *requests.ListAgentsRequest {
			agent := model_builders.NewAgentBuilder().WithRole(enums.AgentRoleExpert).Build()
			suite.agentRegistryRequesterMock.
				EXPECT().
				ListAgents(gomock.Any()).
				Return(&responses.ListAgentsResponse{
					PaginateResponse: responses.PaginateResponse{
						Page:            1,
						ElementsPerPage: 10,
						TotalElements:   1,
					},
					Agents: []*responses.AgentResponse{
						{
							ID:                agent.ID,
							Name:              agent.Name,
							Role:              agent.Role,
							Owner:             agent.Owner,
							Creator:           agent.Creator,
							Images:            agent.Images,
							CloneOf:           agent.CloneOf,
							Attributes:        agent.Attributes,
							EnsDomain:         agent.EnsDomain,
							IsInternal:        agent.IsInternal,
							WalletAddress:     agent.WalletAddress,
							ProposalId:        agent.ProposalId,
							CollectionTokenId: agent.CollectionTokenId,
							CollectionAddress: agent.CollectionAddress,
						},
					},
				}, nil)

			suite.llmConnectionRequesterMock.
				EXPECT().
				GetSetupStatus(agent.ID).
				Return(nil, assert.AnError)
			suite.databaseConnectionRequesterMock.
				EXPECT().
				GetSetupStatus(agent.ID).
				Return(nil, nil).
				AnyTimes()

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
}
