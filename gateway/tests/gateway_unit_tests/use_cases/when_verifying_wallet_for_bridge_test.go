package use_cases_test

import (
	"context"
	"strings"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
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
	return &requests.VerifyWalletForBridgeRequest{
		Platform:  enums.BridgePlatformSlack,
		Challenge: uuid.NewString(),
		Signature: uuid.NewString(),
	}
}

func TestWhenVerifyingWalletForBridge(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid because platform is not part of the enum", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingWalletForBridgeBeforeEach(t)
			suite.gatewayMetricsMock.EXPECT().RecordWalletVerification(gomock.Any(), gomock.Not(gomock.Nil()))
			request := generateValidVerifyWalletForBridgeRequest()
			request.Platform = enums.BridgePlatform("none")

			err := suite.sut.Verify(context.Background(), request)

			var validationErr *cb_errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "platform", validationErr.Field)
		})
	})

	t.Run("Given no caller id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an unauthorized error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingWalletForBridgeBeforeEach(t)
			suite.gatewayMetricsMock.EXPECT().RecordWalletVerification(gomock.Any(), gomock.Not(gomock.Nil()))

			err := suite.sut.Verify(context.Background(), generateValidVerifyWalletForBridgeRequest())

			var unauthorizedErr *cb_errors.UnauthorizedError
			assert.ErrorAs(t, err, &unauthorizedErr)
			assert.Equal(t, "missing_caller_id", unauthorizedErr.Code)
		})
	})

	t.Run("Given the signature cannot be recovered", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenVerifyingWalletForBridgeTestingSuite, request *requests.VerifyWalletForBridgeRequest) {
			suite.gatewayMetricsMock.EXPECT().RecordWalletVerification(gomock.Any(), gomock.Not(gomock.Nil()))
			suite.mockSignatureVerifier.
				EXPECT().
				Recover(request.Challenge, request.Signature).
				Return(common.Address{}, assert.AnError)
		}

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingWalletForBridgeBeforeEach(t)
			ctx := helpers.ContextWithCallerId(
				context.Background(),
				strings.ToLower(common_tests.RandomEthAddress().Hex()),
			)
			request := generateValidVerifyWalletForBridgeRequest()
			initSuite(suite, request)

			err := suite.sut.Verify(ctx, request)

			var validationErr *cb_errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "signature", validationErr.Field)
		})
	})

	t.Run("Given the signature does not match the caller", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenVerifyingWalletForBridgeTestingSuite, request *requests.VerifyWalletForBridgeRequest) {
			suite.gatewayMetricsMock.EXPECT().RecordWalletVerification(gomock.Any(), gomock.Not(gomock.Nil()))
			suite.mockSignatureVerifier.
				EXPECT().
				Recover(request.Challenge, request.Signature).
				Return(common.HexToAddress("0x1111111111111111111111111111111111111111"), nil)
		}

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingWalletForBridgeBeforeEach(t)
			ctx := helpers.ContextWithCallerId(
				context.Background(),
				strings.ToLower(common_tests.RandomEthAddress().Hex()),
			)
			request := generateValidVerifyWalletForBridgeRequest()
			initSuite(suite, request)

			err := suite.sut.Verify(ctx, request)

			var validationErr *cb_errors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "signature", validationErr.Field)
		})
	})

	t.Run("Given the platform bridge fails to set the wallet ownership", func(t *testing.T) {
		t.Parallel()

		initSuite := func(
			suite *WhenVerifyingWalletForBridgeTestingSuite,
			request *requests.VerifyWalletForBridgeRequest,
			callerId string,
		) {
			suite.gatewayMetricsMock.EXPECT().RecordWalletVerification(gomock.Any(), gomock.Not(gomock.Nil()))
			suite.mockSignatureVerifier.
				EXPECT().
				Recover(request.Challenge, request.Signature).
				Return(common.HexToAddress(callerId), nil)
			suite.mockPlatformBridgeRequester.
				EXPECT().
				SetWalletOwnershipToUser(gomock.Any(), request.Platform, request.Challenge, callerId).
				Return(assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingWalletForBridgeBeforeEach(t)
			callerId := strings.ToLower(common_tests.RandomEthAddress().Hex())
			ctx := helpers.ContextWithCallerId(context.Background(), callerId)
			request := generateValidVerifyWalletForBridgeRequest()
			initSuite(suite, request, callerId)

			err := suite.sut.Verify(ctx, request)

			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Given the signature matches the caller", func(t *testing.T) {
		t.Parallel()

		initSuite := func(
			suite *WhenVerifyingWalletForBridgeTestingSuite,
			request *requests.VerifyWalletForBridgeRequest,
			callerId string,
		) {
			suite.gatewayMetricsMock.EXPECT().RecordWalletVerification(gomock.Any(), gomock.Nil())
			suite.mockSignatureVerifier.
				EXPECT().
				Recover(request.Challenge, request.Signature).
				Return(common.HexToAddress(callerId), nil)
			suite.mockPlatformBridgeRequester.
				EXPECT().
				SetWalletOwnershipToUser(gomock.Any(), request.Platform, request.Challenge, callerId).
				Return(nil)
		}

		t.Run("Should set the caller as the wallet owner", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingWalletForBridgeBeforeEach(t)
			callerId := strings.ToLower(common_tests.RandomEthAddress().Hex())
			ctx := helpers.ContextWithCallerId(context.Background(), callerId)
			request := generateValidVerifyWalletForBridgeRequest()
			initSuite(suite, request, callerId)

			err := suite.sut.Verify(ctx, request)

			assert.NoError(t, err)
		})
	})
}
