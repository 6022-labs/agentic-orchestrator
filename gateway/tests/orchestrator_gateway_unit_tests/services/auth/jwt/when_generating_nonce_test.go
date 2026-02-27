package jwt_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/auth/jwt"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGeneratingNonceTestSuite struct {
	cacheStorageMock *interfaces_mocks.MockCacheStorageInterface

	sut *jwt.NonceGenerator
}

func WhenGeneratingNonceBeforeEach(t *testing.T) *WhenGeneratingNonceTestSuite {
	mockController := gomock.NewController(t)
	cacheStorageMock := interfaces_mocks.NewMockCacheStorageInterface(mockController)

	sut := jwt.NewNonceGenerator(cacheStorageMock)

	return &WhenGeneratingNonceTestSuite{
		cacheStorageMock: cacheStorageMock,

		sut: sut,
	}
}

func TestWhenGeneratingNonce(t *testing.T) {
	t.Parallel()

	t.Run("Given nonce generated", func(t *testing.T) {
		t.Parallel()

		t.Run("Should generate nonce", func(t *testing.T) {
			t.Parallel()

			suite := WhenGeneratingNonceBeforeEach(t)

			response1, _ := suite.sut.GenerateNonce()
			response2, _ := suite.sut.GenerateNonce()

			assert.NotEqual(t, response1, response2)
		})
	})
}
