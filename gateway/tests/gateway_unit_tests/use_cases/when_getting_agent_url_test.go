package use_cases_test

import (
	"context"

	"errors"
	"testing"

	gw_errors "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingAgentUrlTestingSuite struct {
	sut *use_cases.GetAgentUrl

	mockAgentRegistryRequester      *interfaces_mocks.MockAgentRegistryRequesterInterface
	mockConversationBrokerRequester *interfaces_mocks.MockConversationBrokerRequesterInterface
}

func WhenGettingAgentUrlBeforeEach(t *testing.T) *WhenGettingAgentUrlTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentRegistryRequester := interfaces_mocks.NewMockAgentRegistryRequesterInterface(mockController)
	mockConversationBrokerRequester := interfaces_mocks.NewMockConversationBrokerRequesterInterface(mockController)

	sut := use_cases.NewGetAgentUrl(mockAgentRegistryRequester, mockConversationBrokerRequester)

	return &WhenGettingAgentUrlTestingSuite{
		sut: sut,

		mockAgentRegistryRequester:      mockAgentRegistryRequester,
		mockConversationBrokerRequester: mockConversationBrokerRequester,
	}
}

func generateValidGetAgentUrlRequest() *requests.AgentByCollectionTokenIdRequest {
	return &requests.AgentByCollectionTokenIdRequest{
		AgentCollectionAddress: "0x39DD75a025145770a8b8278Cc8BA705116443Ea7",
		AgentCollectionTokenId: 1,
	}
}

func TestWhenGettingAgentUrl(t *testing.T) {
	t.Parallel()

	t.Run("Given a collection address without the 0x prefix", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error without calling the registry", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentUrlBeforeEach(t)
			request := generateValidGetAgentUrlRequest()
			request.AgentCollectionAddress = "39DD75a025145770a8b8278Cc8BA705116443Ea7"

			response, err := suite.sut.Get(context.Background(), request)

			assert.Nil(t, response)
			assert.IsType(t, &gw_errors.ValidationError{}, err)
		})
	})

	t.Run("Given a checksummed collection address", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentUrlTestingSuite) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), "0x39dd75a025145770a8b8278cc8ba705116443ea7", uint64(1)).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should look up the lowercased address", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentUrlBeforeEach(t)
			initSuite(suite)
			request := generateValidGetAgentUrlRequest()

			_, _ = suite.sut.Get(context.Background(), request)

			assert.Equal(t, "0x39dd75a025145770a8b8278cc8ba705116443ea7", request.AgentCollectionAddress)
		})
	})

	t.Run("Given the registry fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentUrlTestingSuite) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), "0x39dd75a025145770a8b8278cc8ba705116443ea7", uint64(1)).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentUrlBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Get(context.Background(), generateValidGetAgentUrlRequest())

			assert.Error(t, err)
		})
	})

	t.Run("Given the agent does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentUrlTestingSuite) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), "0x39dd75a025145770a8b8278cc8ba705116443ea7", uint64(1)).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentUrlBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Get(context.Background(), generateValidGetAgentUrlRequest())

			assert.IsType(t, &gw_errors.NotFoundError{}, err)
		})
	})

	t.Run("Given the agent has no ens domain", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentUrlTestingSuite) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), "0x39dd75a025145770a8b8278cc8ba705116443ea7", uint64(1)).
				Return(&responses.AgentResponse{}, nil).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentUrlBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Get(context.Background(), generateValidGetAgentUrlRequest())

			assert.IsType(t, &gw_errors.NotFoundError{}, err)
		})
	})

	t.Run("Given the agent has an ens domain", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentUrlTestingSuite) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), "0x39dd75a025145770a8b8278cc8ba705116443ea7", uint64(1)).
				Return(&responses.AgentResponse{EnsDomain: "ada.experts.80002.6022.eth"}, nil).
				Times(1)
			suite.mockConversationBrokerRequester.EXPECT().
				GetAgentUrl(gomock.Any(), "ada.experts.80002.6022.eth").
				Return(&responses.AgentUrlResponse{Url: "https://ada.example"}, nil).
				Times(1)
		}

		t.Run("Should return the resolved url", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentUrlBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Get(context.Background(), generateValidGetAgentUrlRequest())

			assert.NoError(t, err)
			assert.Equal(t, "https://ada.example", response.Url)
		})
	})
}
