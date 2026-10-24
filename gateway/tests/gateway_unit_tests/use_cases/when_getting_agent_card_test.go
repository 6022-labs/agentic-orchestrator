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

type WhenGettingAgentCardTestingSuite struct {
	sut *use_cases.GetAgentCard

	mockAgentRegistryRequester      *interfaces_mocks.MockAgentRegistryRequesterInterface
	mockConversationBrokerRequester *interfaces_mocks.MockConversationBrokerRequesterInterface
}

func WhenGettingAgentCardBeforeEach(t *testing.T) *WhenGettingAgentCardTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentRegistryRequester := interfaces_mocks.NewMockAgentRegistryRequesterInterface(mockController)
	mockConversationBrokerRequester := interfaces_mocks.NewMockConversationBrokerRequesterInterface(mockController)

	sut := use_cases.NewGetAgentCard(mockAgentRegistryRequester, mockConversationBrokerRequester)

	return &WhenGettingAgentCardTestingSuite{
		sut: sut,

		mockAgentRegistryRequester:      mockAgentRegistryRequester,
		mockConversationBrokerRequester: mockConversationBrokerRequester,
	}
}

func generateValidGetAgentCardRequest() *requests.AgentByCollectionTokenIdRequest {
	return &requests.AgentByCollectionTokenIdRequest{
		AgentCollectionAddress: "0x39DD75a025145770a8b8278Cc8BA705116443Ea7",
		AgentCollectionTokenId: 1,
	}
}

func TestWhenGettingAgentCard(t *testing.T) {
	t.Parallel()

	t.Run("Given a collection address without the 0x prefix", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error without calling the registry", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentCardBeforeEach(t)
			request := generateValidGetAgentCardRequest()
			request.AgentCollectionAddress = "39DD75a025145770a8b8278Cc8BA705116443Ea7"

			response, err := suite.sut.Get(context.Background(), request)

			assert.Nil(t, response)
			assert.IsType(t, &gw_errors.ValidationError{}, err)
		})
	})

	t.Run("Given the registry fails", func(t *testing.T) {
		t.Parallel()

		registryErr := errors.New("boom")

		initSuite := func(suite *WhenGettingAgentCardTestingSuite) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), "0x39dd75a025145770a8b8278cc8ba705116443ea7", uint64(1)).
				Return(nil, registryErr).
				Times(1)
		}

		t.Run("Should propagate the wrapped error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentCardBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Get(context.Background(), generateValidGetAgentCardRequest())

			assert.Nil(t, response)
			assert.ErrorIs(t, err, registryErr)
		})
	})

	t.Run("Given the agent does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentCardTestingSuite) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), "0x39dd75a025145770a8b8278cc8ba705116443ea7", uint64(1)).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentCardBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Get(context.Background(), generateValidGetAgentCardRequest())

			assert.IsType(t, &gw_errors.NotFoundError{}, err)
		})
	})

	t.Run("Given the agent has no ens domain", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentCardTestingSuite) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), "0x39dd75a025145770a8b8278cc8ba705116443ea7", uint64(1)).
				Return(&responses.AgentResponse{}, nil).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentCardBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Get(context.Background(), generateValidGetAgentCardRequest())

			assert.IsType(t, &gw_errors.NotFoundError{}, err)
		})
	})

	t.Run("Given the conversation broker fails", func(t *testing.T) {
		t.Parallel()

		brokerErr := errors.New("broker down")

		initSuite := func(suite *WhenGettingAgentCardTestingSuite) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), "0x39dd75a025145770a8b8278cc8ba705116443ea7", uint64(1)).
				Return(&responses.AgentResponse{EnsDomain: "ada.experts.80002.6022.eth"}, nil).
				Times(1)
			suite.mockConversationBrokerRequester.EXPECT().
				GetAgentCard(gomock.Any(), "ada.experts.80002.6022.eth").
				Return(nil, brokerErr).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentCardBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Get(context.Background(), generateValidGetAgentCardRequest())

			assert.Nil(t, response)
			assert.ErrorIs(t, err, brokerErr)
		})
	})

	t.Run("Given a checksummed collection address and an agent with an ens domain", func(t *testing.T) {
		t.Parallel()

		card := &responses.AgentCardResponse{Name: "Ada", EnsDomain: "ada.experts.80002.6022.eth"}

		initSuite := func(suite *WhenGettingAgentCardTestingSuite) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), "0x39dd75a025145770a8b8278cc8ba705116443ea7", uint64(1)).
				Return(&responses.AgentResponse{EnsDomain: "ada.experts.80002.6022.eth"}, nil).
				Times(1)
			suite.mockConversationBrokerRequester.EXPECT().
				GetAgentCard(gomock.Any(), "ada.experts.80002.6022.eth").
				Return(card, nil).
				Times(1)
		}

		t.Run("Should look up the lowercased address and return the card", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentCardBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Get(context.Background(), generateValidGetAgentCardRequest())

			assert.NoError(t, err)
			assert.Same(t, card, response)
		})
	})
}
