package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenExportingAgentTestSuite struct {
	agentRegistryRequesterMock      *interfaces_mocks.MockAgentRegistryRequesterInterface
	agentPrivateKeyStorageMock      *interfaces_mocks.MockAgentPrivateKeyStorageInterface
	agentCollectionRequesterMock    *interfaces_mocks.MockAgentCollectionRequesterInterface
	llmConnectionRequesterMock      *interfaces_mocks.MockLLMConnectionRequesterInterface
	databaseConnectionRequesterMock *interfaces_mocks.MockDatabaseConnectionRequesterInterface

	sut *use_cases.ExportAgent
}

func WhenExportingAgentTestSuiteBeforeEach(t *testing.T) *WhenExportingAgentTestSuite {
	// Create mocks
	mockController := gomock.NewController(t, gomock.WithOverridableExpectations())

	agentRegistryRequesterMock := interfaces_mocks.NewMockAgentRegistryRequesterInterface(mockController)
	agentPrivateKeyStorageMock := interfaces_mocks.NewMockAgentPrivateKeyStorageInterface(mockController)
	agentCollectionRequesterMock := interfaces_mocks.NewMockAgentCollectionRequesterInterface(mockController)
	llmConnectionRequesterMock := interfaces_mocks.NewMockLLMConnectionRequesterInterface(mockController)
	databaseConnectionRequesterMock := interfaces_mocks.NewMockDatabaseConnectionRequesterInterface(mockController)

	sut := use_cases.NewExportAgent(
		agentRegistryRequesterMock,
		agentPrivateKeyStorageMock,
		agentCollectionRequesterMock,
		llmConnectionRequesterMock,
		databaseConnectionRequesterMock,
	)

	return &WhenExportingAgentTestSuite{
		sut: sut,

		agentRegistryRequesterMock:      agentRegistryRequesterMock,
		agentPrivateKeyStorageMock:      agentPrivateKeyStorageMock,
		agentCollectionRequesterMock:    agentCollectionRequesterMock,
		llmConnectionRequesterMock:      llmConnectionRequesterMock,
		databaseConnectionRequesterMock: databaseConnectionRequesterMock,
	}
}

func TestWhenExportingAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given error happens while fetching agent from registry", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenExportingAgentTestSuite) *requests.AgentByCollectionTokenIdRequest {
			wallet := common_tests.RandomEthAddress()

			suite.agentRegistryRequesterMock.
				EXPECT().
				GetAgent(gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError)

			return &requests.AgentByCollectionTokenIdRequest{
				AuthenticatedWallet:    &wallet,
				AgentCollectionAddress: common_tests.RandomEthAddress().Hex(),
				AgentCollectionTokenId: 1,
			}
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			basics := WhenExportingAgentTestSuiteBeforeEach(t)
			req := initSuite(basics)

			_, err := basics.sut.Export(req)

			assert.NotNil(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given agent is not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenExportingAgentTestSuite) *requests.AgentByCollectionTokenIdRequest {
			wallet := common_tests.RandomEthAddress()
			collectionAddress := common_tests.RandomEthAddress().Hex()
			var collectionTokenId uint64 = 42

			suite.agentRegistryRequesterMock.
				EXPECT().
				GetAgent(collectionAddress, collectionTokenId).
				Return(nil, nil)

			signerAddress := common_tests.RandomEthAddress()
			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentSigner(gomock.Any(), gomock.Any()).
				Return(&signerAddress, nil)

			return &requests.AgentByCollectionTokenIdRequest{
				AuthenticatedWallet:    &wallet,
				AgentCollectionAddress: collectionAddress,
				AgentCollectionTokenId: collectionTokenId,
			}
		}

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			basics := WhenExportingAgentTestSuiteBeforeEach(t)
			req := initSuite(basics)

			_, err := basics.sut.Export(req)

			assert.NotNil(t, err)
			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "agentCollectionTokenId", validationErr.Field)
			assert.Equal(t, "agent not found", validationErr.Message)
		})
	})

	t.Run("Given agent is not internal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenExportingAgentTestSuite) *requests.AgentByCollectionTokenIdRequest {
			wallet := common_tests.RandomEthAddress()
			agent := model_builders.NewAgentBuilder().WithIsInternal(false).Build()
			collectionAddress := agent.CollectionAddress
			collectionTokenId := *agent.CollectionTokenId

			suite.agentRegistryRequesterMock.
				EXPECT().
				GetAgent(collectionAddress, collectionTokenId).
				Return(agent, nil)

			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentSigner(gomock.Any(), gomock.Any()).
				Return(&wallet, nil)

			return &requests.AgentByCollectionTokenIdRequest{
				AuthenticatedWallet:    &wallet,
				AgentCollectionAddress: collectionAddress,
				AgentCollectionTokenId: collectionTokenId,
			}
		}

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			basics := WhenExportingAgentTestSuiteBeforeEach(t)
			req := initSuite(basics)

			_, err := basics.sut.Export(req)

			assert.NotNil(t, err)
			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "agentCollectionTokenId", validationErr.Field)
			assert.Equal(t, "Agent is not internal", validationErr.Message)
		})
	})

	t.Run("Given authenticated wallet is not the agent signer", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenExportingAgentTestSuite) *requests.AgentByCollectionTokenIdRequest {
			wallet := common_tests.RandomEthAddress()
			signerAddress := common_tests.RandomEthAddress()
			agent := model_builders.NewAgentBuilder().WithIsInternal(true).Build()
			collectionAddress := agent.CollectionAddress
			collectionTokenId := *agent.CollectionTokenId

			suite.agentRegistryRequesterMock.
				EXPECT().
				GetAgent(collectionAddress, collectionTokenId).
				Return(agent, nil)

			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentSigner(gomock.Any(), gomock.Any()).
				Return(&signerAddress, nil)

			return &requests.AgentByCollectionTokenIdRequest{
				AuthenticatedWallet:    &wallet,
				AgentCollectionAddress: collectionAddress,
				AgentCollectionTokenId: collectionTokenId,
			}
		}

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			basics := WhenExportingAgentTestSuiteBeforeEach(t)
			req := initSuite(basics)

			_, err := basics.sut.Export(req)

			assert.NotNil(t, err)
			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "agentCollectionTokenId", validationErr.Field)
			assert.Equal(t, "You are not the owner of this agent", validationErr.Message)
		})
	})

	t.Run("Given agent export is successful", func(t *testing.T) {
		t.Parallel()

		var expectedPrivateKey string

		initSuite := func(suite *WhenExportingAgentTestSuite) *requests.AgentByCollectionTokenIdRequest {
			wallet := common_tests.RandomEthAddress()
			agent := model_builders.NewAgentBuilder().WithIsInternal(true).Build()
			collectionAddress := agent.CollectionAddress
			collectionTokenId := *agent.CollectionTokenId

			suite.agentRegistryRequesterMock.
				EXPECT().
				GetAgent(collectionAddress, collectionTokenId).
				Return(agent, nil)

			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentSigner(gomock.Any(), gomock.Any()).
				Return(&wallet, nil)

			expectedPrivateKey = "0xdeadbeef1234567890"
			suite.agentPrivateKeyStorageMock.
				EXPECT().
				Get(agent.ID).
				Return(&expectedPrivateKey, nil)

			suite.agentPrivateKeyStorageMock.
				EXPECT().
				Delete(agent.ID).
				Return(nil)

			suite.databaseConnectionRequesterMock.
				EXPECT().
				Remove(agent.ID).
				Return(nil)

			suite.llmConnectionRequesterMock.
				EXPECT().
				Remove(agent.ID).
				Return(nil)

			return &requests.AgentByCollectionTokenIdRequest{
				AuthenticatedWallet:    &wallet,
				AgentCollectionAddress: collectionAddress,
				AgentCollectionTokenId: collectionTokenId,
			}
		}

		t.Run("Should return the private key", func(t *testing.T) {
			t.Parallel()

			basics := WhenExportingAgentTestSuiteBeforeEach(t)
			req := initSuite(basics)

			response, err := basics.sut.Export(req)

			assert.Nil(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, expectedPrivateKey, response.PrivateKey)
		})
	})

	t.Run("Given private key is not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenExportingAgentTestSuite) *requests.AgentByCollectionTokenIdRequest {
			wallet := common_tests.RandomEthAddress()
			agent := model_builders.NewAgentBuilder().WithIsInternal(true).Build()
			collectionAddress := agent.CollectionAddress
			collectionTokenId := *agent.CollectionTokenId

			suite.agentRegistryRequesterMock.
				EXPECT().
				GetAgent(collectionAddress, collectionTokenId).
				Return(agent, nil)

			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentSigner(gomock.Any(), gomock.Any()).
				Return(&wallet, nil)

			suite.agentPrivateKeyStorageMock.
				EXPECT().
				Get(agent.ID).
				Return(nil, nil)

			return &requests.AgentByCollectionTokenIdRequest{
				AuthenticatedWallet:    &wallet,
				AgentCollectionAddress: collectionAddress,
				AgentCollectionTokenId: collectionTokenId,
			}
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			basics := WhenExportingAgentTestSuiteBeforeEach(t)
			req := initSuite(basics)

			_, err := basics.sut.Export(req)

			assert.NotNil(t, err)
			assert.Equal(t, "private key not found", err.Error())
		})
	})

	t.Run("Given error happens while getting agent signer", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenExportingAgentTestSuite) *requests.AgentByCollectionTokenIdRequest {
			wallet := common_tests.RandomEthAddress()
			collectionAddress := common_tests.RandomEthAddress().Hex()
			var collectionTokenId uint64 = 1

			agent := model_builders.NewAgentBuilder().WithIsInternal(true).Build()

			suite.agentRegistryRequesterMock.
				EXPECT().
				GetAgent(collectionAddress, collectionTokenId).
				Return(agent, nil)

			suite.agentCollectionRequesterMock.
				EXPECT().
				GetAgentSigner(gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError)

			return &requests.AgentByCollectionTokenIdRequest{
				AuthenticatedWallet:    &wallet,
				AgentCollectionAddress: collectionAddress,
				AgentCollectionTokenId: collectionTokenId,
			}
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			basics := WhenExportingAgentTestSuiteBeforeEach(t)
			req := initSuite(basics)

			_, err := basics.sut.Export(req)

			assert.NotNil(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})
}

func TestWhenExportingAgentWithNilAuthenticatedWallet(t *testing.T) {
	t.Parallel()

	t.Run("Given request has nil authenticated wallet", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error from EnsureAuthenticatedWallet", func(t *testing.T) {
			t.Parallel()

			req := &requests.AgentByCollectionTokenIdRequest{
				AuthenticatedWallet:    nil,
				AgentCollectionAddress: common_tests.RandomEthAddress().Hex(),
				AgentCollectionTokenId: 1,
			}

			err := req.EnsureAuthenticatedWallet()

			assert.NotNil(t, err)
			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "auth", validationErr.Field)
		})
	})

	t.Run("Given request has zero authenticated wallet", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error from EnsureAuthenticatedWallet", func(t *testing.T) {
			t.Parallel()

			zeroAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")
			req := &requests.AgentByCollectionTokenIdRequest{
				AuthenticatedWallet:    &zeroAddress,
				AgentCollectionAddress: common_tests.RandomEthAddress().Hex(),
				AgentCollectionTokenId: 1,
			}

			err := req.EnsureAuthenticatedWallet()

			assert.NotNil(t, err)
			var validationErr *errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "auth", validationErr.Field)
		})
	})
}
