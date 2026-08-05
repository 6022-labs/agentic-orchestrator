package use_cases_test

import (
	"errors"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenListingConversationsTestingSuite struct {
	sut *use_cases.ListConversations

	mockConversationRepository *repositories_mocks.MockConversationRepositoryInterface
}

func WhenListingConversationsBeforeEach(t *testing.T) *WhenListingConversationsTestingSuite {
	mockController := gomock.NewController(t)

	mockConversationRepository := repositories_mocks.NewMockConversationRepositoryInterface(mockController)

	sut := use_cases.NewListConversations(mockConversationRepository)

	return &WhenListingConversationsTestingSuite{
		sut: sut,

		mockConversationRepository: mockConversationRepository,
	}
}

func TestWhenListingConversations(t *testing.T) {
	t.Parallel()

	t.Run("Given no pagination is provided", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingConversationsTestingSuite) {
			suite.mockConversationRepository.EXPECT().
				List(uint(1), uint(10)).
				Return(nil, nil).
				Times(1)
			suite.mockConversationRepository.EXPECT().
				Count().
				Return(uint(0), nil).
				Times(1)
		}

		t.Run("Should default to page 1 with 10 elements", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingConversationsBeforeEach(t)
			initSuite(suite)

			resp, err := suite.sut.List(&requests.ListConversationsRequest{})

			assert.NoError(t, err)
			assert.Equal(t, uint(1), resp.Page)
			assert.Equal(t, uint(10), resp.ElementsPerPage)
		})
	})

	t.Run("Given conversations exist", func(t *testing.T) {
		t.Parallel()

		conversation := model_builders.NewConversationBuilder().Build()

		initSuite := func(suite *WhenListingConversationsTestingSuite) {
			suite.mockConversationRepository.EXPECT().
				List(uint(2), uint(5)).
				Return([]*models.Conversation{conversation}, nil).
				Times(1)
			suite.mockConversationRepository.EXPECT().
				Count().
				Return(uint(6), nil).
				Times(1)
		}

		t.Run("Should return the page and the total count", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingConversationsBeforeEach(t)
			initSuite(suite)

			resp, err := suite.sut.List(&requests.ListConversationsRequest{
				Page:            lo.ToPtr(uint(2)),
				ElementsPerPage: lo.ToPtr(uint(5)),
			})

			assert.NoError(t, err)
			assert.Len(t, resp.Conversations, 1)
			assert.Equal(t, conversation.ID, resp.Conversations[0].ID)
			assert.Equal(t, uint(6), resp.TotalElements)
		})
	})

	t.Run("Given listing conversations fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingConversationsTestingSuite) {
			suite.mockConversationRepository.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingConversationsBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.List(&requests.ListConversationsRequest{})

			assert.Error(t, err)
		})
	})

	t.Run("Given counting conversations fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingConversationsTestingSuite) {
			suite.mockConversationRepository.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(nil, nil).
				Times(1)
			suite.mockConversationRepository.EXPECT().
				Count().
				Return(uint(0), errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingConversationsBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.List(&requests.ListConversationsRequest{})

			assert.Error(t, err)
		})
	})
}
