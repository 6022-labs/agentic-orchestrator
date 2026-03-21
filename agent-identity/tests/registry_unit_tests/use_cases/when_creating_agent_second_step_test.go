package use_cases_test

import (
	"testing"

	registry_errors "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenCreatingAgentSecondStepTestingSuite struct {
	sut *use_cases.CreateAgentSecondStep

	mockMintSigner       *interfaces_mocks.MockMintSignerInterface
	mockAgentRepository  *repositories_mocks.MockAgentRepositoryInterface
	orchestratorSettings *settings.OrchestratorSettings
}

func WhenCreatingAgentSecondStepBeforeEach(t *testing.T) *WhenCreatingAgentSecondStepTestingSuite {
	mockController := gomock.NewController(t)

	mockMintSigner := interfaces_mocks.NewMockMintSignerInterface(mockController)
	mockAgentRepository := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	orchestratorSettings := &settings.OrchestratorSettings{
		PublicKey: common_tests.RandomEthAddress(),
	}

	sut := use_cases.NewCreateAgentSecondStep(mockMintSigner, orchestratorSettings, mockAgentRepository)

	return &WhenCreatingAgentSecondStepTestingSuite{
		sut:                  sut,
		mockMintSigner:       mockMintSigner,
		mockAgentRepository:  mockAgentRepository,
		orchestratorSettings: orchestratorSettings,
	}
}

func generateValidAgentByIdRequest() *requests.AgentByIdRequest {
	return &requests.AgentByIdRequest{
		AgentId:             uuid.New(),
		AuthenticatedWallet: common_tests.RandomEthAddress(),
	}
}

func TestWhenCreatingAgentSecondStep(t *testing.T) {
	t.Parallel()

	t.Run("Given request validation fails", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error when agent ID is nil", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentSecondStepBeforeEach(t)

			req := &requests.AgentByIdRequest{
				AgentId:             uuid.Nil,
				AuthenticatedWallet: common_tests.RandomEthAddress(),
			}
			resp, err := basics.sut.Execute(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.IsType(t, &registry_errors.ValidationError{}, err)
		})

		t.Run("Should return a validation error when authenticated wallet is empty", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentSecondStepBeforeEach(t)

			req := &requests.AgentByIdRequest{
				AgentId:             uuid.New(),
				AuthenticatedWallet: common.Address{},
			}
			resp, err := basics.sut.Execute(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.IsType(t, &registry_errors.ValidationError{}, err)
		})
	})

	t.Run("Given error happens while fetching agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentSecondStepTestingSuite) {
			suite.mockAgentRepository.EXPECT().
				GetById(gomock.Any()).
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentSecondStepBeforeEach(t)
			initSuite(basics)

			req := generateValidAgentByIdRequest()
			resp, err := basics.sut.Execute(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given agent is not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentSecondStepTestingSuite) {
			suite.mockAgentRepository.EXPECT().
				GetById(gomock.Any()).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentSecondStepBeforeEach(t)
			initSuite(basics)

			req := generateValidAgentByIdRequest()
			resp, err := basics.sut.Execute(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.IsType(t, &registry_errors.NotFoundError{}, err)
		})
	})

	t.Run("Given agent is already minted", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentSecondStepBeforeEach(t)

			authenticatedWallet := common_tests.RandomEthAddress()
			agent := model_builders.NewAgentBuilder().
				WithCreator(authenticatedWallet.Hex()).
				WithCollectionTokenId(ptrUint64(42)).
				Build()

			basics.mockAgentRepository.EXPECT().
				GetById(agent.ID).
				Return(agent, nil).
				Times(1)

			req := &requests.AgentByIdRequest{
				AgentId:             agent.ID,
				AuthenticatedWallet: authenticatedWallet,
			}
			resp, err := basics.sut.Execute(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.IsType(t, &registry_errors.ValidationError{}, err)
		})
	})

	t.Run("Given authenticated wallet is not the creator", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentSecondStepBeforeEach(t)

			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(nil).
				Build()

			basics.mockAgentRepository.EXPECT().
				GetById(agent.ID).
				Return(agent, nil).
				Times(1)

			differentWallet := common_tests.RandomEthAddress()
			req := &requests.AgentByIdRequest{
				AgentId:             agent.ID,
				AuthenticatedWallet: differentWallet,
			}
			resp, err := basics.sut.Execute(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.IsType(t, &registry_errors.ValidationError{}, err)
		})
	})

	t.Run("Given mint signing fails", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentSecondStepBeforeEach(t)

			authenticatedWallet := common_tests.RandomEthAddress()
			agent := model_builders.NewAgentBuilder().
				WithCreator(authenticatedWallet.Hex()).
				WithCollectionTokenId(nil).
				Build()

			basics.mockAgentRepository.EXPECT().
				GetById(agent.ID).
				Return(agent, nil).
				Times(1)
			basics.mockMintSigner.EXPECT().
				Sign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError).
				Times(1)

			req := &requests.AgentByIdRequest{
				AgentId:             agent.ID,
				AuthenticatedWallet: authenticatedWallet,
			}
			resp, err := basics.sut.Execute(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
		})
	})

	t.Run("Given all conditions are met", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return signature, nonce and signer", func(t *testing.T) {
			t.Parallel()
			basics := WhenCreatingAgentSecondStepBeforeEach(t)

			authenticatedWallet := common_tests.RandomEthAddress()
			agent := model_builders.NewAgentBuilder().
				WithCreator(authenticatedWallet.Hex()).
				WithCollectionTokenId(nil).
				Build()

			expectedSignature := "0xsignature"

			basics.mockAgentRepository.EXPECT().
				GetById(agent.ID).
				Return(agent, nil).
				Times(1)
			basics.mockMintSigner.EXPECT().
				Sign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(&expectedSignature, nil).
				Times(1)

			req := &requests.AgentByIdRequest{
				AgentId:             agent.ID,
				AuthenticatedWallet: authenticatedWallet,
			}
			resp, err := basics.sut.Execute(req)

			assert.Nil(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, expectedSignature, resp.Signature)
			assert.NotEmpty(t, resp.Nonce)
			assert.Equal(t, basics.orchestratorSettings.PublicKey.Hex(), resp.Signer)
		})
	})
}

func ptrUint64(v uint64) *uint64 {
	return &v
}
