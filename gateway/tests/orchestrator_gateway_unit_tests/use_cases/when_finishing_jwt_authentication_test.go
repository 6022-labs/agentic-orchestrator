package use_cases_test

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/requests"
	gateway_interfaces_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	chain_id    uint64 = 80002
	public_key  string = "0xd389d997Df5654b87914492E1e0C6861bfe835DE"
	private_key string = "b4c2dd8c816efd854aecde7167b0a5db630854f426c50eb2b8784b9a89c9f95a"
)

type WhenFinishingJwtAuthenticationTestingSuite struct {
	sut *use_cases.FinishJwtAuthentication

	mockCacheStorage           *interfaces_mocks.MockCacheStorageInterface
	mockJwtManager             *gateway_interfaces_mocks.MockJwtManagerInterface
	mockSIWEDetailsExtractor   *gateway_interfaces_mocks.MockSIWEDetailsExtractorInterface
	mockSignedMessageValidator *gateway_interfaces_mocks.MockSignedMessageValidatorInterface
}

func WhenFinishingJwtAuthenticationBeforeEach(t *testing.T) *WhenFinishingJwtAuthenticationTestingSuite {
	mockController := gomock.NewController(t, gomock.WithOverridableExpectations())

	mockCacheStorage := interfaces_mocks.NewMockCacheStorageInterface(mockController)
	mockJwtManager := gateway_interfaces_mocks.NewMockJwtManagerInterface(mockController)
	mockSIWEDetailsExtractor := gateway_interfaces_mocks.NewMockSIWEDetailsExtractorInterface(mockController)
	mockSignedMessageValidator := gateway_interfaces_mocks.NewMockSignedMessageValidatorInterface(mockController)

	sut := use_cases.NewFinishJwtAuthentication(
		mockCacheStorage,
		mockJwtManager,
		mockSIWEDetailsExtractor,
		mockSignedMessageValidator,
	)

	return &WhenFinishingJwtAuthenticationTestingSuite{
		sut: sut,

		mockJwtManager:             mockJwtManager,
		mockCacheStorage:           mockCacheStorage,
		mockSIWEDetailsExtractor:   mockSIWEDetailsExtractor,
		mockSignedMessageValidator: mockSignedMessageValidator,
	}
}

func generateMessage(nonce string) string {
	return strings.Join([]string{
		"localhost:5173 wants you to sign in with your Ethereum account:\n",
		public_key,
		"\n\nPlease sign with your account\n\nURI: http://localhost:5173\nVersion: 1\nChain ID: ",
		fmt.Sprintf("%d", chain_id),
		"\nNonce: ",
		nonce,
		"\nIssued At: 2025-03-02T18:24:14.110Z",
	}, "")
}

func signMessage(message string) (string, error) {
	privateKey, err := crypto.HexToECDSA(private_key)
	if err != nil {
		return "", err
	}

	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(prefixedMessage))
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(signature), nil
}

func generateValidFinishJwtAuthenticationRequest(nonce string) *requests.FinishJwtAuthenticationRequest {
	message := generateMessage(nonce)
	signedMessage, _ := signMessage(message)
	return &requests.FinishJwtAuthenticationRequest{
		Message:   message,
		Signature: signedMessage,
	}
}

