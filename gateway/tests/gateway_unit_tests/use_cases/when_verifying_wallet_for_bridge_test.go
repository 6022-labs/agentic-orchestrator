package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenVerifyingWalletForBridgeTestingSuite struct {
	sut *use_cases.VerifyWalletForBridge

	mockSignedMessageValidator  interfaces_mocks.MockSignedMessageValidatorInterface
	mockPlatformBridgeRequester interfaces_mocks.MockPlatformBridgeRequesterInterface
}

func WhenVerifyingWalletForBridgeBeforeEach(t *testing.T) *WhenVerifyingWalletForBridgeTestingSuite {
	mockController := gomock.NewController(t)

	mockSignedMessageValidator := interfaces_mocks.NewMockSignedMessageValidatorInterface(mockController)
	mockPlatformBridgeRequester := interfaces_mocks.NewMockPlatformBridgeRequesterInterface(mockController)

	sut := use_cases.NewVerifyWalletForBridge(
		mockSignedMessageValidator,
		mockPlatformBridgeRequester,
	)

	return &WhenVerifyingWalletForBridgeTestingSuite{
		sut:                         sut,
		mockSignedMessageValidator:  *mockSignedMessageValidator,
		mockPlatformBridgeRequester: *mockPlatformBridgeRequester,
	}
}

func generateValidVerifyWalletForBridgeRequest() *requests.VerifyWalletForBridgeRequest {
	authenticatedWallet := common_tests.RandomEthAddress()
	challenge := uuid.NewString()
	signature := uuid.NewString()

	return &requests.VerifyWalletForBridgeRequest{
		AuthenticatedWallet: authenticatedWallet,
		Platform:            enums.AgentContextMessagePlatformSlack,
		Challenge:           challenge,
		Signature:           signature,
	}
}

func TestWhenVerifyingWalletForBridge(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because platform is not part of the enum, should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingWalletForBridgeBeforeEach(t)

			request := generateValidVerifyWalletForBridgeRequest()
			request.Platform = enums.AgentContextMessagePlatformNone

			err := suite.sut.Verify(request)
			assert.NotNil(t, err)
		})

		t.Run("Because challenge is empty, should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingWalletForBridgeBeforeEach(t)

			request := generateValidVerifyWalletForBridgeRequest()
			request.Challenge = "   "

			err := suite.sut.Verify(request)
			assert.NotNil(t, err)
		})

		t.Run("Because signature is empty, should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingWalletForBridgeBeforeEach(t)

			request := generateValidVerifyWalletForBridgeRequest()
			request.Signature = "   "

			err := suite.sut.Verify(request)
			assert.NotNil(t, err)
		})
	})

	t.Run("Given signature does not match the authenticated wallet, should return error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenVerifyingWalletForBridgeTestingSuite, request *requests.VerifyWalletForBridgeRequest) {
			suite.mockSignedMessageValidator.
				EXPECT().
				Validate(request.Challenge, request.Signature, request.AuthenticatedWallet.Hex()).
				Return(false)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingWalletForBridgeBeforeEach(t)
			request := generateValidVerifyWalletForBridgeRequest()
			initSuite(suite, request)

			err := suite.sut.Verify(request)
			assert.NotNil(t, err)
		})
	})

	t.Run("Given error occurs while setting user wallet ownership on platform bridge, should return error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenVerifyingWalletForBridgeTestingSuite, request *requests.VerifyWalletForBridgeRequest) {
			suite.mockSignedMessageValidator.
				EXPECT().
				Validate(request.Challenge, request.Signature, request.AuthenticatedWallet.Hex()).
				Return(true)

			suite.mockPlatformBridgeRequester.
				EXPECT().
				SetWalletOwnershipToUser(request.Platform, request.Challenge, request.AuthenticatedWallet.Hex()).
				Return(assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingWalletForBridgeBeforeEach(t)
			request := generateValidVerifyWalletForBridgeRequest()
			initSuite(suite, request)

			err := suite.sut.Verify(request)
			assert.NotNil(t, err)
		})
	})

	t.Run("Given wallet is successfully setup to user", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenVerifyingWalletForBridgeTestingSuite, request *requests.VerifyWalletForBridgeRequest) {
			suite.mockSignedMessageValidator.
				EXPECT().
				Validate(request.Challenge, request.Signature, request.AuthenticatedWallet.Hex()).
				Return(true)

			suite.mockPlatformBridgeRequester.
				EXPECT().
				SetWalletOwnershipToUser(request.Platform, request.Challenge, request.AuthenticatedWallet.Hex()).
				Return(nil)
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingWalletForBridgeBeforeEach(t)
			request := generateValidVerifyWalletForBridgeRequest()
			initSuite(suite, request)

			err := suite.sut.Verify(request)
			assert.Nil(t, err)
		})
	})
}
