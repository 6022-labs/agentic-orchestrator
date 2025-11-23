package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_sc_sync_mocks/services_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenFinishingAgentCreationTestSuite struct {
	agentSynchronizerMock *services_mocks.MockAgentSynchronizerInterface

	sut use_cases.FinishAgentCreationInterface
}

func WhenFinishingAgentCreationBeforeEach(t *testing.T) *WhenFinishingAgentCreationTestSuite {
	mockController := gomock.NewController(t)

	agentSynchronizerMock := services_mocks.NewMockAgentSynchronizerInterface(mockController)

	sut := use_cases.NewFinishAgentCreation(
		agentSynchronizerMock,
	)

	return &WhenFinishingAgentCreationTestSuite{
		sut: sut,

		agentSynchronizerMock: agentSynchronizerMock,
	}
}

func generateValidFinishAgentCreationRequest() *requests.FinishAgentCreationRequest {
	return &requests.FinishAgentCreationRequest{
		CollectionTokenId: 1,
		CollectionAddress: "0x2332cc64591301E10D08f503A75B05D0903D6fe2",
		Signer:            "0x4caBA841628EBDB4101aDB8884611c6E55e46147",
	}
}

func TestWhenFinishingAgentCreation(t *testing.T) {
	t.Parallel()

	t.Run("Given synchronizer returns error", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenFinishingAgentCreationBeforeEach(t)
			request := generateValidFinishAgentCreationRequest()

			suite.agentSynchronizerMock.EXPECT().Sync(gomock.Any(), gomock.Any()).Return(assert.AnError)

			err := suite.sut.Finish(request)
			assert.Error(t, err)
		})
	})

	t.Run("Given synchronizer does not return error", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return err", func(t *testing.T) {
			suite := WhenFinishingAgentCreationBeforeEach(t)
			request := generateValidFinishAgentCreationRequest()

			suite.agentSynchronizerMock.EXPECT().Sync(gomock.Any(), gomock.Any()).Return(nil)

			err := suite.sut.Finish(request)
			assert.NoError(t, err)
		})
	})
}
