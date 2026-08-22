package services_test

import (
	"context"

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

type WhenGettingAgentTestingSuite struct {
	sut *services.AgentRegistryRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenGettingAgentBeforeEach(t *testing.T) *WhenGettingAgentTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	return &WhenGettingAgentTestingSuite{
		sut:        services.NewAgentRegistryRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenGettingAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given the agent exists", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentTestingSuite) uuid.UUID {
			agentId := uuid.New()

			suite.mockClient.EXPECT().
				GetAgentAgentCollectionAddressAgentCollectionTokenIdWithResponse(gomock.Any(), "0xabc", 7).
				Return(&clients.GetAgentAgentCollectionAddressAgentCollectionTokenIdResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusOK},
					JSON200: &clients.ResponsesAgentResponse{
						Id:                lo.ToPtr(agentId.String()),
						Name:              lo.ToPtr("Ada"),
						EnsDomain:         lo.ToPtr("ada.eth"),
						Attributes:        lo.ToPtr(map[string]string{"role": "clone"}),
						CollectionAddress: lo.ToPtr("0xabc"),
						CollectionTokenId: lo.ToPtr(7),
					},
				}, nil).
				Times(1)

			return agentId
		}

		t.Run("Should return the mapped agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentBeforeEach(t)
			agentId := initSuite(suite)

			result, err := suite.sut.GetAgent(context.Background(), "0xabc", 7)

			assert.NoError(t, err)
			if assert.NotNil(t, result) {
				assert.Equal(t, agentId, result.ID)
				assert.Equal(t, "Ada", result.Name)
				assert.Equal(t, "ada.eth", result.EnsDomain)
				assert.Equal(t, map[string]string{"role": "clone"}, result.Attributes)
				assert.Equal(t, "0xabc", result.CollectionAddress)
				assert.Equal(t, uint64(7), *result.CollectionTokenId)
			}
		})
	})

	t.Run("Given the agent does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentTestingSuite) {
			suite.mockClient.EXPECT().
				GetAgentAgentCollectionAddressAgentCollectionTokenIdWithResponse(gomock.Any(), "0xabc", 7).
				Return(&clients.GetAgentAgentCollectionAddressAgentCollectionTokenIdResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusNotFound},
				}, nil).
				Times(1)
		}

		t.Run("Should return nil without error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.GetAgent(context.Background(), "0xabc", 7)

			assert.NoError(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Given the registry returns an unexpected status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentTestingSuite) {
			suite.mockClient.EXPECT().
				GetAgentAgentCollectionAddressAgentCollectionTokenIdWithResponse(gomock.Any(), "0xabc", 7).
				Return(&clients.GetAgentAgentCollectionAddressAgentCollectionTokenIdResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusInternalServerError},
					Body:         []byte("boom"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.GetAgent(context.Background(), "0xabc", 7)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentTestingSuite) {
			suite.mockClient.EXPECT().
				GetAgentAgentCollectionAddressAgentCollectionTokenIdWithResponse(gomock.Any(), "0xabc", 7).
				Return(nil, errors.New("registry down")).
				Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.GetAgent(context.Background(), "0xabc", 7)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
}
