package use_cases_test

import (
	"context"
	"strings"

	"errors"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type WhenGettingAgentPromptTemplateTestingSuite struct {
	sut *use_cases.GetAgentPromptTemplate

	mockAgentPromptTemplateRepository *repositories_mocks.MockAgentPromptTemplateRepositoryInterface
}

func WhenGettingAgentPromptTemplateBeforeEach(t *testing.T) *WhenGettingAgentPromptTemplateTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentPromptTemplateRepository := repositories_mocks.NewMockAgentPromptTemplateRepositoryInterface(
		mockController,
	)

	sut := use_cases.NewGetAgentPromptTemplate(mockAgentPromptTemplateRepository)

	return &WhenGettingAgentPromptTemplateTestingSuite{
		sut: sut,

		mockAgentPromptTemplateRepository: mockAgentPromptTemplateRepository,
	}
}

func generateValidGetAgentPromptTemplateRequest() *requests.GetAgentPromptTemplateRequest {
	return &requests.GetAgentPromptTemplateRequest{
		CollectionAddress: "0xabc00000000000000000000000000000000000de",
		CollectionTokenId: 1,
		Key:               "clone-system",
	}
}

func TestWhenGettingAgentPromptTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because collectionAddress has no 0x prefix", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			request := generateValidGetAgentPromptTemplateRequest()
			request.CollectionAddress = strings.TrimPrefix(request.CollectionAddress, "0x")

			_, err := suite.sut.Get(context.Background(), request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given an uppercase key", func(t *testing.T) {
		t.Parallel()

		template := model_builders.NewAgentPromptTemplateBuilder().
			WithCollectionAddress("0xabc00000000000000000000000000000000000de").
			WithCollectionTokenId(1).
			WithKey(enums.AgentPromptTemplateKeyCloneSystem).
			Build()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite) {
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(gomock.Any(), gomock.Any(), uint64(1), enums.AgentPromptTemplateKeyCloneSystem).
				Return(template, nil).
				Times(1)
		}

		t.Run("Should look the template up by the canonical key", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)
			request := generateValidGetAgentPromptTemplateRequest()
			request.Key = "CLONE-SYSTEM"

			got, err := suite.sut.Get(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, template.ID, got.ID)
		})
	})

	t.Run("Given a checksummed collection address", func(t *testing.T) {
		t.Parallel()

		template := model_builders.NewAgentPromptTemplateBuilder().
			WithCollectionAddress("0xabc00000000000000000000000000000000000de").
			WithCollectionTokenId(1).
			Build()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite) {
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(gomock.Any(), "0xabc00000000000000000000000000000000000de", uint64(1), template.Key).
				Return(template, nil).
				Times(1)
		}

		t.Run("Should look the template up by the lowercased address", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)
			request := generateValidGetAgentPromptTemplateRequest()
			request.CollectionAddress = "0xAbC00000000000000000000000000000000000dE"

			got, err := suite.sut.Get(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, template.ID, got.ID)
		})
	})

	t.Run("Given repository fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite) {
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Get(context.Background(), generateValidGetAgentPromptTemplateRequest())

			assert.Error(t, err)
		})
	})

	t.Run("Given template does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite) {
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Get(context.Background(), generateValidGetAgentPromptTemplateRequest())

			var notFoundErr *cb_errors.NotFoundError
			require.ErrorAs(t, err, &notFoundErr)
			assert.Equal(t, "prompt_template_not_found", notFoundErr.Code)
		})
	})

	t.Run("Given template exists", func(t *testing.T) {
		t.Parallel()

		template := model_builders.NewAgentPromptTemplateBuilder().
			WithCollectionAddress("0xabc00000000000000000000000000000000000de").
			WithCollectionTokenId(1).
			Build()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite) {
			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByAgentAndKey(gomock.Any(), "0xabc00000000000000000000000000000000000de", uint64(1), template.Key).
				Return(template, nil).
				Times(1)
		}

		t.Run("Should return the template", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			got, err := suite.sut.Get(context.Background(), generateValidGetAgentPromptTemplateRequest())

			assert.NoError(t, err)
			assert.Equal(t, template.ID, got.ID)
		})
	})
}
