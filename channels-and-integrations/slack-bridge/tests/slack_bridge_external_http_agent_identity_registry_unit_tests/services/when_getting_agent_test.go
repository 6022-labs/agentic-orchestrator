package services_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_http_agent_identity_registry/clients"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_http_agent_identity_registry/services"
	clients_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_external_http_agent_identity_registry_mocks/clients_mocks"
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

	t.Run("Given the agent exists with a role attribute", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentTestingSuite) {
			suite.mockClient.EXPECT().GetAgentAgentCollectionAddressAgentCollectionTokenIdWithResponse(gomock.Any(), "0xabc", 7).Return(&clients.GetAgentAgentCollectionAddressAgentCollectionTokenIdResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &clients.ResponsesAgentResponse{
					Owner:             lo.ToPtr("0xowner"),
					EnsDomain:         lo.ToPtr("ada.eth"),
					Images:            lo.ToPtr(map[string]string{"icon": "cid123"}),
					Attributes:        lo.ToPtr(map[string]string{"role": "clone"}),
					CollectionAddress: lo.ToPtr("0xabc"),
					CollectionTokenId: lo.ToPtr(7),
				},
			}, nil).Times(1)
		}

		t.Run("Should return the mapped slack agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.GetAgent("0xabc", 7)

			assert.NoError(t, err)
			if assert.NotNil(t, result) {
				assert.Equal(t, enums.AgentRoleClone, result.Role)
				assert.Equal(t, "0xowner", result.OwnerWallet)
				assert.Equal(t, "ada.eth", result.EnsDomain)
				assert.Equal(t, uint64(7), result.CollectionTokenId)
				assert.Equal(t, "0xabc", result.CollectionAddress)
				assert.Equal(t, "cid123", *result.IconImageCid)
			}
		})
	})

	t.Run("Given the agent has no role attribute", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentTestingSuite) {
			suite.mockClient.EXPECT().GetAgentAgentCollectionAddressAgentCollectionTokenIdWithResponse(gomock.Any(), "0xabc", 7).Return(&clients.GetAgentAgentCollectionAddressAgentCollectionTokenIdResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &clients.ResponsesAgentResponse{
					Attributes: lo.ToPtr(map[string]string{}),
				},
			}, nil).Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.GetAgent("0xabc", 7)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Given the agent does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentTestingSuite) {
			suite.mockClient.EXPECT().GetAgentAgentCollectionAddressAgentCollectionTokenIdWithResponse(gomock.Any(), "0xabc", 7).Return(&clients.GetAgentAgentCollectionAddressAgentCollectionTokenIdResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusNotFound},
				Body:         []byte("not found"),
			}, nil).Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.GetAgent("0xabc", 7)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentTestingSuite) {
			suite.mockClient.EXPECT().GetAgentAgentCollectionAddressAgentCollectionTokenIdWithResponse(gomock.Any(), "0xabc", 7).Return(nil, errors.New("registry down")).Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.GetAgent("0xabc", 7)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
}