func TestWhenFinishingJwtAuthentication(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because message is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenFinishingJwtAuthenticationBeforeEach(t)

			request := generateValidFinishJwtAuthenticationRequest("")
			request.Message = ""
			_, err := suite.sut.Finish(request)

			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because signature is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenFinishingJwtAuthenticationBeforeEach(t)

			request := generateValidFinishJwtAuthenticationRequest("")
			request.Signature = ""
			_, err := suite.sut.Finish(request)

			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})
	})

	t.Run("Given error while extracting siwe message details", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenFinishingJwtAuthenticationTestingSuite) {
			suite.mockSIWEDetailsExtractor.EXPECT().ExtractSIWEMessageDetails(gomock.Any()).Return(nil, nil, nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenFinishingJwtAuthenticationBeforeEach(t)
			initSuite(suite)

			request := generateValidFinishJwtAuthenticationRequest("nonce")
			_, err := suite.sut.Finish(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given error while reading nonce from redis", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenFinishingJwtAuthenticationTestingSuite, nonce string) {
			suite.mockSIWEDetailsExtractor.EXPECT().ExtractSIWEMessageDetails(gomock.Any()).Return(&public_key, &chain_id, &nonce, nil)
			suite.mockCacheStorage.EXPECT().Read(gomock.Any()).Return("", assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			nonce := "nonce"

			request := generateValidFinishJwtAuthenticationRequest(nonce)

			suite := WhenFinishingJwtAuthenticationBeforeEach(t)
			initSuite(suite, nonce)

			_, err := suite.sut.Finish(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given invalid nonce in redis does not match the one received", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenFinishingJwtAuthenticationTestingSuite, nonce string) {
			suite.mockSIWEDetailsExtractor.EXPECT().ExtractSIWEMessageDetails(gomock.Any()).Return(&public_key, &chain_id, &nonce, nil)
			suite.mockCacheStorage.EXPECT().Read(gomock.Any()).Return(nonce, nil)
			suite.mockSignedMessageValidator.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			baseNonce := "nonce"
			signedNonce := "invalid_signed_nonce"

			suite := WhenFinishingJwtAuthenticationBeforeEach(t)
			initSuite(suite, baseNonce)

			request := generateValidFinishJwtAuthenticationRequest(signedNonce)
			_, err := suite.sut.Finish(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given error while creating jwt token", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenFinishingJwtAuthenticationTestingSuite, nonce string) {
			suite.mockSIWEDetailsExtractor.EXPECT().ExtractSIWEMessageDetails(gomock.Any()).Return(&public_key, &chain_id, &nonce, nil)
			suite.mockCacheStorage.EXPECT().Read(gomock.Any()).Return(nonce, nil)
			suite.mockSignedMessageValidator.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
			suite.mockJwtManager.EXPECT().CreateJwtToken(gomock.Any(), gomock.Any()).Return("", assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			nonce := "nonce"

			suite := WhenFinishingJwtAuthenticationBeforeEach(t)
			initSuite(suite, nonce)

			request := generateValidFinishJwtAuthenticationRequest(nonce)
			_, err := suite.sut.Finish(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given valid request and signed nonce", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenFinishingJwtAuthenticationTestingSuite, nonce string) {
			suite.mockSIWEDetailsExtractor.EXPECT().ExtractSIWEMessageDetails(gomock.Any()).Return(&public_key, &chain_id, &nonce, nil)
			suite.mockCacheStorage.EXPECT().Read(gomock.Any()).Return(nonce, nil)
			suite.mockSignedMessageValidator.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
			suite.mockJwtManager.EXPECT().CreateJwtToken(gomock.Any(), gomock.Any()).Return("jwt", nil)
			suite.mockCacheStorage.EXPECT().Delete(gomock.Any()).Return(nil)
		}

		t.Run("Should return jwt token", func(t *testing.T) {
			t.Parallel()

			nonce := "nonce"

			suite := WhenFinishingJwtAuthenticationBeforeEach(t)
			initSuite(suite, nonce)

			request := generateValidFinishJwtAuthenticationRequest(nonce)
			response, err := suite.sut.Finish(request)

			assert.Nil(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, "jwt", response.Jwt)
		})

		t.Run("Should delete nonce from redis", func(t *testing.T) {
			t.Parallel()

			nonce := "nonce"

			suite := WhenFinishingJwtAuthenticationBeforeEach(t)
			initSuite(suite, nonce)

			suite.mockCacheStorage.
				EXPECT().
				Delete(gomock.Any()).
				DoAndReturn(func(key string) error {
					nonceRedisPath := fmt.Sprintf(use_cases.NonceLoginRedisPathTemplate, public_key)

					assert.Equal(t, nonceRedisPath, key)
					return nil
				})

			request := generateValidFinishJwtAuthenticationRequest(nonce)
			_, err := suite.sut.Finish(request)

			assert.Nil(t, err)
		})
	})
}
