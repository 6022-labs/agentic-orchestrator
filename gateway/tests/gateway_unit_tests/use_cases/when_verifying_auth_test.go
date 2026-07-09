package use_cases_test

import (
	"testing"
	"time"

	cb_errors "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenVerifyingAuthTestingSuite struct {
	sut *use_cases.VerifyAuth

	mockNonceStore             *interfaces_mocks.MockNonceStoreInterface
	mockSignedMessageValidator *interfaces_mocks.MockSignedMessageValidatorInterface
	mockJwtManager             *interfaces_mocks.MockJwtManagerInterface
}

func WhenVerifyingAuthBeforeEach(t *testing.T) *WhenVerifyingAuthTestingSuite {
	mockController := gomock.NewController(t)

	mockNonceStore := interfaces_mocks.NewMockNonceStoreInterface(mockController)
	mockSignedMessageValidator := interfaces_mocks.NewMockSignedMessageValidatorInterface(mockController)
	mockJwtManager := interfaces_mocks.NewMockJwtManagerInterface(mockController)

	sut := use_cases.NewVerifyAuth(mockNonceStore, mockSignedMessageValidator, mockJwtManager)

	return &WhenVerifyingAuthTestingSuite{
		sut: sut,

		mockNonceStore:             mockNonceStore,
		mockSignedMessageValidator: mockSignedMessageValidator,
		mockJwtManager:             mockJwtManager,
	}
}

func generateValidVerifyAuthRequest() *requests.VerifyAuthRequest {
	return &requests.VerifyAuthRequest{
		Address:   "0x4bca9E17a8326d283d1fC68A725B46b42c2D9146",
		Signature: "0x" + string(make132Hex()),
		Nonce:     "abc123",
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

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because address is malformed", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAuthBeforeEach(t)
			request := generateValidVerifyAuthRequest()
			request.Address = "not-an-address"

			_, err := suite.sut.Verify(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given the nonce is unknown or already used", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenVerifyingAuthTestingSuite) {
			suite.mockNonceStore.EXPECT().Consume(gomock.Any()).Return(false).Times(1)
		}

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAuthBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Verify(generateValidVerifyAuthRequest())

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given the signature does not match the address", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenVerifyingAuthTestingSuite, request *requests.VerifyAuthRequest) {
			suite.mockNonceStore.EXPECT().Consume(request.Nonce).Return(true).Times(1)
			suite.mockSignedMessageValidator.EXPECT().
				Validate(request.Nonce, request.Signature, request.Address).
				Return(false).
				Times(1)
		}

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAuthBeforeEach(t)
			request := generateValidVerifyAuthRequest()
			initSuite(suite, request)

			_, err := suite.sut.Verify(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given the signed nonce is valid", func(t *testing.T) {
		t.Parallel()

		expiresAt := time.Now().Add(time.Hour)

		initSuite := func(suite *WhenVerifyingAuthTestingSuite, request *requests.VerifyAuthRequest) {
			suite.mockNonceStore.EXPECT().Consume(request.Nonce).Return(true).Times(1)
			suite.mockSignedMessageValidator.EXPECT().
				Validate(request.Nonce, request.Signature, request.Address).
				Return(true).
				Times(1)
			suite.mockJwtManager.EXPECT().
				CreateJwtToken(request.Address).
				Return("jwt-token", expiresAt, nil).
				Times(1)
		}

		t.Run("Should return the token and expiry", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingAuthBeforeEach(t)
			request := generateValidVerifyAuthRequest()
			initSuite(suite, request)

			response, err := suite.sut.Verify(request)

			assert.NoError(t, err)
			assert.Equal(t, "jwt-token", response.Token)
			assert.Equal(t, expiresAt, response.ExpiresAt)
		})
	})
}
