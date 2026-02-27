package use_cases_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/requests"
	gateway_interfaces_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenInitializingJwtAuthenticationTestingSuite struct {
	sut *use_cases.InitializeJwtAuthentication

	mockCacheStorage   *interfaces_mocks.MockCacheStorageInterface
	mockNonceGenerator *gateway_interfaces_mocks.MockNonceGeneratorInterface
}

func WhenInitializingJwtAuthenticationBeforeEach(t *testing.T) *WhenInitializingJwtAuthenticationTestingSuite {
	mockController := gomock.NewController(t)

	mockCacheStorage := interfaces_mocks.NewMockCacheStorageInterface(mockController)
	mockNonceGenerator := gateway_interfaces_mocks.NewMockNonceGeneratorInterface(mockController)

	sut := use_cases.NewInitializeJwtAuthentication(
		mockCacheStorage,
		mockNonceGenerator,
	)

	return &WhenInitializingJwtAuthenticationTestingSuite{
		sut: sut,

		mockCacheStorage:   mockCacheStorage,
		mockNonceGenerator: mockNonceGenerator,
	}
}

func generateInitializeJwtAuthenticationRequest() *requests.InitializeJwtAuthenticationRequest {
	return &requests.InitializeJwtAuthenticationRequest{
		Address: "0x1234567890abcdef1234567890abcdef12345678",
	}
}

func TestWhenInitializingJwtAuthentication(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because address is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenInitializingJwtAuthenticationBeforeEach(t)

			request := generateInitializeJwtAuthenticationRequest()
			request.Address = ""
			_, err := suite.sut.Init(request)

			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because address is not valid ethereum address", func(t *testing.T) {
			t.Parallel()

			suite := WhenInitializingJwtAuthenticationBeforeEach(t)

			request := generateInitializeJwtAuthenticationRequest()
			request.Address = "0xinvalid-address"
			_, err := suite.sut.Init(request)

			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})
	})

	t.Run("Given error while generating nonce", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenInitializingJwtAuthenticationTestingSuite) {
			suite.mockNonceGenerator.EXPECT().GenerateNonce().Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()
			suite := WhenInitializingJwtAuthenticationBeforeEach(t)
			initSuite(suite)

			request := generateInitializeJwtAuthenticationRequest()

			_, err := suite.sut.Init(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given error while storing nonce in redis", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenInitializingJwtAuthenticationTestingSuite) {
			generatedNonce := "generated-nonce"
			suite.mockNonceGenerator.EXPECT().GenerateNonce().Return(&generatedNonce, nil)
			suite.mockCacheStorage.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return(assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenInitializingJwtAuthenticationBeforeEach(t)
			initSuite(suite)

			request := generateInitializeJwtAuthenticationRequest()

			_, err := suite.sut.Init(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given valid request and no error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenInitializingJwtAuthenticationTestingSuite) {
			generatedNonce := "generated-nonce"
			suite.mockNonceGenerator.EXPECT().GenerateNonce().Return(&generatedNonce, nil)
		}

		t.Run("Should store nonce in redis 5 minutes expiration", func(t *testing.T) {
			t.Parallel()

			var actualDuration *time.Duration

			suite := WhenInitializingJwtAuthenticationBeforeEach(t)
			initSuite(suite)
			suite.mockCacheStorage.
				EXPECT().
				Store(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(nonce string, address string, duration *time.Duration) error {
					actualDuration = duration
					return nil
				})

			request := generateInitializeJwtAuthenticationRequest()

			_, err := suite.sut.Init(request)

			assert.Nil(t, err)
			assert.NotNil(t, actualDuration)
			assert.Equal(t, time.Minute*5, *actualDuration)
		})

		t.Run("Should store nonce in redis with correct path", func(t *testing.T) {
			t.Parallel()

			var actualPath string

			suite := WhenInitializingJwtAuthenticationBeforeEach(t)
			initSuite(suite)
			suite.mockCacheStorage.
				EXPECT().
				Store(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(path string, value string, duration *time.Duration) error {
					actualPath = path
					return nil
				})

			request := generateInitializeJwtAuthenticationRequest()

			_, err := suite.sut.Init(request)

			expectedPath := fmt.Sprintf(use_cases.NonceLoginRedisPathTemplate, request.Address)

			assert.Nil(t, err)
			assert.NotEmpty(t, actualPath)
			assert.Equal(t, expectedPath, actualPath)
		})

		t.Run("Should return nonce", func(t *testing.T) {
			t.Parallel()

			suite := WhenInitializingJwtAuthenticationBeforeEach(t)
			initSuite(suite)
			suite.mockCacheStorage.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			request := generateInitializeJwtAuthenticationRequest()

			response, err := suite.sut.Init(request)

			assert.Nil(t, err)
			assert.NotNil(t, response)
			assert.NotEmpty(t, response.Nonce)
		})
	})
}
