package use_cases_test

import (
	"context"
	"strings"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	metrics_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/metrics_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenVerifyingWalletForBridgeTestingSuite struct {
	gatewayMetricsMock *metrics_mocks.MockGatewayMetricsInterface
	sut                *use_cases.VerifyWalletForBridge

	mockSignatureVerifier       interfaces_mocks.MockSignatureVerifierInterface
	mockPlatformBridgeRequester interfaces_mocks.MockPlatformBridgeRequesterInterface
}

func WhenVerifyingWalletForBridgeBeforeEach(t *testing.T) *WhenVerifyingWalletForBridgeTestingSuite {
	mockController := gomock.NewController(t)

	mockSignatureVerifier := interfaces_mocks.NewMockSignatureVerifierInterface(mockController)
	mockPlatformBridgeRequester := interfaces_mocks.NewMockPlatformBridgeRequesterInterface(mockController)

	gatewayMetricsMock := metrics_mocks.NewMockGatewayMetricsInterface(mockController)

	sut := use_cases.NewVerifyWalletForBridge(
		gatewayMetricsMock,
		mockSignatureVerifier,
		mockPlatformBridgeRequester,
	)

	return &WhenVerifyingWalletForBridgeTestingSuite{
		sut:                         sut,
		gatewayMetricsMock:          gatewayMetricsMock,
		mockSignatureVerifier:       *mockSignatureVerifier,
		mockPlatformBridgeRequester: *mockPlatformBridgeRequester,
	}
}

func generateValidVerifyWalletForBridgeRequest() *requests.VerifyWalletForBridgeRequest {
	authenticatedWallet := common_tests.RandomEthAddress()
	challenge := uuid.NewString()
	signature := uuid.NewString()

	return &requests.VerifyWalletForBridgeRequest{
		AuthenticatedWallet: authenticatedWallet,
		Platform:            enums.BridgePlatformSlack,
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
			suite.gatewayMetricsMock.EXPECT().RecordWalletVerification(gomock.Any(), gomock.Not(gomock.Nil()))

			request := generateValidVerifyWalletForBridgeRequest()
			request.Platform = enums.BridgePlatform("none")

			err := suite.sut.Verify(context.Background(), request)
			assert.NotNil(t, err)
		})
	})

	t.Run("Given signature does not match the authenticated wallet, should return error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenVerifyingWalletForBridgeTestingSuite, request *requests.VerifyWalletForBridgeRequest) {
			suite.mockSignatureVerifier.
				EXPECT().
				Recover(request.Challenge, request.Signature).
				Return(common.HexToAddress("0x1111111111111111111111111111111111111111"), nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingWalletForBridgeBeforeEach(t)
			suite.gatewayMetricsMock.EXPECT().RecordWalletVerification(gomock.Any(), gomock.Not(gomock.Nil()))
			request := generateValidVerifyWalletForBridgeRequest()
			initSuite(suite, request)

			err := suite.sut.Verify(context.Background(), request)
			assert.NotNil(t, err)
		})
	})

	t.Run(
		"Given error occurs while setting user wallet ownership on platform bridge, should return error",
		func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenVerifyingWalletForBridgeTestingSuite, request *requests.VerifyWalletForBridgeRequest) {
				suite.mockSignatureVerifier.
					EXPECT().
					Recover(request.Challenge, request.Signature).
					Return(request.AuthenticatedWallet, nil)

				suite.mockPlatformBridgeRequester.
					EXPECT().
					SetWalletOwnershipToUser(
						gomock.Any(),
						request.Platform, request.Challenge, strings.ToLower(request.AuthenticatedWallet.Hex()),
					).
					Return(assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				suite := WhenVerifyingWalletForBridgeBeforeEach(t)
				suite.gatewayMetricsMock.EXPECT().RecordWalletVerification(gomock.Any(), gomock.Not(gomock.Nil()))
				request := generateValidVerifyWalletForBridgeRequest()
				initSuite(suite, request)

				err := suite.sut.Verify(context.Background(), request)
				assert.NotNil(t, err)
			})
		},
	)

	t.Run("Given wallet is successfully setup to user", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenVerifyingWalletForBridgeTestingSuite, request *requests.VerifyWalletForBridgeRequest) {
			suite.mockSignatureVerifier.
				EXPECT().
				Recover(request.Challenge, request.Signature).
				Return(request.AuthenticatedWallet, nil)

			suite.mockPlatformBridgeRequester.
				EXPECT().
				SetWalletOwnershipToUser(
					gomock.Any(),
					request.Platform, request.Challenge, strings.ToLower(request.AuthenticatedWallet.Hex()),
				).
				Return(nil)
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingWalletForBridgeBeforeEach(t)
			suite.gatewayMetricsMock.EXPECT().RecordWalletVerification(gomock.Any(), gomock.Nil())
			request := generateValidVerifyWalletForBridgeRequest()
			initSuite(suite, request)

			err := suite.sut.Verify(context.Background(), request)
			assert.Nil(t, err)
		})
	})
}
