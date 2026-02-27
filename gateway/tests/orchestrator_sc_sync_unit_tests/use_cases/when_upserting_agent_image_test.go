package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenUpsertingAgentImageTestingSuite struct {
	sut *use_cases.UpsertAgentImage

	agentRepositoryMock *repositories_mocks.MockAgentRepositoryInterface
}

func WhenUpsertingAgentImageBeforeEach(t *testing.T) *WhenUpsertingAgentImageTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)

	sut := use_cases.NewUpsertAgentImage(
		agentRepositoryMock,
	)

	return &WhenUpsertingAgentImageTestingSuite{
		sut: sut,

		agentRepositoryMock: agentRepositoryMock,
	}
}

func generateValidUpsertAgentImageRequest() *requests.UpsertAgentImageRequest {
	collectionAddress := common_tests.RandomEthAddress().Hex()
	return &requests.UpsertAgentImageRequest{
		CollectionTokenId: 1,
		CollectionAddress: collectionAddress,

		Key:      "default",
		NewValue: "azertyuiop",
	}
}

func TestWhenUpsertingAgentImage(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because collection address is not a valid Ethereum address", func(t *testing.T) {
			t.Parallel()

			suite := WhenUpsertingAgentImageBeforeEach(t)
			request := generateValidUpsertAgentImageRequest()
			request.CollectionAddress = "invalid_address"

			err := suite.sut.Upsert(request)

			assert.Error(t, err)
		})

		t.Run("Because collection token ID is not valid", func(t *testing.T) {
			t.Parallel()

			suite := WhenUpsertingAgentImageBeforeEach(t)
			request := generateValidUpsertAgentImageRequest()
			request.CollectionTokenId = 0

			err := suite.sut.Upsert(request)

			assert.Error(t, err)
		})

		t.Run("Because key is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenUpsertingAgentImageBeforeEach(t)
			request := generateValidUpsertAgentImageRequest()
			request.Key = ""

			err := suite.sut.Upsert(request)

			assert.Error(t, err)
		})

		t.Run("Because image value is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenUpsertingAgentImageBeforeEach(t)
			request := generateValidUpsertAgentImageRequest()
			request.NewValue = ""

			err := suite.sut.Upsert(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting agent by collection address and token ID", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenUpsertingAgentImageTestingSuite) {
			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenUpsertingAgentImageBeforeEach(t)
			initSuite(suite)

			request := generateValidUpsertAgentImageRequest()

			err := suite.sut.Upsert(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent does not exist for the provided collection address and token ID", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenUpsertingAgentImageTestingSuite) {
			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), gomock.Any()).
				Return(nil, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenUpsertingAgentImageBeforeEach(t)
			initSuite(suite)

			request := generateValidUpsertAgentImageRequest()

			err := suite.sut.Upsert(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while saving agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenUpsertingAgentImageTestingSuite) {
			agent := model_builders.NewAgentBuilder().Build()
			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), gomock.Any()).
				Return(agent, nil)

			suite.agentRepositoryMock.
				EXPECT().
				Save(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenUpsertingAgentImageBeforeEach(t)
			initSuite(suite)

			request := generateValidUpsertAgentImageRequest()

			err := suite.sut.Upsert(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent exists and image is successfully upserted", func(t *testing.T) {
		t.Parallel()

		t.Run("And image already exists", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenUpsertingAgentImageTestingSuite) *models.Agent {
				agent := model_builders.NewAgentBuilder().WithImages(map[string]string{
					"default": "old_cid",
				}).Build()
				suite.agentRepositoryMock.
					EXPECT().
					GetByCollectionAddressAndCollectionTokenId(gomock.Any(), gomock.Any()).
					Return(agent, nil)

				suite.agentRepositoryMock.
					EXPECT().
					Save(gomock.Any()).
					Return(agent, nil)

				return agent
			}

			t.Run("Should override existing image", func(t *testing.T) {
				t.Parallel()

				suite := WhenUpsertingAgentImageBeforeEach(t)
				agent := initSuite(suite)

				request := generateValidUpsertAgentImageRequest()
				request.NewValue = "new_cid"

				err := suite.sut.Upsert(request)

				assert.NoError(t, err)
				assert.Equal(t, "new_cid", agent.Images["default"])
			})
		})

		t.Run("And image does not exist", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenUpsertingAgentImageTestingSuite) *models.Agent {
				agent := model_builders.NewAgentBuilder().WithImages(map[string]string{}).
					Build()
				suite.agentRepositoryMock.
					EXPECT().
					GetByCollectionAddressAndCollectionTokenId(gomock.Any(), gomock.Any()).
					Return(agent, nil)

				suite.agentRepositoryMock.
					EXPECT().
					Save(gomock.Any()).
					Return(agent, nil)

				return agent
			}

			t.Run("Should add new image", func(t *testing.T) {
				t.Parallel()

				suite := WhenUpsertingAgentImageBeforeEach(t)
				agent := initSuite(suite)

				request := generateValidUpsertAgentImageRequest()
				request.NewValue = "new_cid"

				err := suite.sut.Upsert(request)

				assert.NoError(t, err)
				assert.Equal(t, "new_cid", agent.Images["default"])
			})
		})
	})
}
