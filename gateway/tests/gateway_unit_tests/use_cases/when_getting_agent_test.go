package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingAgentTestingSuite struct {
	sut *use_cases.GetAgent

	agentRegistryRequesterMock      *interfaces_mocks.MockAgentRegistryRequesterInterface
	llmConnectionRequesterMock      *interfaces_mocks.MockLLMConnectionRequesterInterface
	databaseConnectionRequesterMock *interfaces_mocks.MockDatabaseConnectionRequesterInterface
}

func WhenGettingAgentTestingSuiteBeforeEach(t *testing.T) *WhenGettingAgentTestingSuite {
	mockController := gomock.NewController(t)

	agentRegistryRequesterMock := interfaces_mocks.NewMockAgentRegistryRequesterInterface(mockController)
	llmConnectionRequesterMock := interfaces_mocks.NewMockLLMConnectionRequesterInterface(mockController)
	databaseConnectionRequesterMock := interfaces_mocks.NewMockDatabaseConnectionRequesterInterface(mockController)

	sut := use_cases.NewGetAgent(
		agentRegistryRequesterMock,
		llmConnectionRequesterMock,
		databaseConnectionRequesterMock,
	)

	return &WhenGettingAgentTestingSuite{
		sut:                             sut,
		agentRegistryRequesterMock:      agentRegistryRequesterMock,
		llmConnectionRequesterMock:      llmConnectionRequesterMock,
		databaseConnectionRequesterMock: databaseConnectionRequesterMock,
	}
}

func TestWhenGettingAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given error happens while fetching agent from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentTestingSuite) {
			suite.agentRegistryRequesterMock.EXPECT().GetAgent(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenGettingAgentTestingSuiteBeforeEach(t)
			initSuite(basics)

			req := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()
			resp, err := basics.sut.Get(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given agent is not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentTestingSuite) {
			suite.agentRegistryRequesterMock.EXPECT().GetAgent(gomock.Any(), gomock.Any()).Return(nil, nil)
		}

		t.Run("Should return nil response and nil error", func(t *testing.T) {
			t.Parallel()
			basics := WhenGettingAgentTestingSuiteBeforeEach(t)
			initSuite(basics)

			req := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()
			resp, err := basics.sut.Get(req)

			assert.Nil(t, resp)
			assert.Nil(t, err)
		})
	})

	t.Run("Given agent is fetched successfully", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentTestingSuite) *responses.AgentResponse {
			agent := model_builders.NewAgentBuilder().Build()
			suite.agentRegistryRequesterMock.EXPECT().GetAgent(gomock.Any(), gomock.Any()).Return(&responses.AgentResponse{
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
			suite.llmConnectionRequesterMock.EXPECT().GetSetupStatus(agent.ID).Return(&llmConnectionStatus, nil)
			suite.databaseConnectionRequesterMock.EXPECT().GetSetupStatus(agent.ID).Return(&dbConnectionStatus, nil)
			return agent
		}

		t.Run("Should return agent response with setup statuses", func(t *testing.T) {
			t.Parallel()
			basics := WhenGettingAgentTestingSuiteBeforeEach(t)
			agent := initSuite(basics)

			req := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()
			resp, err := basics.sut.Get(req)

			assert.Nil(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, agent.ID, resp.ID)
			assert.Equal(t, agent.Name, resp.Name)
			assert.Equal(t, agent.Owner, resp.Owner)
			assert.Equal(t, agent.WalletAddress, resp.WalletAddress)
			assert.NotNil(t, resp.LLMSetupStatus)
			assert.Equal(t, "llm-ok", *resp.LLMSetupStatus)
			assert.NotNil(t, resp.DatabaseSetupStatus)
			assert.Equal(t, "db-ok", *resp.DatabaseSetupStatus)
		})
	})

	t.Run("Given error happens while fetching setup statuses", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentTestingSuite) *responses.AgentResponse {
			agent := model_builders.NewAgentBuilder().Build()
			suite.agentRegistryRequesterMock.EXPECT().GetAgent(gomock.Any(), gomock.Any()).Return(&responses.AgentResponse{
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
			}, nil)
			suite.llmConnectionRequesterMock.EXPECT().GetSetupStatus(agent.ID).Return(nil, assert.AnError)
			suite.databaseConnectionRequesterMock.EXPECT().GetSetupStatus(agent.ID).Return(nil, nil).AnyTimes()
			return agent
		}

		t.Run("Should return error from setup status", func(t *testing.T) {
			t.Parallel()
			basics := WhenGettingAgentTestingSuiteBeforeEach(t)
			_ = initSuite(basics)

			req := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()
			resp, err := basics.sut.Get(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})
}
