package use_cases_test

import (
	"context"

	"testing"
	"time"

	cb_errors "github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenVerifyingAuthTestingSuite struct {
	sut *use_cases.VerifyAuth

	mockNonceStore         *interfaces_mocks.MockNonceStoreInterface
	mockSiweMessageBuilder *interfaces_mocks.MockSiweMessageBuilderInterface
	mockSignatureVerifier  *interfaces_mocks.MockSignatureVerifierInterface
	mockJwtManager         *interfaces_mocks.MockJwtManagerInterface
}

func WhenVerifyingAuthBeforeEach(t *testing.T) *WhenVerifyingAuthTestingSuite {
	mockController := gomock.NewController(t)

	mockNonceStore := interfaces_mocks.NewMockNonceStoreInterface(mockController)
	mockSiweMessageBuilder := interfaces_mocks.NewMockSiweMessageBuilderInterface(mockController)
	mockSignatureVerifier := interfaces_mocks.NewMockSignatureVerifierInterface(mockController)
	mockJwtManager := interfaces_mocks.NewMockJwtManagerInterface(mockController)

	sut := use_cases.NewVerifyAuth(mockNonceStore, mockSiweMessageBuilder, mockSignatureVerifier, mockJwtManager)

	return &WhenVerifyingAuthTestingSuite{
		sut: sut,

		mockNonceStore:         mockNonceStore,
		mockSiweMessageBuilder: mockSiweMessageBuilder,
		mockSignatureVerifier:  mockSignatureVerifier,
		mockJwtManager:         mockJwtManager,
	}
}

func generateValidVerifyAuthRequest() *requests.VerifyAuthRequest {
	return &requests.VerifyAuthRequest{
		Address:   "0x4bca9E17a8326d283d1fC68A725B46b42c2D9146",
		Signature: "0x" + string(make132Hex()),
	}
}

func make132Hex() []byte {
	out := make([]byte, 130)
	for i := range out {
		out[i] = 'a'
	}
	return out
}

func TestWhenVerifyingAuth(t *testing.T) {
	t.Parallel()

	wallet := common.HexToAddress("0x4bca9E17a8326d283d1fC68A725B46b42c2D9146")
	entry := &interfaces.ChallengeEntry{
		Nonce:     "abc123",
		ChainId:   80002,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	siweMessage := "localhost:5173 wants you to sign in with your Ethereum account: ..."
	expectSiweMessage := func(suite *WhenVerifyingAuthTestingSuite) {
		suite.mockSiweMessageBuilder.EXPECT().
			Build(wallet, uint64(80002), "abc123", entry.IssuedAt, entry.ExpiresAt).
			Return(siweMessage)
	}

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because address is malformed", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAuthBeforeEach(t)
			request := generateValidVerifyAuthRequest()
			request.Address = "not-an-address"

			_, err := suite.sut.Execute(context.Background(), request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given the challenge is unknown or already used", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenVerifyingAuthTestingSuite) {
			suite.mockNonceStore.EXPECT().Consume(gomock.Any(), wallet).Return(nil).Times(1)
		}

		t.Run("Should return an unauthorized error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAuthBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(context.Background(), generateValidVerifyAuthRequest())

			var unauthorizedErr *cb_errors.UnauthorizedError
			assert.ErrorAs(t, err, &unauthorizedErr)
			assert.Equal(t, "challenge_not_found", unauthorizedErr.Code)
		})
	})

	t.Run("Given the signature does not match the address", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenVerifyingAuthTestingSuite, request *requests.VerifyAuthRequest) {
			suite.mockNonceStore.EXPECT().Consume(gomock.Any(), wallet).Return(entry).Times(1)
			expectSiweMessage(suite)
			suite.mockSignatureVerifier.EXPECT().
				Recover(siweMessage, request.Signature).
				Return(common.HexToAddress("0x1111111111111111111111111111111111111111"), nil).
				Times(1)
		}

		t.Run("Should return an unauthorized error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAuthBeforeEach(t)
			request := generateValidVerifyAuthRequest()
			initSuite(suite, request)

			_, err := suite.sut.Execute(context.Background(), request)

			var unauthorizedErr *cb_errors.UnauthorizedError
			assert.ErrorAs(t, err, &unauthorizedErr)
			assert.Equal(t, "signature_mismatch", unauthorizedErr.Code)
		})
	})

	t.Run("Given the signed nonce is valid", func(t *testing.T) {
		t.Parallel()

		expiresAt := time.Now().Add(time.Hour)

		initSuite := func(suite *WhenVerifyingAuthTestingSuite, request *requests.VerifyAuthRequest) {
			suite.mockNonceStore.EXPECT().Consume(gomock.Any(), wallet).Return(entry).Times(1)
			expectSiweMessage(suite)
			suite.mockSignatureVerifier.EXPECT().
				Recover(siweMessage, request.Signature).
				Return(wallet, nil).
				Times(1)
			suite.mockJwtManager.EXPECT().
				CreateJwtToken("0x4bca9e17a8326d283d1fc68a725b46b42c2d9146").
				Return("jwt-token", expiresAt, nil).
				Times(1)
		}

		t.Run("Should return the token and expiry", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAuthBeforeEach(t)
			request := generateValidVerifyAuthRequest()
			initSuite(suite, request)

			response, err := suite.sut.Execute(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, "jwt-token", response.Token)
			assert.Equal(t, expiresAt, response.ExpiresAt)
		})
	})
}
