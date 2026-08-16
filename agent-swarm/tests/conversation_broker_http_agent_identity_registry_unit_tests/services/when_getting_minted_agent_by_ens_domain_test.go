package services_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_agent_identity_registry/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_agent_identity_registry/services"
	clients_mocks "github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_http_agent_identity_registry_mocks/clients_mocks"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingMintedAgentByEnsDomainTestingSuite struct {
	sut *services.AgentRegistryRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenGettingMintedAgentByEnsDomainBeforeEach(t *testing.T) *WhenGettingMintedAgentByEnsDomainTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	return &WhenGettingMintedAgentByEnsDomainTestingSuite{
		sut:        services.NewAgentRegistryRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenGettingMintedAgentByEnsDomain(t *testing.T) {
	t.Parallel()

	t.Run("Given the agent exists", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingMintedAgentByEnsDomainTestingSuite) uuid.UUID {
			agentId := uuid.New()

			suite.mockClient.EXPECT().GetAgentEnsEnsDomainWithResponse(gomock.Any(), "ada.eth").Return(&clients.GetAgentEnsEnsDomainResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &clients.ResponsesAgentResponse{
					Id:        lo.ToPtr(agentId.String()),
					Name:      lo.ToPtr("Ada"),
					EnsDomain: lo.ToPtr("ada.eth"),
				},
			}, nil).Times(1)

			return agentId
		}

		t.Run("Should return the mapped agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingMintedAgentByEnsDomainBeforeEach(t)
			agentId := initSuite(suite)

			result, err := suite.sut.GetMintedAgentByEnsDomain("ada.eth")

			assert.NoError(t, err)
			if assert.NotNil(t, result) {
				assert.Equal(t, agentId, result.ID)
				assert.Equal(t, "ada.eth", result.EnsDomain)
			}
		})
	})

	t.Run("Given the agent does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingMintedAgentByEnsDomainTestingSuite) {
			suite.mockClient.EXPECT().GetAgentEnsEnsDomainWithResponse(gomock.Any(), "ada.eth").Return(&clients.GetAgentEnsEnsDomainResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusNotFound},
			}, nil).Times(1)
		}

		t.Run("Should return nil without error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingMintedAgentByEnsDomainBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.GetMintedAgentByEnsDomain("ada.eth")

			assert.NoError(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Given the registry returns an unexpected status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingMintedAgentByEnsDomainTestingSuite) {
			suite.mockClient.EXPECT().GetAgentEnsEnsDomainWithResponse(gomock.Any(), "ada.eth").Return(&clients.GetAgentEnsEnsDomainResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusInternalServerError},
				Body:         []byte("boom"),
			}, nil).Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingMintedAgentByEnsDomainBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.GetMintedAgentByEnsDomain("ada.eth")

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingMintedAgentByEnsDomainTestingSuite) {
			suite.mockClient.EXPECT().GetAgentEnsEnsDomainWithResponse(gomock.Any(), "ada.eth").Return(nil, errors.New("registry down")).Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingMintedAgentByEnsDomainBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.GetMintedAgentByEnsDomain("ada.eth")

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
}
