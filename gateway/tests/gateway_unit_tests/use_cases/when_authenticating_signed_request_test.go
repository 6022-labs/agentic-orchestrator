package use_cases_test

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	cb_errors "github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenAuthenticatingSignedRequestTestingSuite struct {
	sut *use_cases.AuthenticateSignedRequest

	mockSignedRequestMessageBuilder *interfaces_mocks.MockSignedRequestMessageBuilderInterface
	mockSignatureVerifier           *interfaces_mocks.MockSignatureVerifierInterface
	mockCacheStorage                *interfaces_mocks.MockCacheStorageInterface
}

func WhenAuthenticatingSignedRequestBeforeEach(t *testing.T) *WhenAuthenticatingSignedRequestTestingSuite {
	mockController := gomock.NewController(t)

	mockSignedRequestMessageBuilder := interfaces_mocks.NewMockSignedRequestMessageBuilderInterface(mockController)
	mockSignatureVerifier := interfaces_mocks.NewMockSignatureVerifierInterface(mockController)
	mockCacheStorage := interfaces_mocks.NewMockCacheStorageInterface(mockController)

	sut := use_cases.NewAuthenticateSignedRequest(
		&settings.SiweSettings{SignedRequestMaxAgeSeconds: 60},
		mockSignedRequestMessageBuilder,
		mockSignatureVerifier,
		mockCacheStorage,
	)

	return &WhenAuthenticatingSignedRequestTestingSuite{
		sut: sut,

		mockSignedRequestMessageBuilder: mockSignedRequestMessageBuilder,
		mockSignatureVerifier:           mockSignatureVerifier,
		mockCacheStorage:                mockCacheStorage,
	}
}

func generateValidAuthenticateSignedRequestRequest(signedAt time.Time) *requests.AuthenticateSignedRequestRequest {
	return &requests.AuthenticateSignedRequestRequest{
		Method:    "GET",
		Path:      "/agents",
		Query:     "page=2",
		Timestamp: strconv.FormatInt(signedAt.Unix(), 10),
		Nonce:     "0123456789abcdef",
		Signature: "0x" + strings.Repeat("ab", 65),
	}
}

