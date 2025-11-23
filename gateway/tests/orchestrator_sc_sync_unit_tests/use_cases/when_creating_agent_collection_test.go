package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_sc_sync_mocks/services_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenCreatingAgentCollectionTestSuite struct {
	sut *use_cases.CreateAgentCollection

	agentCollectionSynchronizerMock *services_mocks.MockAgentCollectionSynchronizerInterface
}

func WhenCreatingAgentCollectionBeforeEach(t *testing.T) *WhenCreatingAgentCollectionTestSuite {
	mockController := gomock.NewController(t)

	agentCollectionSynchronizerMock := services_mocks.NewMockAgentCollectionSynchronizerInterface(mockController)

	sut := use_cases.NewCreateAgentCollection(agentCollectionSynchronizerMock)

	return &WhenCreatingAgentCollectionTestSuite{
		sut: sut,

		agentCollectionSynchronizerMock: agentCollectionSynchronizerMock,
	}
}

func generateValidCreateAgentCollectionRequest() *requests.CreateAgentCollectionRequest {
	return &requests.CreateAgentCollectionRequest{
		CollectionAddress: common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678").Hex(),
		CollectionName:    "Test Collection",
		CollectionSymbol:  "TST",
	}
}

func TestWhenCreatingAgentCollection(t *testing.T) {
	t.Parallel()

	t.Run("Given synchronizer returns error", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentCollectionBeforeEach(t)
			request := generateValidCreateAgentCollectionRequest()

			suite.agentCollectionSynchronizerMock.EXPECT().Sync(gomock.Any(), gomock.Any(), gomock.Any()).Return(assert.AnError)

			err := suite.sut.Create(request)
			assert.Error(t, err)
		})
	})

	t.Run("Given synchronizer does not return error", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return err", func(t *testing.T) {
			suite := WhenCreatingAgentCollectionBeforeEach(t)
			request := generateValidCreateAgentCollectionRequest()

			suite.agentCollectionSynchronizerMock.EXPECT().Sync(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			err := suite.sut.Create(request)
			assert.NoError(t, err)
		})
	})
}