func TestWhenAuthenticatingSignedRequest(t *testing.T) {
	t.Parallel()

	signedMessage := "6022 orchestrator signed request ..."

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()
			suite := WhenAuthenticatingSignedRequestBeforeEach(t)
			request := generateValidAuthenticateSignedRequestRequest(time.Now())
			request.Nonce = "short"

			_, err := suite.sut.Execute(context.Background(), request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given a timestamp older than the allowed window", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an unauthorized error", func(t *testing.T) {
			t.Parallel()
			suite := WhenAuthenticatingSignedRequestBeforeEach(t)
			request := generateValidAuthenticateSignedRequestRequest(time.Now().Add(-2 * time.Minute))

			_, err := suite.sut.Execute(context.Background(), request)

			var unauthorizedErr *cb_errors.UnauthorizedError
			assert.ErrorAs(t, err, &unauthorizedErr)
			assert.Equal(t, "signed_request_expired", unauthorizedErr.Code)
		})
	})

	t.Run("Given a timestamp further in the future than the allowed window", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an unauthorized error", func(t *testing.T) {
			t.Parallel()
			suite := WhenAuthenticatingSignedRequestBeforeEach(t)
			request := generateValidAuthenticateSignedRequestRequest(time.Now().Add(2 * time.Minute))

			_, err := suite.sut.Execute(context.Background(), request)

			var unauthorizedErr *cb_errors.UnauthorizedError
			assert.ErrorAs(t, err, &unauthorizedErr)
			assert.Equal(t, "signed_request_expired", unauthorizedErr.Code)
		})
	})

	t.Run("Given the signature cannot be recovered", func(t *testing.T) {
		t.Parallel()

		request := generateValidAuthenticateSignedRequestRequest(time.Now())
		initSuite := func(suite *WhenAuthenticatingSignedRequestTestingSuite) {
			suite.mockSignedRequestMessageBuilder.EXPECT().
				Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(signedMessage).
				Times(1)
			suite.mockSignatureVerifier.EXPECT().
				Recover(signedMessage, request.Signature).
				Return(common_tests.RandomEthAddress(), assert.AnError).
				Times(1)
		}

		t.Run("Should return an unauthorized error", func(t *testing.T) {
			t.Parallel()
			suite := WhenAuthenticatingSignedRequestBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(context.Background(), request)

			var unauthorizedErr *cb_errors.UnauthorizedError
			assert.ErrorAs(t, err, &unauthorizedErr)
			assert.Equal(t, "signature_invalid", unauthorizedErr.Code)
		})
	})

	t.Run("Given the nonce store fails", func(t *testing.T) {
		t.Parallel()

		request := generateValidAuthenticateSignedRequestRequest(time.Now())
		wallet := common_tests.RandomEthAddress()
		initSuite := func(suite *WhenAuthenticatingSignedRequestTestingSuite) {
			suite.mockSignedRequestMessageBuilder.EXPECT().
				Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(signedMessage).
				Times(1)
			suite.mockSignatureVerifier.EXPECT().Recover(signedMessage, request.Signature).Return(wallet, nil).Times(1)
			suite.mockCacheStorage.EXPECT().
				StoreIfAbsent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(false, assert.AnError).
				Times(1)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()
			suite := WhenAuthenticatingSignedRequestBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(context.Background(), request)

			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Given the nonce was already used", func(t *testing.T) {
		t.Parallel()

		request := generateValidAuthenticateSignedRequestRequest(time.Now())
		wallet := common_tests.RandomEthAddress()
		initSuite := func(suite *WhenAuthenticatingSignedRequestTestingSuite) {
			suite.mockSignedRequestMessageBuilder.EXPECT().
				Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(signedMessage).
				Times(1)
			suite.mockSignatureVerifier.EXPECT().Recover(signedMessage, request.Signature).Return(wallet, nil).Times(1)
			suite.mockCacheStorage.EXPECT().
				StoreIfAbsent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(false, nil).
				Times(1)
		}

		t.Run("Should return an unauthorized error", func(t *testing.T) {
			t.Parallel()
			suite := WhenAuthenticatingSignedRequestBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(context.Background(), request)

			var unauthorizedErr *cb_errors.UnauthorizedError
			assert.ErrorAs(t, err, &unauthorizedErr)
			assert.Equal(t, "nonce_already_used", unauthorizedErr.Code)
		})
	})

	t.Run("Given a fresh signed request", func(t *testing.T) {
		t.Parallel()

		signedAt := time.Now()
		request := generateValidAuthenticateSignedRequestRequest(signedAt)
		wallet := common_tests.RandomEthAddress()
		walletAddress := strings.ToLower(wallet.Hex())
		initSuite := func(suite *WhenAuthenticatingSignedRequestTestingSuite) {
			suite.mockSignatureVerifier.EXPECT().Recover(signedMessage, request.Signature).Return(wallet, nil).Times(1)
		}
		expectAnyBuild := func(suite *WhenAuthenticatingSignedRequestTestingSuite) {
			suite.mockSignedRequestMessageBuilder.EXPECT().
				Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(signedMessage).
				Times(1)
		}
		expectStored := func(suite *WhenAuthenticatingSignedRequestTestingSuite) {
			suite.mockCacheStorage.EXPECT().
				StoreIfAbsent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(true, nil).
				Times(1)
		}

		t.Run("Should build the message from the request", func(t *testing.T) {
			t.Parallel()
			suite := WhenAuthenticatingSignedRequestBeforeEach(t)
			initSuite(suite)
			expectStored(suite)
			suite.mockSignedRequestMessageBuilder.EXPECT().
				Build("GET", "/agents", "page=2", signedAt.Unix(), "0123456789abcdef", gomock.Nil()).
				Return(signedMessage).
				Times(1)

			_, err := suite.sut.Execute(context.Background(), request)

			assert.Nil(t, err)
		})

		t.Run("Should store the nonce for twice the window", func(t *testing.T) {
			t.Parallel()
			suite := WhenAuthenticatingSignedRequestBeforeEach(t)
			initSuite(suite)
			expectAnyBuild(suite)
			suite.mockCacheStorage.EXPECT().
				StoreIfAbsent(
					gomock.Any(),
					"auth/signed-requests/"+walletAddress+"/0123456789abcdef",
					gomock.Any(),
					2*time.Minute,
				).
				Return(true, nil).
				Times(1)

			_, err := suite.sut.Execute(context.Background(), request)

			assert.Nil(t, err)
		})

		t.Run("Should return the lowercased signer wallet", func(t *testing.T) {
			t.Parallel()
			suite := WhenAuthenticatingSignedRequestBeforeEach(t)
			initSuite(suite)
			expectAnyBuild(suite)
			expectStored(suite)

			response, err := suite.sut.Execute(context.Background(), request)

			assert.Nil(t, err)
			if assert.NotNil(t, response) {
				assert.Equal(t, walletAddress, response.Wallet)
			}
		})
	})
}
